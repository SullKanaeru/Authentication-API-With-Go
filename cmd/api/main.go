package main

import (
	"log"
	"authentication_api/internal/config"
	"authentication_api/internal/handler"
	"authentication_api/internal/repository"
	"authentication_api/internal/service"
	"authentication_api/internal/middleware"
	"authentication_api/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.LoadConfig()

	db := database.ConnectDB()
	database.MigrateDB(db)

	userRepo := repository.NewUserRepository(db)

	notifService := service.NewNotificationService(
		config.GetEnv("BREVO_API_KEY", ""),
		config.GetEnv("SENDER_EMAIL", ""),
		config.GetEnv("SENDER_NAME", "Aplikasi Saya"),
		config.GetEnv("FONNTE_TOKEN", ""),
	)
	
	authService := service.NewAuthService(userRepo, config.GetEnv("JWT_SECRET", "default_secret"), notifService)
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userRepo)

	app := fiber.New()
	app.Use(logger.New()) 

	app.Use(cors.New(cors.Config{
        AllowOrigins:     "http://localhost:5173, https://domain-frontend-kamu.com",
        AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
        AllowCredentials: true,
    }))

	api := app.Group("/api/v1")
	auth := api.Group("/auth")
	
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)
	auth.Post("/verify", authHandler.VerifyOTP)

	users := api.Group("/users", middleware.RequireAuth())
	users.Get("/me", userHandler.GetProfile)

	port := config.GetEnv("PORT", "3000")
	log.Printf("Server berjalan di port %s", port)
	log.Fatal(app.Listen(":" + port))
}