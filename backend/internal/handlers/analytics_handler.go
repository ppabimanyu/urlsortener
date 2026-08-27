package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"backend/internal/middleware"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AnalyticsHandler struct {
	analyticsService services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

func (h *AnalyticsHandler) GetOverview(c *gin.Context) {
	userIDVal, _ := c.Get(middleware.ContextUserIDKey)
	userID := userIDVal.(uuid.UUID)

	overview, err := h.analyticsService.GetOverview(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": overview})
}

func (h *AnalyticsHandler) GetURLAnalytics(c *gin.Context) {
	userIDVal, _ := c.Get(middleware.ContextUserIDKey)
	userID := userIDVal.(uuid.UUID)

	urlID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL ID"})
		return
	}

	days, _ := strconv.Atoi(c.DefaultQuery("days", "14"))

	stats, err := h.analyticsService.GetURLAnalytics(urlID, userID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *AnalyticsHandler) ExportCSV(c *gin.Context) {
	userIDVal, _ := c.Get(middleware.ContextUserIDKey)
	userID := userIDVal.(uuid.UUID)

	urlID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL ID"})
		return
	}

	csvBytes, filename, err := h.analyticsService.ExportURLAnalyticsCSV(urlID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.Data(http.StatusOK, "text/csv", csvBytes)
}
