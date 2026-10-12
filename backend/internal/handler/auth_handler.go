package handler

import (
	"errors"
	"net/http"

	"live-coding-starter/backend/internal/dto"
	"live-coding-starter/backend/internal/middleware"
	"live-coding-starter/backend/internal/service"
	"live-coding-starter/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	router.POST("/login", h.Login)
	router.POST("/logout", h.Logout)
	router.GET("/me", authMiddleware, h.Me)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	result, err := h.authService.Login(req)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Login berhasil", result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	response.Success(c, http.StatusOK, "Logout berhasil", nil)
}

func (h *AuthHandler) Me(c *gin.Context) {
	currentUser, ok := middleware.CurrentUser(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	user, err := h.authService.GetCurrentUser(currentUser.ID)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Authenticated user retrieved successfully", user)
}

func handleAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		response.Error(c, http.StatusUnauthorized, "Email atau password salah.", nil)
	case errors.Is(err, service.ErrInactiveAccount):
		response.Error(c, http.StatusForbidden, "Akun Anda telah dinonaktifkan. Silakan hubungi administrator.", nil)
	case errors.Is(err, service.ErrUnauthorized):
		response.Error(c, http.StatusUnauthorized, "Authentication required", nil)
	default:
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
	}
}
