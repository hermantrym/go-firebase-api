package model

import "github.com/hermantrym/go-firebase-api/internal/role"

// User represents the data model for a user in the application.
// It includes struct tags for JSON serialization, Firestore mapping, and validation.
type User struct {
	// ID is the unique identifier for the user.
	// The `firestore:"-"` tag means this field is not stored within the Firestore document itself,
	// as the ID is used as the document's name.
	ID string `json:"id,omitempty" firestore:"-"`

	// Name is the user's full name.
	// It is a required field with a minimum length of 2 and a maximum of 100 characters.
	Name string `json:"name" firestore:"name" validate:"required,min=2,max=100"`

	// Email is the user's email address.
	// It is a required field and must be a valid email format.
	Email string `json:"email" firestore:"email" validate:"required,email"`

	// Role defines the user's authorization level (e.g., "admin", "user").
	Role role.Role `json:"role" firestore:"role"`
}
