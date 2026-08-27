package handlers

import (
	"net/http"
	"strconv"
	"time"

	"backend/internal/middleware"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type URLHandler struct {
	urlService services.URLService
}

func NewURLHandler(urlService services.URLService) *URLHandler {
	return &URLHandler{urlService: urlService}
}

type CreateURLRequest struct {
	OriginalURL string     `json:"original_url" binding:"required"`
	CustomSlug  string     `json:"custom_slug"`
	Title       string     `json:"title"`
	Password    string     `json:"password"`
	ExpiresAt   *time.Time `json:"expires_at"`
	ClickLimit  *int       `json:"click_limit"`
}

type UpdateURLRequest struct {
	OriginalURL string     `json:"original_url"`
	CustomSlug  string     `json:"custom_slug"`
	Title       string     `json:"title"`
	Password    *string    `json:"password"`
	ExpiresAt   *time.Time `json:"expires_at"`
	ClickLimit  *int       `json:"click_limit"`
	IsActive    *bool      `json:"is_active"`
}

func (h *URLHandler) Create(c *gin.Context) {
	userIDVal, _ := c.Get(middleware.ContextUserIDKey)
	userID := userIDVal.(uuid.UUID)

	var req CreateURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	urlItem, err := h.urlService.CreateURL(
		userID,
		req.OriginalURL,
		req.CustomSlug,
		req.Title,
		req.Password,
		req.ExpiresAt,
		req.ClickLimit,
	)
	if err != nil {
		if err == services.ErrSlugAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err == services.ErrInvalidURL || err == services.ErrInvalidSlug {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "URL shortened successfully",
		"data":    urlItem,
	})
}

func (h *URLHandler) GetAll(c *gin.Context) {
	userIDVal, _ := c.Get(middleware.ContextUserIDKey)
	userID := userIDVal.(uuid.UUID)

	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var activeFilter *bool
	if activeStr := c.Query("is_active"); activeStr != "" {
		if val, err := strconv.ParseBool(activeStr); err == nil {
			activeFilter = &val
		}
	}

	urls, total, err := h.urlService.GetURLs(userID, search, activeFilter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  urls,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *URLHandler) GetByID(c *gin.Context) {
	userIDVal, _ := c.Get(middleware.ContextUserIDKey)
	userID := userIDVal.(uuid.UUID)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL ID"})
		return
	}

	urlItem, err := h.urlService.GetURLByID(id, userID)
	if err != nil {
		if err == services.ErrURLNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": urlItem})
}

func (h *URLHandler) Update(c *gin.Context) {
	userIDVal, _ := c.Get(middleware.ContextUserIDKey)
	userID := userIDVal.(uuid.UUID)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL ID"})
		return
	}

	var req UpdateURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedURL, err := h.urlService.UpdateURL(
		id,
		userID,
		req.OriginalURL,
		req.CustomSlug,
		req.Title,
		req.Password,
		req.ExpiresAt,
		req.ClickLimit,
		req.IsActive,
	)
	if err != nil {
		if err == services.ErrURLNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == services.ErrSlugAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err == services.ErrInvalidURL || err == services.ErrInvalidSlug {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "URL updated successfully",
		"data":    updatedURL,
	})
}

func (h *URLHandler) Delete(c *gin.Context) {
	userIDVal, _ := c.Get(middleware.ContextUserIDKey)
	userID := userIDVal.(uuid.UUID)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL ID"})
		return
	}

	if err := h.urlService.DeleteURL(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL deleted successfully"})
}

func (h *URLHandler) ToggleStatus(c *gin.Context) {
	userIDVal, _ := c.Get(middleware.ContextUserIDKey)
	userID := userIDVal.(uuid.UUID)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL ID"})
		return
	}

	urlItem, err := h.urlService.ToggleStatus(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "URL status toggled",
		"data":    urlItem,
	})
}
