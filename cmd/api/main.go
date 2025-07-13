package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/hermantrym/go-firebase-api/internal/auth"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hermantrym/go-firebase-api/internal/config"
	"github.com/hermantrym/go-firebase-api/internal/handler"
	"github.com/hermantrym/go-firebase-api/internal/repository"
	"github.com/hermantrym/go-firebase-api/internal/service"
	"github.com/joho/godotenv"
)

// main is the entry point for the application.
// It initializes the configuration, database connection, dependency injection,
// router, and starts the HTTP server.
func main() {
	// Load Environment Variables
	// Load variables from the .env file.
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize Services & Dependencies
	// Initialize the Firestore client connection.
	firestoreClient := config.InitializeFirebase()
	// Ensure the client is closed gracefully when the application exits.
	defer func() {
		if err := firestoreClient.Close(); err != nil {
			log.Printf("ERROR: Failed to close Firestore client: %v", err)
		}
	}()
	// Create a new instance of the validator.
	validate := validator.New()

	// Dependency Injection
	// Wire together the application layers.
	userRepo := repository.NewUserRepository(firestoreClient)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, validate)
	authHandler := handler.NewAuthHandler(userService)

	// Setup Router (Gin)
	r := gin.Default()

	// --- PUBLIC ROUTES ---
	// Routes that can be accessed without authentication/token.
	r.POST("/login", authHandler.Login)
	r.POST("/users", userHandler.CreateUser) // Endpoint for user registration.

	// --- PROTECTED ROUTES ---
	// This group of routes requires a valid JWT.
	authorized := r.Group("/")
	authorized.Use(auth.AuthMiddleware())
	{
		// The endpoint to get user details is now protected.
		authorized.GET("/users/:id", userHandler.GetUser)
	}

	// --- PROTECTED ADMIN ROUTES ---
	// This group of routes is protected by two layers of middleware:
	// AuthMiddleware() - Ensures the user has a valid JWT.
	// RoleAuthMiddleware("admin") - Ensures the user has the 'admin' role.
	adminRoutes := r.Group("/admin")
	adminRoutes.Use(auth.AuthMiddleware())
	adminRoutes.Use(auth.RoleAuthMiddleware("admin"))
	{
		adminRoutes.GET("/users", userHandler.GetAllUsers)
		adminRoutes.POST("/users", userHandler.AdminCreateUser)
	}

	// Run Server
	log.Println("Server is running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
