package repository

import (
	"time"

	"live-coding-starter/backend/internal/model"

	"gorm.io/gorm"
)

type UserListFilter struct {
	Search string
	Status string
	Role   string
	Page   int
	Limit  int
}

type DashboardStats struct {
	TotalUsers    int64
	ActiveUsers   int64
	InactiveUsers int64
	NewUsers      int64
}

type UserRepository interface {
	FindByID(id uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	List(filter UserListFilter) ([]model.User, int64, error)
	GetDashboardStats() (DashboardStats, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("LOWER(email) = LOWER(?)", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) List(filter UserListFilter) ([]model.User, int64, error) {
	query := r.db.Model(&model.User{})
	query = applyUserListFilter(query, filter)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []model.User
	offset := (filter.Page - 1) * filter.Limit
	err := query.Order("created_at DESC").Limit(filter.Limit).Offset(offset).Find(&users).Error

	return users, total, err
}

func (r *userRepository) GetDashboardStats() (DashboardStats, error) {
	var stats DashboardStats
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	if err := r.db.Model(&model.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return stats, err
	}
	if err := r.db.Model(&model.User{}).Where("is_active = ?", true).Count(&stats.ActiveUsers).Error; err != nil {
		return stats, err
	}
	if err := r.db.Model(&model.User{}).Where("is_active = ?", false).Count(&stats.InactiveUsers).Error; err != nil {
		return stats, err
	}
	if err := r.db.Model(&model.User{}).Where("created_at >= ?", thirtyDaysAgo).Count(&stats.NewUsers).Error; err != nil {
		return stats, err
	}

	return stats, nil
}

func applyUserListFilter(query *gorm.DB, filter UserListFilter) *gorm.DB {
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", search, search)
	}

	if filter.Status == "active" {
		query = query.Where("is_active = ?", true)
	}
	if filter.Status == "inactive" {
		query = query.Where("is_active = ?", false)
	}

	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}

	return query
}
