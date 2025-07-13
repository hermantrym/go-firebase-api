package service

import (
	"context"
	"github.com/hermantrym/go-firebase-api/internal/apierror"
	"github.com/hermantrym/go-firebase-api/internal/auth"
	"log"

	"github.com/hermantrym/go-firebase-api/internal/model"
	"github.com/hermantrym/go-firebase-api/internal/repository"
)

// UserService defines the interface for user-related business logic.
type UserService interface {
	RegisterUser(ctx context.Context, user model.User) (*model.User, error)
	FindUserByID(ctx context.Context, id string) (*model.User, error)
	LoginUser(ctx context.Context, email string) (string, error)
}

// userService is the concrete implementation of the UserService interface.
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new instance of userService.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{userRepo: repo}
}

// RegisterUser handles the business logic for creating a new user.
func (s *userService) RegisterUser(ctx context.Context, user model.User) (*model.User, error) {
	return s.userRepo.CreateUser(ctx, user)
}

// LoginUser handles the user login process.
// It finds a user by email and generates a JWT if the user is found.
func (s *userService) LoginUser(ctx context.Context, email string) (string, error) {
	// Find the user by email.
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		// Return the error from the repository layer.
		return "", err
	}

	// If the user is found, generate a JWT.
	token, err := auth.GenerateJWT(user.ID, user.Email)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		return "", apierror.NewInternalServerError("Failed to generate authentication token")
	}

	return token, nil
}

// FindUserByID retrieves a user by their unique ID.
func (s *userService) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.GetUser(ctx, id)
}
