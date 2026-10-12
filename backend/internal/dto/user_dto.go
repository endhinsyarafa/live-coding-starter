package dto

type CreateUserRequest struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password"`
	Role            string `json:"role" binding:"required"`
	IsActive        *bool  `json:"is_active"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required"`
	IsActive bool   `json:"is_active"`
}

type UpdateUserStatusRequest struct {
	IsActive bool `json:"is_active"`
}

type DashboardStatsResponse struct {
	TotalUsers    int64 `json:"total_users"`
	ActiveUsers   int64 `json:"active_users"`
	InactiveUsers int64 `json:"inactive_users"`
	NewUsers      int64 `json:"new_users"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}
