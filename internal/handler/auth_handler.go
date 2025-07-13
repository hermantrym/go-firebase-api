package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hermantrym/go-firebase-api/internal/apierror"
	"github.com/hermantrym/go-firebase-api/internal/service"
	"net/http"
)

// AuthHandler handles HTTP requests related to authentication.
type AuthHandler struct {
	userService service.UserService
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAuthHandler(svc service.UserService) *AuthHandler {
	return &AuthHandler{userService: svc}
}

// LoginRequest defines the expected JSON request body for the login endpoint.
type LoginRequest struct {
	// Email is the user's email address, required for login.
	Email string `json: "email" binding:"required,email"`
}

// Login handles the user login request. It validates the request body,
// calls the user service to generate a JWT, and returns the token upon success.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	// Bind and validate the incoming JSON payload.
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := apierror.NewBadRequestError("Invalid request body: email is required and must be valid")
		c.JSON(apiErr.Code, apiErr)
		return
	}

	// Call the service to perform the login logic and generate a token.
	token, err := h.userService.LoginUser(c.Request.Context(), req.Email)
	if err != nil {
		var apiErr *apierror.APIError
		// Check if the error is a custom APIError (e.g., NotFoundError) for a specific response.
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.Code, apiErr)
		} else {
			// Fallback for unexpected errors.
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
		}
		return
	}

	// Return the token in the response.
	c.JSON(http.StatusOK, gin.H{"token": token})
}
