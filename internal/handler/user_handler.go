package handler

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/hermantrym/go-firebase-api/internal/apierror"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hermantrym/go-firebase-api/internal/model"
	"github.com/hermantrym/go-firebase-api/internal/service"
)

// UserHandler handles HTTP requests related to users.
type UserHandler struct {
	userService service.UserService
	validate    *validator.Validate
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler(svc service.UserService, val *validator.Validate) *UserHandler {
	return &UserHandler{
		userService: svc,
		validate:    val,
	}
}

// CreateUser handles the POST /users endpoint.
// It parses the user data from the request body, validates it,
// and passes it to the user service for creation.
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user model.User

	// Bind the incoming JSON payload to the user struct.
	if err := c.ShouldBindJSON(&user); err != nil {
		apiErr := apierror.NewBadRequestError("Invalid JSON format")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Validate the user struct based on the defined tags.
	if err := h.validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	// Call the service to register the user.
	createdUser, err := h.userService.RegisterUser(c.Request.Context(), user)
	if err != nil {
		var apiErr *apierror.APIError
		// Check if the error is a custom APIError for specific HTTP responses.
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.Code, apiErr)
		} else {
			// Fallback for unexpected errors.
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
		}
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

// AdminCreateUser handles the POST /admin/users endpoint.
// This allows an administrator to create a new user, potentially with a specific role.
// It validates the incoming user data before creation.
func (h *UserHandler) AdminCreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		apiErr := apierror.NewBadRequestError("Invalid JSON format")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Validate the user struct based on the defined tags.
	if err := h.validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	// Call the service to register the user.
	createdUser, err := h.userService.AdminRegisterUser(c.Request.Context(), user)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.Code, apiErr)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
		}
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

// GetUser handles the GET /users/:id endpoint.
// It retrieves a user by the ID provided in the URL path.
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	user, err := h.userService.FindUserByID(c.Request.Context(), userID)

	if err != nil {
		var apiErr *apierror.APIError
		// Check if the error is a custom APIError (e.g., NotFoundError).
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.Code, apiErr)
		} else {
			// Fallback for unexpected errors.
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "An unexpected error occurred",
			})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetAllUsers handles the GET /admin/users endpoint.
// It retrieves a list of all users in the system.
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.FindAllUsers(c.Request.Context())
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.Code, apiErr)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
		}
		return
	}

	c.JSON(http.StatusOK, users)
}

// formatValidationErrors transforms validation errors from the validator library
// into a more readable map[string]string format for client consumption.
func formatValidationErrors(err error) map[string]string {
	errorsMap := make(map[string]string)

	// Type assert the error to access the slice of validation errors.
	for _, fieldErr := range err.(validator.ValidationErrors) {
		errorsMap[fieldErr.Field()] = "Field validation for '" + fieldErr.Field() + "' failed on the '" + fieldErr.Tag() + "' tag"
	}

	return errorsMap
}
