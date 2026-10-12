package router

import (
	"net/http"

	"live-coding-starter/backend/internal/config"
	"live-coding-starter/backend/internal/handler"
	"live-coding-starter/backend/internal/middleware"
	"live-coding-starter/backend/internal/repository"
	"live-coding-starter/backend/internal/service"

	"github.com/gin-gonic/gin"
)

func Setup(
	cfg config.Config,
	productHandler *handler.ProductHandler,
	authHandler *handler.AuthHandler,
	adminHandler *handler.AdminHandler,
	tokenService service.TokenService,
	userRepository repository.UserRepository,
) *gin.Engine {
	router := gin.Default()

	router.Use(corsMiddleware(cfg.AllowedOrigin))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1")
	authMiddleware := middleware.Auth(tokenService, userRepository)

	auth := api.Group("/auth")
	authHandler.RegisterRoutes(auth, authMiddleware)

	products := api.Group("/products")
	productHandler.RegisterRoutes(products)

	admin := api.Group("/admin")
	admin.Use(authMiddleware, middleware.AdminOnly())
	adminHandler.RegisterRoutes(admin)

	return router
}

func corsMiddleware(allowedOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == allowedOrigin {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Vary", "Origin")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
