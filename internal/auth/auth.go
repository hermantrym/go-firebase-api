package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hermantrym/go-firebase-api/internal/apierror"
)

// JWTClaims defines the custom claims to be stored in the JWT payload.
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new signed JWT for a given user.
// It relies on the JWT_SECRET_KEY environment variable.
func GenerateJWT(userID, email string) (string, error) {
	// Retrieve the secret key from environment variables.
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return "", errors.New("JWT_SECRET_KEY environment variable not set")
	}

	// Set the token's expiration time (e.g., 24 hours).
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the JWT claims, including custom and registered claims.
	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-firebase-api",
		},
	}

	// Create a new token with the claims and HS256 signing method.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key to get the complete token string.
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// AuthMiddleware creates a gin middleware to verify the JWT from the Authorization header.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := os.Getenv("JWT_SECRET_KEY")
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			err := apierror.NewAPIError(http.StatusUnauthorized, "Authorization header is required")
			c.AbortWithStatusJSON(err.Code, err)
			return
		}

		// The token is expected in the format "Bearer <token>".
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			err := apierror.NewAPIError(http.StatusUnauthorized, "Authorization header format must be Bearer {token}")
			c.AbortWithStatusJSON(err.Code, err)
			return
		}

		tokenString := parts[1]
		claims := &JWTClaims{}

		// Parse and validate the token.
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Provide the key for signature verification.
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			apiErr := apierror.NewAPIError(http.StatusUnauthorized, "Invalid or expired token")
			c.AbortWithStatusJSON(apiErr.Code, apiErr)
			return
		}

		// Store the user ID in the context for use by subsequent handlers.
		c.Set("userID", claims.UserID)

		// Continue to the next handler.
		c.Next()
	}
}
