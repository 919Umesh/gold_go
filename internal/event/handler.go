package event

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetCompanyEvents godoc
// @Summary Get events for a company
// @Description Get list of corporate events (AGM, dividend, etc.) for a specific company
// @Tags events
// @Produce json
// @Param company_id path string true "Company ID"
// @Param limit query int false "Limit (default 50)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /market/companies/{company_id}/events [get]
func (h *Handler) GetCompanyEvents(c *gin.Context) {
	companyID := c.Param("company_id")
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	events, err := h.service.GetEventsByCompany(companyID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events": events,
		"count":  len(events),
	})
}

// GetUpcomingEvents godoc
// @Summary Get upcoming events
// @Description Get list of all upcoming corporate events across all companies
// @Tags events
// @Produce json
// @Param limit query int false "Limit (default 50)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /market/events/upcoming [get]
func (h *Handler) GetUpcomingEvents(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	events, err := h.service.GetUpcomingEvents(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events": events,
		"count":  len(events),
	})
}

// GetEventsByType godoc
// @Summary Get events by type
// @Description Get events filtered by event type (agm, dividend, bonus_share, etc.)
// @Tags events
// @Produce json
// @Param type path string true "Event type"
// @Param limit query int false "Limit (default 50)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /market/events/type/{type} [get]
func (h *Handler) GetEventsByType(c *gin.Context) {
	eventType := c.Param("type")
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	events, err := h.service.GetEventsByType(eventType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events": events,
		"count":  len(events),
	})
}

// GetAllEvents godoc
// @Summary Get all events
// @Description Get list of all corporate events with pagination
// @Tags events
// @Produce json
// @Param limit query int false "Limit (default 50)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /market/events [get]
func (h *Handler) GetAllEvents(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	events, err := h.service.GetAllEvents(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events": events,
		"count":  len(events),
	})
}
