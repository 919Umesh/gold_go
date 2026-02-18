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
	operator string 
	value    interface{}
}

type orGroup struct {
	conditions []sqlFilter
}

func parseSelectQuery(query string, args []interface{}) (*sqlParts, error) {
	parts := &sqlParts{columns: "*"}

	q := strings.TrimSpace(query)
	q = regexp.MustCompile(`\s+`).ReplaceAllString(q, " ")

	upper := strings.ToUpper(q)
	if !strings.HasPrefix(upper, "SELECT ") {
		return nil, fmt.Errorf("not a SELECT query")
	}

	fromIdx := caseInsensitiveIndex(q, " FROM ")
	if fromIdx < 0 {
		return nil, fmt.Errorf("missing FROM clause")
	}
	parts.columns = strings.TrimSpace(q[7:fromIdx])

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

	
	if idx := caseInsensitiveIndex(rest, "WHERE "); idx >= 0 {
		whereRest := rest[idx+6:]
		whereEnd := len(whereRest)
		for _, kw := range []string{" ORDER ", " LIMIT ", " OFFSET "} {
			kwIdx := caseInsensitiveIndex(whereRest, kw)
			if kwIdx >= 0 && kwIdx < whereEnd {
				whereEnd = kwIdx
			}
		}
		whereClause := strings.TrimSpace(whereRest[:whereEnd])
		rest = strings.TrimSpace(whereRest[whereEnd:])

		filters, orFilters, err := parseWhereClause(whereClause, args)
		if err != nil {
			return nil, err
		}
		parts.filters = filters
		parts.orFilters = orFilters
	}

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

		orderParts := strings.Fields(orderClause)
		if len(orderParts) >= 1 {
			parts.orderBy = orderParts[0]
			parts.ascending = true 
			if len(orderParts) >= 2 && strings.ToUpper(orderParts[1]) == "DESC" {
				parts.ascending = false
			}
		}
	}

	if idx := caseInsensitiveIndex(rest, "LIMIT "); idx >= 0 {
		limitRest := rest[idx+6:]
		limitEnd := len(limitRest)
		if offIdx := caseInsensitiveIndex(limitRest, " OFFSET "); offIdx >= 0 && offIdx < limitEnd {
			limitEnd = offIdx
		}
		limitVal := strings.TrimSpace(limitRest[:limitEnd])
		rest = strings.TrimSpace(limitRest[limitEnd:])

		parts.limit = resolveIntParam(limitVal, args)
	}

	if idx := caseInsensitiveIndex(rest, "OFFSET "); idx >= 0 {
		offsetVal := strings.TrimSpace(rest[idx+7:])
		parts.offset = resolveIntParam(offsetVal, args)
	}

	return parts, nil
}


func parseWhereClause(clause string, args []interface{}) ([]sqlFilter, []orGroup, error) {
	var filters []sqlFilter
	var orGroups []orGroup

	conditions := splitByTopLevelAND(clause)

	for _, cond := range conditions {
		cond = strings.TrimSpace(cond)
		if cond == "" {
			continue
		}

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

func parseSingleCondition(cond string, args []interface{}) (sqlFilter, error) {
	cond = strings.TrimSpace(cond)

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

func resolveParam(val string, args []interface{}) interface{} {
	val = strings.TrimSpace(val)
	if strings.HasPrefix(val, "$") {
		idx := 0
		fmt.Sscanf(val, "$%d", &idx)
		if idx > 0 && idx <= len(args) {
			return args[idx-1]
		}
	}
	if (strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'")) ||
		(strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"")) {
		return val[1 : len(val)-1]
	}
	return val
}

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
			i += 4 
		} else {
			current += string(clause[i])
		}
	}
	if current != "" {
		parts = append(parts, strings.TrimSpace(current))
	}
	return parts
}

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
			i += 3 
		} else {
			current += string(clause[i])
		}
	}
	if current != "" {
		parts = append(parts, strings.TrimSpace(current))
	}
	return parts
}


func caseInsensitiveIndex(s, substr string) int {
	return strings.Index(strings.ToUpper(s), strings.ToUpper(substr))
}


func parseInsertQuery(query string, args []interface{}) (string, map[string]interface{}, error) {
	q := strings.TrimSpace(query)
	q = regexp.MustCompile(`\s+`).ReplaceAllString(q, " ")

	upper := strings.ToUpper(q)
	if !strings.HasPrefix(upper, "INSERT INTO ") {
		return "", nil, fmt.Errorf("not an INSERT query")
	}

	rest := q[12:]
	parenIdx := strings.Index(rest, "(")
	if parenIdx < 0 {
		return "", nil, fmt.Errorf("missing column list in INSERT")
	}
	table := strings.TrimSpace(rest[:parenIdx])

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
	valuesIdx := caseInsensitiveIndex(rest, "VALUES ")
	if valuesIdx < 0 {
		return "", nil, fmt.Errorf("missing VALUES clause")
	}
	valuesRest := rest[valuesIdx+7:]
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

func parseUpdateQuery(query string, args []interface{}) (string, map[string]interface{}, map[string]interface{}, error) {
	q := strings.TrimSpace(query)
	q = regexp.MustCompile(`\s+`).ReplaceAllString(q, " ")

	upper := strings.ToUpper(q)
	if !strings.HasPrefix(upper, "UPDATE ") {
		return "", nil, nil, fmt.Errorf("not an UPDATE query")
	}

	setIdx := caseInsensitiveIndex(q, " SET ")
	if setIdx < 0 {
		return "", nil, nil, fmt.Errorf("missing SET clause")
	}
	table := strings.TrimSpace(q[7:setIdx])

	rest := q[setIdx+5:]
	whereIdx := caseInsensitiveIndex(rest, " WHERE ")
	var setClause, whereClause string
	if whereIdx >= 0 {
		setClause = strings.TrimSpace(rest[:whereIdx])
		whereRest := rest[whereIdx+7:]
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

func parseDeleteQuery(query string, args []interface{}) (string, map[string]interface{}, error) {
	q := strings.TrimSpace(query)
	q = regexp.MustCompile(`\s+`).ReplaceAllString(q, " ")

	upper := strings.ToUpper(q)
	if !strings.HasPrefix(upper, "DELETE FROM ") {
		return "", nil, fmt.Errorf("not a DELETE query")
	}

	rest := q[12:]
	whereIdx := caseInsensitiveIndex(rest, " WHERE ")
	if whereIdx < 0 {
		table := strings.TrimSpace(rest)
		return table, map[string]interface{}{}, nil
	}
	table := strings.TrimSpace(rest[:whereIdx])
	whereClause := strings.TrimSpace(rest[whereIdx+7:])

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

func (c *Client) ExecuteQueryRow(query string, dest interface{}, args ...interface{}) error {
	parts, err := parseSelectQuery(query, args)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	params := url.Values{}
	params.Set("select", parts.columns)

	for _, f := range parts.filters {
		params.Add(f.column, fmt.Sprintf("%s.%v", f.operator, f.value))
	}

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

func (c *Client) ExecuteQuery(query string, dest interface{}, args ...interface{}) error {
	parts, err := parseSelectQuery(query, args)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}
	params := url.Values{}
	params.Set("select", parts.columns)

	for _, f := range parts.filters {
		params.Add(f.column, fmt.Sprintf("%s.%v", f.operator, f.value))
	}

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

type Client struct {
	BaseURL    string       
	APIKey     string     
	HTTPClient *http.Client 
}

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

func (c *Client) TableURL(table string) string {
	return fmt.Sprintf("%s/rest/v1/%s", c.BaseURL, table)
}


func (c *Client) StorageURL() string {
	return fmt.Sprintf("%s/storage/v1", c.BaseURL)
}

func (c *Client) SetHeaders(req *http.Request) {
	req.Header.Set("apikey", c.APIKey)
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
