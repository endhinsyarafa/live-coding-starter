package service

import (
	"errors"
	"math"
	"strings"
	"time"

	"live-coding-starter/backend/internal/dto"
	"live-coding-starter/backend/internal/model"
	"live-coding-starter/backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrInactiveAccount      = errors.New("inactive account")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrForbidden            = errors.New("forbidden")
	ErrUserNotFound         = errors.New("user not found")
	ErrDuplicateEmail       = errors.New("email already exists")
	ErrInvalidRole          = errors.New("invalid role")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrCannotDeactivateSelf = errors.New("cannot deactivate own account")
	ErrInvalidPagination    = errors.New("invalid pagination")
)

type AuthService interface {
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	GetCurrentUser(userID uint) (*dto.UserResponse, error)
	BuildUserResponse(user model.User) dto.UserResponse
}

type UserService interface {
	GetDashboardStats() (*dto.DashboardStatsResponse, error)
	ListUsers(filter repository.UserListFilter) ([]dto.UserResponse, dto.PaginationMeta, error)
	GetUserByID(id uint) (*dto.UserResponse, error)
	CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error)
	UpdateUser(id uint, currentUserID uint, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	UpdateUserStatus(id uint, currentUserID uint, req dto.UpdateUserStatusRequest) (*dto.UserResponse, error)
}

type authService struct {
	userRepository repository.UserRepository
	tokenService   TokenService
}

type userService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository, tokenService TokenService) AuthService {
	return &authService{
		userRepository: userRepository,
		tokenService:   tokenService,
	}
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrInactiveAccount
	}

	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepository.Update(user); err != nil {
		return nil, err
	}

	token, err := s.tokenService.Generate(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		User:  s.BuildUserResponse(*user),
		Token: token,
	}, nil
}

func (s *authService) GetCurrentUser(userID uint) (*dto.UserResponse, error) {
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUnauthorized
		}

		return nil, err
	}

	if !user.IsActive {
		return nil, ErrInactiveAccount
	}

	response := s.BuildUserResponse(*user)
	return &response, nil
}

func (s *authService) BuildUserResponse(user model.User) dto.UserResponse {
	return buildUserResponse(user)
}

func (s *userService) GetDashboardStats() (*dto.DashboardStatsResponse, error) {
	stats, err := s.userRepository.GetDashboardStats()
	if err != nil {
		return nil, err
	}

	return &dto.DashboardStatsResponse{
		TotalUsers:    stats.TotalUsers,
		ActiveUsers:   stats.ActiveUsers,
		InactiveUsers: stats.InactiveUsers,
		NewUsers:      stats.NewUsers,
	}, nil
}

func (s *userService) ListUsers(filter repository.UserListFilter) ([]dto.UserResponse, dto.PaginationMeta, error) {
	filter = normalizeUserListFilter(filter)
	if filter.Page <= 0 || filter.Limit <= 0 {
		return nil, dto.PaginationMeta{}, ErrInvalidPagination
	}

	users, total, err := s.userRepository.List(filter)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	responses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, buildUserResponse(user))
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(filter.Limit)))
	}

	meta := dto.PaginationMeta{
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return responses, meta, nil
}

func (s *userService) GetUserByID(id uint) (*dto.UserResponse, error) {
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	response := buildUserResponse(*user)
	return &response, nil
}

func (s *userService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	role := strings.TrimSpace(strings.ToLower(req.Role))
	if !isValidRole(role) {
		return nil, ErrInvalidRole
	}

	if len(req.Password) < 8 || (req.ConfirmPassword != "" && req.Password != req.ConfirmPassword) {
		return nil, ErrInvalidPassword
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	if exists, err := s.emailExists(email, 0); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrDuplicateEmail
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	user := &model.User{
		Name:     strings.TrimSpace(req.Name),
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
		IsActive: isActive,
	}

	if err := s.userRepository.Create(user); err != nil {
		return nil, err
	}

	response := buildUserResponse(*user)
	return &response, nil
}

func (s *userService) UpdateUser(id uint, currentUserID uint, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	role := strings.TrimSpace(strings.ToLower(req.Role))
	if !isValidRole(role) {
		return nil, ErrInvalidRole
	}

	user, err := s.userRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	if id == currentUserID && user.IsActive && !req.IsActive {
		return nil, ErrCannotDeactivateSelf
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	if exists, err := s.emailExists(email, id); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrDuplicateEmail
	}

	user.Name = strings.TrimSpace(req.Name)
	user.Email = email
	user.Role = role
	user.IsActive = req.IsActive

	if err := s.userRepository.Update(user); err != nil {
		return nil, err
	}

	response := buildUserResponse(*user)
	return &response, nil
}

func (s *userService) UpdateUserStatus(id uint, currentUserID uint, req dto.UpdateUserStatusRequest) (*dto.UserResponse, error) {
	if id == currentUserID && !req.IsActive {
		return nil, ErrCannotDeactivateSelf
	}

	user, err := s.userRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	user.IsActive = req.IsActive
	if err := s.userRepository.Update(user); err != nil {
		return nil, err
	}

	response := buildUserResponse(*user)
	return &response, nil
}

func (s *userService) emailExists(email string, exceptUserID uint) (bool, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}

	return user.ID != exceptUserID, nil
}

func normalizeUserListFilter(filter repository.UserListFilter) repository.UserListFilter {
	filter.Search = strings.TrimSpace(filter.Search)
	filter.Status = strings.TrimSpace(strings.ToLower(filter.Status))
	filter.Role = strings.TrimSpace(strings.ToLower(filter.Role))

	if filter.Status != "active" && filter.Status != "inactive" {
		filter.Status = ""
	}
	if !isValidRole(filter.Role) {
		filter.Role = ""
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	return filter
}

func isValidRole(role string) bool {
	return role == model.RoleAdmin || role == model.RoleUser
}

func buildUserResponse(user model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Role:        user.Role,
		IsActive:    user.IsActive,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
