package role

// Role is a custom type representing a user role to ensure type safety.
type Role string

// Defines the valid role constants available in the application.
const (
	Admin Role = "admin"
	User  Role = "user"
)

// IsValid checks if the role is one of the predefined valid roles.
// It returns true if the role is valid, and false otherwise.
func (r Role) IsValid() bool {
	switch r {
	case Admin, User:
		return true
	}
	return false
}
