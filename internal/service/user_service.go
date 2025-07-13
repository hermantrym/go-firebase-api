package service

import (
	"context"
	"github.com/hermantrym/go-firebase-api/internal/apierror"
	"github.com/hermantrym/go-firebase-api/internal/auth"
	"github.com/hermantrym/go-firebase-api/internal/role"
	"log"

	"github.com/hermantrym/go-firebase-api/internal/model"
	"github.com/hermantrym/go-firebase-api/internal/repository"
)

// UserService defines the interface for user-related business logic.
type UserService interface {
	RegisterUser(ctx context.Context, user model.User) (*model.User, error)
	AdminRegisterUser(ctx context.Context, user model.User) (*model.User, error)
	FindUserByID(ctx context.Context, id string) (*model.User, error)
	LoginUser(ctx context.Context, email string) (string, error)
	FindAllUsers(ctx context.Context) ([]model.User, error)
}

// userService is the concrete implementation of the UserService interface.
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new instance of userService.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{userRepo: repo}
}

// RegisterUser handles the business logic for creating a new user with a default "user" role.
func (s *userService) RegisterUser(ctx context.Context, user model.User) (*model.User, error) {
	// Always assign the default "user" role for public registrations.
	user.Role = role.User
	return s.userRepo.CreateUser(ctx, user)
}

// AdminRegisterUser handles user creation by an administrator.
// It allows specifying a role, defaulting to "user" if none is provided,
// and validates the role before creation.
func (s *userService) AdminRegisterUser(ctx context.Context, user model.User) (*model.User, error) {
	// If no role is specified in the request, assign the default "user" role.
	if user.Role == "" {
		user.Role = role.User
	}

	// Validate that the provided role is a valid one (e.g., "admin" or "user").
	if !user.Role.IsValid() {
		return nil, apierror.NewBadRequestError("Invalid role specified")
	}

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
	token, err := auth.GenerateJWT(user.ID, user.Email, user.Role)
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

func (s *userService) FindAllUsers(ctx context.Context) ([]model.User, error) {
	return s.userRepo.GetAllUsers(ctx)
}
