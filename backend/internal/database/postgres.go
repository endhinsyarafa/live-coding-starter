package database

import (
	"fmt"
	"strings"

	"live-coding-starter/backend/internal/config"
	"live-coding-starter/backend/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(buildDSN(cfg)), &gorm.Config{})
}

func buildDSN(cfg config.DatabaseConfig) string {
	parts := []string{
		fmt.Sprintf("host=%s", cfg.Host),
		fmt.Sprintf("user=%s", cfg.User),
		fmt.Sprintf("dbname=%s", cfg.Name),
		fmt.Sprintf("port=%s", cfg.Port),
		fmt.Sprintf("sslmode=%s", cfg.SSLMode),
		"TimeZone=Asia/Jakarta",
	}

	if strings.TrimSpace(cfg.Password) != "" {
		parts = append(parts, fmt.Sprintf("password=%s", cfg.Password))
	}

	return strings.Join(parts, " ")
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.Product{}, &model.User{})
}

func SeedProducts(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.Product{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	products := []model.Product{
		{
			Name:        "Keyboard Mechanical",
			Description: "Mechanical keyboard untuk programmer",
			Price:       750000,
			Stock:       10,
		},
		{
			Name:        "Mouse Wireless",
			Description: "Mouse wireless ergonomis untuk kerja harian",
			Price:       250000,
			Stock:       20,
		},
		{
			Name:        "Monitor 24 Inch",
			Description: "Monitor 24 inch untuk setup produktif",
			Price:       1850000,
			Stock:       5,
		},
	}

	return db.Create(&products).Error
}

func SeedAdmin(db *gorm.DB, cfg config.AdminSeedConfig) error {
	var count int64
	if err := db.Model(&model.User{}).Where("LOWER(email) = LOWER(?)", cfg.Email).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := model.User{
		Name:     cfg.Name,
		Email:    strings.ToLower(strings.TrimSpace(cfg.Email)),
		Password: string(hashedPassword),
		Role:     model.RoleAdmin,
		IsActive: true,
	}

	return db.Create(&admin).Error
}
