package model

import "time"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"not null"`
	Email       string     `json:"email" gorm:"uniqueIndex;not null"`
	Password    string     `json:"-" gorm:"not null"`
	Role        string     `json:"role" gorm:"not null;default:user"`
	IsActive    bool       `json:"is_active" gorm:"not null;default:true"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
