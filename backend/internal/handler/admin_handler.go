package handler

import (
	"errors"
	"net/http"
	"strconv"

	"live-coding-starter/backend/internal/dto"
	"live-coding-starter/backend/internal/middleware"
	"live-coding-starter/backend/internal/repository"
	"live-coding-starter/backend/internal/service"
	"live-coding-starter/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	userService service.UserService
}

func NewAdminHandler(userService service.UserService) *AdminHandler {
	return &AdminHandler{userService: userService}
}

func (h *AdminHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/dashboard", h.GetDashboard)
	router.GET("/users", h.GetUsers)
	router.GET("/users/:id", h.GetUser)
	router.POST("/users", h.CreateUser)
	router.PATCH("/users/:id", h.UpdateUser)
	router.PATCH("/users/:id/status", h.UpdateUserStatus)
}

func (h *AdminHandler) GetDashboard(c *gin.Context) {
	stats, err := h.userService.GetDashboardStats()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve dashboard", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Dashboard retrieved successfully", stats)
}

func (h *AdminHandler) GetUsers(c *gin.Context) {
	page := parseQueryInt(c, "page", 1)
	limit := parseQueryInt(c, "limit", 10)

	users, meta, err := h.userService.ListUsers(repository.UserListFilter{
		Search: c.Query("search"),
		Status: c.Query("status"),
		Role:   c.Query("role"),
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		handleAdminError(c, err)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "Users retrieved successfully", users, meta)
}

func (h *AdminHandler) GetUser(c *gin.Context) {
	id, ok := parseUserID(c)
	if !ok {
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		handleAdminError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "User retrieved successfully", user)
}

func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		handleAdminError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "User berhasil ditambahkan.", user)
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	id, ok := parseUserID(c)
	if !ok {
		return
	}

	currentUser, ok := middleware.CurrentUser(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.userService.UpdateUser(id, currentUser.ID, req)
	if err != nil {
		handleAdminError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Data user berhasil diperbarui.", user)
}

func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	id, ok := parseUserID(c)
	if !ok {
		return
	}

	currentUser, ok := middleware.CurrentUser(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	var req dto.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.userService.UpdateUserStatus(id, currentUser.ID, req)
	if err != nil {
		handleAdminError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Status user berhasil diperbarui", user)
}

func parseQueryInt(c *gin.Context, key string, fallback int) int {
	value, err := strconv.Atoi(c.DefaultQuery(key, strconv.Itoa(fallback)))
	if err != nil {
		return fallback
	}

	return value
}

func parseUserID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, "Invalid user ID", nil)
		return 0, false
	}

	return uint(id), true
}

func handleAdminError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		response.Error(c, http.StatusNotFound, "User not found", nil)
	case errors.Is(err, service.ErrDuplicateEmail):
		response.Error(c, http.StatusConflict, "Email sudah digunakan.", nil)
	case errors.Is(err, service.ErrInvalidRole):
		response.Error(c, http.StatusBadRequest, "Role tidak valid.", nil)
	case errors.Is(err, service.ErrInvalidPassword):
		response.Error(c, http.StatusBadRequest, "Password tidak valid.", nil)
	case errors.Is(err, service.ErrCannotDeactivateSelf):
		response.Error(c, http.StatusBadRequest, "Anda tidak dapat menonaktifkan akun sendiri.", nil)
	case errors.Is(err, service.ErrInvalidPagination):
		response.Error(c, http.StatusBadRequest, "Pagination tidak valid.", nil)
	default:
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
	}
}
