package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

// =============================================================================
// SQL Query Execution via PostgREST
// =============================================================================
//
// These methods accept actual SQL query strings and translate them into
// PostgREST API calls under the hood. Write your queries like real SQL:
//
//     // SELECT single row
//     var user models.User
//     err := client.ExecuteQueryRow("SELECT * FROM users WHERE id = $1", &user, userID)
//
//     // SELECT multiple rows
//     var users []models.User
//     err := client.ExecuteQuery("SELECT * FROM users WHERE role = $1 ORDER BY name ASC LIMIT 10", &users, "admin")
//
//     // SELECT with OR conditions
//     var companies []models.Company
//     err := client.ExecuteQuery("SELECT * FROM companies WHERE is_active = $1 AND (symbol ILIKE $2 OR name ILIKE $3) ORDER BY market_cap DESC LIMIT $4", &companies, true, "%NABIL%", "%NABIL%", 10)
//
//     // INSERT and return the inserted row
//     var user models.User
//     err := client.ExecuteInsert("INSERT INTO users (full_name, email) VALUES ($1, $2) RETURNING *", &user, "John", "john@example.com")
//
//     // UPDATE with RETURNING
//     var user models.User
//     err := client.ExecuteUpdate("UPDATE users SET full_name = $1 WHERE id = $2 RETURNING *", &user, "Jane", userID)
//
//     // DELETE
//     err := client.ExecuteDelete("DELETE FROM users WHERE id = $1", userID)
//
// =============================================================================

// sqlParts holds the parsed components of a SQL query
type sqlParts struct {
	table     string
	columns   string
	filters   []sqlFilter
	orFilters []orGroup
	orderBy   string
	ascending bool
	limit     int
	offset    int
}

type sqlFilter struct {
	column   string
	operator string // eq, gt, gte, lt, lte, neq, ilike, like
	value    interface{}
}

type orGroup struct {
	conditions []sqlFilter
}

// parseSelectQuery parses a SQL SELECT query and extracts its components
func parseSelectQuery(query string, args []interface{}) (*sqlParts, error) {
	parts := &sqlParts{columns: "*"}

	// Normalize whitespace
	q := strings.TrimSpace(query)
	q = regexp.MustCompile(`\s+`).ReplaceAllString(q, " ")

	// Must start with SELECT
	upper := strings.ToUpper(q)
	if !strings.HasPrefix(upper, "SELECT ") {
		return nil, fmt.Errorf("not a SELECT query")
	}

	// Extract columns: SELECT <columns> FROM
	fromIdx := caseInsensitiveIndex(q, " FROM ")
	if fromIdx < 0 {
		return nil, fmt.Errorf("missing FROM clause")
	}
	parts.columns = strings.TrimSpace(q[7:fromIdx])

	// Extract table name: FROM <table> [WHERE|ORDER|LIMIT|$]
	rest := strings.TrimSpace(q[fromIdx+6:])
	tableEnd := len(rest)
	for _, kw := range []string{" WHERE ", " ORDER ", " LIMIT ", " OFFSET "} {
		idx := caseInsensitiveIndex(rest, kw)
		if idx >= 0 && idx < tableEnd {
			tableEnd = idx
		}
	}
	parts.table = strings.TrimSpace(rest[:tableEnd])
	rest = strings.TrimSpace(rest[tableEnd:])

	// Extract WHERE clause
	if idx := caseInsensitiveIndex(rest, "WHERE "); idx >= 0 {
		whereRest := rest[idx+6:]
		// Find end of WHERE clause (ORDER BY, LIMIT, OFFSET, or end)
		whereEnd := len(whereRest)
		for _, kw := range []string{" ORDER ", " LIMIT ", " OFFSET "} {
			kwIdx := caseInsensitiveIndex(whereRest, kw)
			if kwIdx >= 0 && kwIdx < whereEnd {
				whereEnd = kwIdx
			}
		}
		whereClause := strings.TrimSpace(whereRest[:whereEnd])
		rest = strings.TrimSpace(whereRest[whereEnd:])

		// Parse WHERE conditions
		filters, orFilters, err := parseWhereClause(whereClause, args)
		if err != nil {
			return nil, err
		}
		parts.filters = filters
		parts.orFilters = orFilters
	}

	// Extract ORDER BY
	if idx := caseInsensitiveIndex(rest, "ORDER BY "); idx >= 0 {
		orderRest := rest[idx+9:]
		orderEnd := len(orderRest)
		for _, kw := range []string{" LIMIT ", " OFFSET "} {
			kwIdx := caseInsensitiveIndex(orderRest, kw)
			if kwIdx >= 0 && kwIdx < orderEnd {
				orderEnd = kwIdx
			}
		}
		orderClause := strings.TrimSpace(orderRest[:orderEnd])
		rest = strings.TrimSpace(orderRest[orderEnd:])

		// Parse "column ASC/DESC"
		orderParts := strings.Fields(orderClause)
		if len(orderParts) >= 1 {
			parts.orderBy = orderParts[0]
			parts.ascending = true // default ASC
			if len(orderParts) >= 2 && strings.ToUpper(orderParts[1]) == "DESC" {
				parts.ascending = false
			}
		}
	}

	// Extract LIMIT
	if idx := caseInsensitiveIndex(rest, "LIMIT "); idx >= 0 {
		limitRest := rest[idx+6:]
		limitEnd := len(limitRest)
		if offIdx := caseInsensitiveIndex(limitRest, " OFFSET "); offIdx >= 0 && offIdx < limitEnd {
			limitEnd = offIdx
		}
		limitVal := strings.TrimSpace(limitRest[:limitEnd])
		rest = strings.TrimSpace(limitRest[limitEnd:])

		// Could be a $N placeholder or a literal number
		parts.limit = resolveIntParam(limitVal, args)
	}

	// Extract OFFSET
	if idx := caseInsensitiveIndex(rest, "OFFSET "); idx >= 0 {
		offsetVal := strings.TrimSpace(rest[idx+7:])
		parts.offset = resolveIntParam(offsetVal, args)
	}

	return parts, nil
}

// parseWhereClause parses WHERE conditions with AND, OR, and comparison operators
func parseWhereClause(clause string, args []interface{}) ([]sqlFilter, []orGroup, error) {
	var filters []sqlFilter
	var orGroups []orGroup

	// Handle parenthesized OR groups: ... AND (col1 ILIKE $2 OR col2 ILIKE $3) AND ...
	// Split by top-level AND first, respecting parentheses
	conditions := splitByTopLevelAND(clause)

	for _, cond := range conditions {
		cond = strings.TrimSpace(cond)
		if cond == "" {
			continue
		}

		// Check if it's a parenthesized OR group: (cond1 OR cond2)
		if strings.HasPrefix(cond, "(") && strings.HasSuffix(cond, ")") {
			inner := cond[1 : len(cond)-1]
			orParts := splitByOR(inner)
			if len(orParts) > 1 {
				var group orGroup
				for _, orPart := range orParts {
					f, err := parseSingleCondition(strings.TrimSpace(orPart), args)
					if err != nil {
						return nil, nil, err
					}
					group.conditions = append(group.conditions, f)
				}
				orGroups = append(orGroups, group)
				continue
			}
			// Single condition in parens, just parse normally
			cond = inner
		}

		f, err := parseSingleCondition(cond, args)
		if err != nil {
			return nil, nil, err
		}
		filters = append(filters, f)
	}

	return filters, orGroups, nil
}

// parseSingleCondition parses: column OPERATOR $N or column OPERATOR value
func parseSingleCondition(cond string, args []interface{}) (sqlFilter, error) {
	cond = strings.TrimSpace(cond)

	// Supported operators: =, !=, <>, >, >=, <, <=, ILIKE, LIKE
	operators := []struct {
		sql     string
		postgOp string
	}{
		{">=", "gte"},
		{"<=", "lte"},
		{"!=", "neq"},
		{"<>", "neq"},
		{">", "gt"},
		{"<", "lt"},
		{"=", "eq"},
	}

	// Try text operators first (case-insensitive)
	upperCond := strings.ToUpper(cond)
	for _, textOp := range []struct {
		sql     string
		postgOp string
	}{
		{" ILIKE ", "ilike"},
		{" LIKE ", "like"},
	} {
		idx := strings.Index(upperCond, textOp.sql)
		if idx >= 0 {
			col := strings.TrimSpace(cond[:idx])
			val := strings.TrimSpace(cond[idx+len(textOp.sql):])
			resolved := resolveParam(val, args)
			return sqlFilter{column: col, operator: textOp.postgOp, value: resolved}, nil
		}
	}

	// Try symbol operators
	for _, op := range operators {
		idx := strings.Index(cond, op.sql)
		if idx >= 0 {
			col := strings.TrimSpace(cond[:idx])
			val := strings.TrimSpace(cond[idx+len(op.sql):])
			resolved := resolveParam(val, args)
			return sqlFilter{column: col, operator: op.postgOp, value: resolved}, nil
		}
	}

	return sqlFilter{}, fmt.Errorf("could not parse condition: %s", cond)
}

// resolveParam resolves $N placeholders to actual arg values
func resolveParam(val string, args []interface{}) interface{} {
	val = strings.TrimSpace(val)
	if strings.HasPrefix(val, "$") {
		idx := 0
		fmt.Sscanf(val, "$%d", &idx)
		if idx > 0 && idx <= len(args) {
			return args[idx-1]
		}
	}
	// Remove surrounding quotes if any
	if (strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'")) ||
		(strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"")) {
		return val[1 : len(val)-1]
	}
	return val
}

// resolveIntParam resolves a parameter to an int
func resolveIntParam(val string, args []interface{}) int {
	val = strings.TrimSpace(val)
	if strings.HasPrefix(val, "$") {
		idx := 0
		fmt.Sscanf(val, "$%d", &idx)
		if idx > 0 && idx <= len(args) {
			switch v := args[idx-1].(type) {
			case int:
				return v
			case int64:
				return int(v)
			case float64:
				return int(v)
			}
		}
	}
	n := 0
	fmt.Sscanf(val, "%d", &n)
	return n
}

// splitByTopLevelAND splits a WHERE clause by AND, respecting parentheses
func splitByTopLevelAND(clause string) []string {
	var parts []string
	depth := 0
	current := ""
	upper := strings.ToUpper(clause)

	for i := 0; i < len(clause); i++ {
		if clause[i] == '(' {
			depth++
			current += string(clause[i])
		} else if clause[i] == ')' {
			depth--
			current += string(clause[i])
		} else if depth == 0 && i+5 <= len(upper) && upper[i:i+5] == " AND " {
			parts = append(parts, strings.TrimSpace(current))
			current = ""
			i += 4 // skip " AND "
		} else {
			current += string(clause[i])
		}
	}
	if current != "" {
		parts = append(parts, strings.TrimSpace(current))
	}
	return parts
}

// splitByOR splits a clause by OR (not inside parentheses)
func splitByOR(clause string) []string {
	var parts []string
	depth := 0
	current := ""
	upper := strings.ToUpper(clause)

	for i := 0; i < len(clause); i++ {
		if clause[i] == '(' {
			depth++
			current += string(clause[i])
		} else if clause[i] == ')' {
			depth--
			current += string(clause[i])
		} else if depth == 0 && i+4 <= len(upper) && upper[i:i+4] == " OR " {
			parts = append(parts, strings.TrimSpace(current))
			current = ""
			i += 3 // skip " OR "
		} else {
			current += string(clause[i])
		}
	}
	if current != "" {
		parts = append(parts, strings.TrimSpace(current))
	}
	return parts
}

// caseInsensitiveIndex finds the index of substr in s, case-insensitive
func caseInsensitiveIndex(s, substr string) int {
	return strings.Index(strings.ToUpper(s), strings.ToUpper(substr))
}

// parseInsertQuery extracts table name and column/value pairs from INSERT INTO ... VALUES ...
func parseInsertQuery(query string, args []interface{}) (string, map[string]interface{}, error) {
	q := strings.TrimSpace(query)
	q = regexp.MustCompile(`\s+`).ReplaceAllString(q, " ")

	upper := strings.ToUpper(q)
	if !strings.HasPrefix(upper, "INSERT INTO ") {
		return "", nil, fmt.Errorf("not an INSERT query")
	}

	// Extract table: INSERT INTO <table> (
	rest := q[12:]
	parenIdx := strings.Index(rest, "(")
	if parenIdx < 0 {
		return "", nil, fmt.Errorf("missing column list in INSERT")
	}
	table := strings.TrimSpace(rest[:parenIdx])

	// Extract columns: (...columns...) VALUES (...)
	rest = rest[parenIdx:]
	closeParen := strings.Index(rest, ")")
	if closeParen < 0 {
		return "", nil, fmt.Errorf("malformed INSERT column list")
	}
	colStr := rest[1:closeParen]
	columns := strings.Split(colStr, ",")
	for i := range columns {
		columns[i] = strings.TrimSpace(columns[i])
	}

	// Extract values: VALUES ($1, $2, ...)
	valuesIdx := caseInsensitiveIndex(rest, "VALUES ")
	if valuesIdx < 0 {
		return "", nil, fmt.Errorf("missing VALUES clause")
	}
	valuesRest := rest[valuesIdx+7:]
	// Find ( and )
	vOpen := strings.Index(valuesRest, "(")
	vClose := strings.LastIndex(valuesRest, ")")
	if vOpen < 0 || vClose < 0 {
		return "", nil, fmt.Errorf("malformed VALUES clause")
	}
	valStr := valuesRest[vOpen+1 : vClose]
	values := strings.Split(valStr, ",")

	if len(columns) != len(values) {
		return "", nil, fmt.Errorf("column count (%d) != value count (%d)", len(columns), len(values))
	}

	data := make(map[string]interface{})
	for i, col := range columns {
		val := strings.TrimSpace(values[i])
		data[col] = resolveParam(val, args)
	}

	return table, data, nil
}

// parseUpdateQuery extracts table, SET data, and WHERE filters from UPDATE ... SET ... WHERE ...
func parseUpdateQuery(query string, args []interface{}) (string, map[string]interface{}, map[string]interface{}, error) {
	q := strings.TrimSpace(query)
	q = regexp.MustCompile(`\s+`).ReplaceAllString(q, " ")

	upper := strings.ToUpper(q)
	if !strings.HasPrefix(upper, "UPDATE ") {
		return "", nil, nil, fmt.Errorf("not an UPDATE query")
	}

	// Extract table: UPDATE <table> SET
	setIdx := caseInsensitiveIndex(q, " SET ")
	if setIdx < 0 {
		return "", nil, nil, fmt.Errorf("missing SET clause")
	}
	table := strings.TrimSpace(q[7:setIdx])

	// Extract SET clause: SET col1 = $1, col2 = $2 WHERE ...
	rest := q[setIdx+5:]
	whereIdx := caseInsensitiveIndex(rest, " WHERE ")
	var setClause, whereClause string
	if whereIdx >= 0 {
		setClause = strings.TrimSpace(rest[:whereIdx])
		whereRest := rest[whereIdx+7:]
		// Remove RETURNING if present
		retIdx := caseInsensitiveIndex(whereRest, " RETURNING ")
		if retIdx >= 0 {
			whereClause = strings.TrimSpace(whereRest[:retIdx])
		} else {
			whereClause = strings.TrimSpace(whereRest)
		}
	} else {
		retIdx := caseInsensitiveIndex(rest, " RETURNING ")
		if retIdx >= 0 {
			setClause = strings.TrimSpace(rest[:retIdx])
		} else {
			setClause = strings.TrimSpace(rest)
		}
	}

	// Parse SET pairs: col1 = $1, col2 = $2
	data := make(map[string]interface{})
	setPairs := strings.Split(setClause, ",")
	for _, pair := range setPairs {
		eqIdx := strings.Index(pair, "=")
		if eqIdx < 0 {
			continue
		}
		col := strings.TrimSpace(pair[:eqIdx])
		val := strings.TrimSpace(pair[eqIdx+1:])
		data[col] = resolveParam(val, args)
	}

	// Parse WHERE filters
	filters := make(map[string]interface{})
	if whereClause != "" {
		whereParts := splitByTopLevelAND(whereClause)
		for _, part := range whereParts {
			f, err := parseSingleCondition(part, args)
			if err != nil {
				return "", nil, nil, err
			}
			filters[f.column] = f.value
		}
	}

	return table, data, filters, nil
}

// parseDeleteQuery extracts table and WHERE filters from DELETE FROM ... WHERE ...
func parseDeleteQuery(query string, args []interface{}) (string, map[string]interface{}, error) {
	q := strings.TrimSpace(query)
	q = regexp.MustCompile(`\s+`).ReplaceAllString(q, " ")

	upper := strings.ToUpper(q)
	if !strings.HasPrefix(upper, "DELETE FROM ") {
		return "", nil, fmt.Errorf("not a DELETE query")
	}

	// Extract table: DELETE FROM <table> WHERE
	rest := q[12:]
	whereIdx := caseInsensitiveIndex(rest, " WHERE ")
	if whereIdx < 0 {
		table := strings.TrimSpace(rest)
		return table, map[string]interface{}{}, nil
	}
	table := strings.TrimSpace(rest[:whereIdx])
	whereClause := strings.TrimSpace(rest[whereIdx+7:])

	// Parse WHERE filters
	filters := make(map[string]interface{})
	whereParts := splitByTopLevelAND(whereClause)
	for _, part := range whereParts {
		f, err := parseSingleCondition(part, args)
		if err != nil {
			return "", nil, err
		}
		filters[f.column] = f.value
	}

	return table, filters, nil
}

// =============================================================================
// ExecuteQueryRow — Execute a SELECT query returning a single row
//
// Usage:
//
//	var user models.User
//	err := client.ExecuteQueryRow("SELECT * FROM users WHERE id = $1", &user, userID)
//
// =============================================================================
func (c *Client) ExecuteQueryRow(query string, dest interface{}, args ...interface{}) error {
	parts, err := parseSelectQuery(query, args)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	params := url.Values{}
	params.Set("select", parts.columns)

	// Apply filters
	for _, f := range parts.filters {
		params.Add(f.column, fmt.Sprintf("%s.%v", f.operator, f.value))
	}

	// Apply OR groups
	for _, og := range parts.orFilters {
		var orParts []string
		for _, f := range og.conditions {
			orParts = append(orParts, fmt.Sprintf("%s.%s.%v", f.column, f.operator, f.value))
		}
		params.Add("or", fmt.Sprintf("(%s)", strings.Join(orParts, ",")))
	}

	urlStr := fmt.Sprintf("%s?%s", c.TableURL(parts.table), params.Encode())

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.SetHeaders(req)
	req.Header.Set("Accept", "application/vnd.pgrst.object+json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode == 406 {
		return fmt.Errorf("not found")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("query error (status %d): %s", resp.StatusCode, string(respBody))
	}

	if dest != nil {
		if err := json.Unmarshal(respBody, dest); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// =============================================================================
// ExecuteQuery — Execute a SELECT query returning multiple rows
//
// Usage:
//
//	var users []models.User
//	err := client.ExecuteQuery("SELECT * FROM users WHERE role = $1 ORDER BY name ASC LIMIT 10", &users, "admin")
//
// =============================================================================
func (c *Client) ExecuteQuery(query string, dest interface{}, args ...interface{}) error {
	parts, err := parseSelectQuery(query, args)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	params := url.Values{}
	params.Set("select", parts.columns)

	// Apply filters
	for _, f := range parts.filters {
		params.Add(f.column, fmt.Sprintf("%s.%v", f.operator, f.value))
	}

	// Apply OR groups
	for _, og := range parts.orFilters {
		var orParts []string
		for _, f := range og.conditions {
			orParts = append(orParts, fmt.Sprintf("%s.%s.%v", f.column, f.operator, f.value))
		}
		params.Add("or", fmt.Sprintf("(%s)", strings.Join(orParts, ",")))
	}

	if parts.orderBy != "" {
		direction := "asc"
		if !parts.ascending {
			direction = "desc"
		}
		params.Set("order", fmt.Sprintf("%s.%s", parts.orderBy, direction))
	}

	if parts.limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", parts.limit))
	}

	if parts.offset > 0 {
		params.Set("offset", fmt.Sprintf("%d", parts.offset))
	}

	urlStr := fmt.Sprintf("%s?%s", c.TableURL(parts.table), params.Encode())

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.SetHeaders(req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("query error (status %d): %s", resp.StatusCode, string(respBody))
	}

	if dest != nil {
		if err := json.Unmarshal(respBody, dest); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// =============================================================================
// ExecuteInsert — Execute an INSERT query and return the created row
//
// Usage:
//
//	var user models.User
//	err := client.ExecuteInsert("INSERT INTO users (full_name, email) VALUES ($1, $2) RETURNING *", &user, "John", "john@example.com")
//
// =============================================================================
func (c *Client) ExecuteInsert(query string, dest interface{}, args ...interface{}) error {
	table, data, err := parseInsertQuery(query, args)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	urlStr := fmt.Sprintf("%s?select=*", c.TableURL(table))

	req, err := http.NewRequest("POST", urlStr, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.SetHeaders(req)
	req.Header.Set("Prefer", "return=representation")
	req.Header.Set("Accept", "application/vnd.pgrst.object+json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("insert error (status %d): %s", resp.StatusCode, string(respBody))
	}

	if dest != nil {
		if err := json.Unmarshal(respBody, dest); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// =============================================================================
// ExecuteUpdate — Execute an UPDATE query and return the updated row
//
// Usage:
//
//	var user models.User
//	err := client.ExecuteUpdate("UPDATE users SET full_name = $1 WHERE id = $2 RETURNING *", &user, "Jane", userID)
//
// =============================================================================
func (c *Client) ExecuteUpdate(query string, dest interface{}, args ...interface{}) error {
	table, data, filters, err := parseUpdateQuery(query, args)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	params := url.Values{}
	params.Set("select", "*")

	for col, val := range filters {
		params.Add(col, fmt.Sprintf("eq.%v", val))
	}

	urlStr := fmt.Sprintf("%s?%s", c.TableURL(table), params.Encode())

	req, err := http.NewRequest("PATCH", urlStr, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.SetHeaders(req)
	req.Header.Set("Prefer", "return=representation")
	req.Header.Set("Accept", "application/vnd.pgrst.object+json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode == 406 {
		return fmt.Errorf("not found")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("update error (status %d): %s", resp.StatusCode, string(respBody))
	}

	if dest != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, dest); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// =============================================================================
// ExecuteDelete — Execute a DELETE query
//
// Usage:
//
//	err := client.ExecuteDelete("DELETE FROM users WHERE id = $1", userID)
//
// =============================================================================
func (c *Client) ExecuteDelete(query string, args ...interface{}) error {
	table, filters, err := parseDeleteQuery(query, args)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	params := url.Values{}
	for col, val := range filters {
		params.Add(col, fmt.Sprintf("eq.%v", val))
	}

	urlStr := fmt.Sprintf("%s?%s", c.TableURL(table), params.Encode())

	req, err := http.NewRequest("DELETE", urlStr, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.SetHeaders(req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// =============================================================================
// Supabase Client
// =============================================================================

type Client struct {
	BaseURL    string       // e.g. "https://abc123.supabase.co"
	APIKey     string       // Your Supabase anon key or service_role key
	HTTPClient *http.Client // Shared HTTP client with timeout
}

// NewClient creates a new Supabase client from environment variables.
func NewClient() (*Client, error) {
	baseURL := getEnv("SUPABASE_URL", "")
	apiKey := getEnv("SUPABASE_ANON_KEY", "")

	if baseURL == "" {
		return nil, fmt.Errorf("SUPABASE_URL is required")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("SUPABASE_ANON_KEY is required")
	}

	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// =============================================================================
// URL Builders
// =============================================================================

// TableURL builds the URL for a table endpoint
// Example: "https://abc123.supabase.co/rest/v1/users"
func (c *Client) TableURL(table string) string {
	return fmt.Sprintf("%s/rest/v1/%s", c.BaseURL, table)
}

// StorageURL builds the Supabase Storage base URL.
// Example: "https://abc123.supabase.co/storage/v1"
func (c *Client) StorageURL() string {
	return fmt.Sprintf("%s/storage/v1", c.BaseURL)
}

// SetHeaders attaches the required Supabase authentication headers to a request.
func (c *Client) SetHeaders(req *http.Request) {
	req.Header.Set("apikey", c.APIKey)
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
}

// =============================================================================
// Helper function
// =============================================================================

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
