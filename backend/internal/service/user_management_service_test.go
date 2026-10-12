package service

import (
	"errors"
	"testing"

	"live-coding-starter/backend/internal/config"
	"live-coding-starter/backend/internal/dto"
	"live-coding-starter/backend/internal/model"
	"live-coding-starter/backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type fakeUserRepository struct {
	users  map[uint]*model.User
	nextID uint
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		users:  map[uint]*model.User{},
		nextID: 1,
	}
}

func (r *fakeUserRepository) FindByID(id uint) (*model.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}

	copiedUser := *user
	return &copiedUser, nil
}

func (r *fakeUserRepository) FindByEmail(email string) (*model.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			copiedUser := *user
			return &copiedUser, nil
		}
	}

	return nil, gorm.ErrRecordNotFound
}

func (r *fakeUserRepository) Create(user *model.User) error {
	user.ID = r.nextID
	r.nextID++
	copiedUser := *user
	r.users[user.ID] = &copiedUser

	return nil
}

func (r *fakeUserRepository) Update(user *model.User) error {
	copiedUser := *user
	r.users[user.ID] = &copiedUser

	return nil
}

func (r *fakeUserRepository) List(filter repository.UserListFilter) ([]model.User, int64, error) {
	users := make([]model.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, *user)
	}

	return users, int64(len(users)), nil
}

func (r *fakeUserRepository) GetDashboardStats() (repository.DashboardStats, error) {
	return repository.DashboardStats{}, nil
}

func TestAuthLoginActiveUser(t *testing.T) {
	repo := newFakeUserRepository()
	user := createFakeUser(t, repo, "admin@example.com", "Password123!", true, model.RoleAdmin)
	authService := NewAuthService(repo, newTestTokenService())

	result, err := authService.Login(dto.LoginRequest{
		Email:    user.Email,
		Password: "Password123!",
	})

	if err != nil {
		t.Fatalf("expected login success, got %v", err)
	}
	if result.Token == "" {
		t.Fatal("expected token to be generated")
	}
	if result.User.Email != user.Email {
		t.Fatalf("expected user email %q, got %q", user.Email, result.User.Email)
	}
}

func TestAuthLoginRejectsWrongPassword(t *testing.T) {
	repo := newFakeUserRepository()
	createFakeUser(t, repo, "admin@example.com", "Password123!", true, model.RoleAdmin)
	authService := NewAuthService(repo, newTestTokenService())

	_, err := authService.Login(dto.LoginRequest{
		Email:    "admin@example.com",
		Password: "wrong-password",
	})

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthLoginRejectsInactiveUser(t *testing.T) {
	repo := newFakeUserRepository()
	createFakeUser(t, repo, "user@example.com", "Password123!", false, model.RoleUser)
	authService := NewAuthService(repo, newTestTokenService())

	_, err := authService.Login(dto.LoginRequest{
		Email:    "user@example.com",
		Password: "Password123!",
	})

	if !errors.Is(err, ErrInactiveAccount) {
		t.Fatalf("expected ErrInactiveAccount, got %v", err)
	}
}

func TestUpdateUserStatusRejectsSelfDeactivate(t *testing.T) {
	repo := newFakeUserRepository()
	user := createFakeUser(t, repo, "admin@example.com", "Password123!", true, model.RoleAdmin)
	userService := NewUserService(repo)

	_, err := userService.UpdateUserStatus(user.ID, user.ID, dto.UpdateUserStatusRequest{
		IsActive: false,
	})

	if !errors.Is(err, ErrCannotDeactivateSelf) {
		t.Fatalf("expected ErrCannotDeactivateSelf, got %v", err)
	}
}

func createFakeUser(t *testing.T, repo *fakeUserRepository, email string, password string, active bool, role string) *model.User {
	t.Helper()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	user := &model.User{
		Name:     "Test User",
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
		IsActive: active,
	}
	if err := repo.Create(user); err != nil {
		t.Fatalf("failed to create fake user: %v", err)
	}

	return user
}

func newTestTokenService() TokenService {
	return NewTokenService(config.JWTConfig{
		Secret:          "test-secret",
		ExpirationHours: 1,
	})
}
