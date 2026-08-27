package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/services"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type RedirectHandler struct {
	urlService services.URLService
	cfg        *config.Config
}

func NewRedirectHandler(urlService services.URLService, cfg *config.Config) *RedirectHandler {
	return &RedirectHandler{
		urlService: urlService,
		cfg:        cfg,
	}
}

type VerifyPasswordRequest struct {
	Code     string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *RedirectHandler) HandleRedirect(c *gin.Context) {
	code := strings.TrimSpace(c.Param("code"))
	if code == "" {
		c.Redirect(http.StatusFound, h.cfg.FrontendURL)
		return
	}

	urlItem, isPasswordRequired, err := h.urlService.ResolveURL(code)
	if err != nil {
		switch err {
		case services.ErrURLNotFound:
			c.Redirect(http.StatusFound, fmt.Sprintf("%s/link-not-found", h.cfg.FrontendURL))
		case services.ErrURLInactive:
			c.Redirect(http.StatusFound, fmt.Sprintf("%s/link-inactive?slug=%s", h.cfg.FrontendURL, code))
		case services.ErrURLExpired:
			c.Redirect(http.StatusFound, fmt.Sprintf("%s/link-expired?slug=%s", h.cfg.FrontendURL, code))
		case services.ErrClickLimitReached:
			c.Redirect(http.StatusFound, fmt.Sprintf("%s/link-limit-reached?slug=%s", h.cfg.FrontendURL, code))
		default:
			c.Redirect(http.StatusFound, fmt.Sprintf("%s/link-error", h.cfg.FrontendURL))
		}
		return
	}

	// If password protected, redirect user to frontend password prompt page
	if isPasswordRequired {
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/p/%s", h.cfg.FrontendURL, code))
		return
	}

	// Record click analytics asynchronously
	h.recordClickData(c, urlItem)

	// Perform 302 temporary redirect to destination
	c.Redirect(http.StatusFound, urlItem.OriginalURL)
}

func (h *RedirectHandler) VerifyPassword(c *gin.Context) {
	var req VerifyPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	urlItem, err := h.urlService.VerifyURLPassword(req.Code, req.Password)
	if err != nil {
		switch err {
		case services.ErrURLNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case services.ErrInvalidURLPassword:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Kata sandi salah"})
		case services.ErrURLInactive:
			c.JSON(http.StatusForbidden, gin.H{"error": "Tautan ini telah dinonaktifkan"})
		case services.ErrURLExpired:
			c.JSON(http.StatusGone, gin.H{"error": "Tautan ini telah kedaluwarsa"})
		case services.ErrClickLimitReached:
			c.JSON(http.StatusGone, gin.H{"error": "Tautan ini telah mencapai batas klik maksimal"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Record click on successful verification
	h.recordClickData(c, urlItem)

	c.JSON(http.StatusOK, gin.H{
		"original_url": urlItem.OriginalURL,
		"title":        urlItem.Title,
	})
}

func (h *RedirectHandler) GetPublicLinkInfo(c *gin.Context) {
	code := strings.TrimSpace(c.Param("code"))
	urlItem, isPasswordRequired, err := h.urlService.ResolveURL(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_code":           urlItem.ShortCode,
		"title":                urlItem.Title,
		"is_password_required": isPasswordRequired,
	})
}

func (h *RedirectHandler) recordClickData(c *gin.Context, urlItem *models.URL) {
	uaString := c.Request.UserAgent()
	deviceType, osName, browser := utils.ParseUserAgent(uaString)
	referrer := utils.CleanReferrer(c.Request.Referer())
	clientIP := utils.ExtractClientIP(c.Request)
	country, countryCode, city := utils.LookupGeoLocation(clientIP)

	click := &models.URLClick{
		ClickedAt:   time.Now(),
		IPAddress:   clientIP,
		UserAgent:   uaString,
		DeviceType:  deviceType,
		OS:          osName,
		Browser:     browser,
		Referrer:    referrer,
		Country:     country,
		CountryCode: countryCode,
		City:        city,
	}

	h.urlService.RecordClick(urlItem.ID, click)
}
