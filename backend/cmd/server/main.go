package main

import (
	"log"

	"live-coding-starter/backend/internal/config"
	"live-coding-starter/backend/internal/database"
	"live-coding-starter/backend/internal/handler"
	"live-coding-starter/backend/internal/repository"
	"live-coding-starter/backend/internal/router"
	"live-coding-starter/backend/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	log.Printf(
		"database config loaded: host=%s port=%s user=%s database=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	log.Println("database connected successfully")

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("failed to run auto migrate: %v", err)
	}
	log.Println("migration completed successfully")

	if cfg.SeedData {
		if err := database.SeedProducts(db); err != nil {
			log.Fatalf("failed to seed products: %v", err)
		}
		log.Println("seed data completed successfully")
	}

	if cfg.SeedAdmin {
		if err := database.SeedAdmin(db, cfg.AdminSeed); err != nil {
			log.Fatalf("failed to seed admin: %v", err)
		}
		log.Println("admin seed completed successfully")
	}

	productRepository := repository.NewProductRepository(db)
	userRepository := repository.NewUserRepository(db)
	tokenService := service.NewTokenService(cfg.JWT)
	productService := service.NewProductService(productRepository)
	authService := service.NewAuthService(userRepository, tokenService)
	userService := service.NewUserService(userRepository)
	productHandler := handler.NewProductHandler(productService)
	authHandler := handler.NewAuthHandler(authService)
	adminHandler := handler.NewAdminHandler(userService)

	app := router.Setup(cfg, productHandler, authHandler, adminHandler, tokenService, userRepository)

	log.Printf("server running on http://localhost:%s", cfg.HTTPPort)
	if err := app.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
