package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/repositories"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load Configurations
	cfg := config.LoadConfig()

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// 2. Initialize Database Connection & Auto-migrate
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 3. Initialize Repositories
	userRepo := repositories.NewUserRepository(db)
	urlRepo := repositories.NewURLRepository(db)
	clickRepo := repositories.NewClickRepository(db)

	// 4. Initialize Services
	authService := services.NewAuthService(userRepo, cfg)
	urlService := services.NewURLService(urlRepo, clickRepo)
	analyticsService := services.NewAnalyticsService(urlRepo, clickRepo)

	// 5. Initialize Handlers
	authHandler := handlers.NewAuthHandler(authService)
	urlHandler := handlers.NewURLHandler(urlService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	redirectHandler := handlers.NewRedirectHandler(urlService, cfg)

	// 6. Setup Router
	r := gin.Default()

	// Global Middlewares
	r.Use(middleware.SetupCORS())

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "urlshortener-backend",
		})
	})

	// API Routes
	api := r.Group("/api/v1")
	{
		// Public Auth
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Public URL verification / lookup
		pub := api.Group("/public")
		{
			pub.POST("/verify-password", redirectHandler.VerifyPassword)
			pub.GET("/link-info/:code", redirectHandler.GetPublicLinkInfo)
		}

		// Protected Routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// Auth
			protected.GET("/auth/me", authHandler.GetMe)

			// URL Management
			urls := protected.Group("/urls")
			{
				urls.GET("", urlHandler.GetAll)
				urls.POST("", urlHandler.Create)
				urls.GET("/:id", urlHandler.GetByID)
				urls.PUT("/:id", urlHandler.Update)
				urls.DELETE("/:id", urlHandler.Delete)
				urls.PATCH("/:id/toggle", urlHandler.ToggleStatus)
			}

			// Analytics
			analytics := protected.Group("/analytics")
			{
				analytics.GET("/overview", analyticsHandler.GetOverview)
				analytics.GET("/urls/:id", analyticsHandler.GetURLAnalytics)
				analytics.GET("/urls/:id/export", analyticsHandler.ExportCSV)
			}
		}
	}

	// Public Redirection Route (e.g. /r/:code or /:code)
	r.GET("/r/:code", redirectHandler.HandleRedirect)

	// Reserved paths for frontend and system routes that shouldn't be treated as shortcodes
	reservedPaths := map[string]bool{
		"":                   true,
		"api":                true,
		"health":             true,
		"assets":             true,
		"favicon.ico":        true,
		"login":              true,
		"register":           true,
		"dashboard":          true,
		"links":              true,
		"analytics":          true,
		"p":                  true,
		"link-not-found":     true,
		"link-inactive":      true,
		"link-expired":       true,
		"link-limit-reached": true,
		"link-error":         true,
	}

	r.GET("/:code", func(c *gin.Context) {
		code := c.Param("code")
		if reservedPaths[code] {
			c.Next()
			return
		}
		redirectHandler.HandleRedirect(c)
	})

	// 7. Serve Frontend SPA if dist directory exists
	distPath := "./dist"
	if _, err := os.Stat(distPath); err == nil {
		r.Static("/assets", distPath+"/assets")
		r.StaticFile("/favicon.ico", distPath+"/favicon.ico")

		// SPA Fallback for non-API routes
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/api") {
				c.JSON(http.StatusNotFound, gin.H{"error": "API route not found"})
				return
			}
			c.File(distPath + "/index.html")
		})
	}

	log.Printf("Server listening on port %s...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
