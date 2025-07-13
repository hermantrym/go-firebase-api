package repository

import (
	"context"
	"errors"
	"github.com/hermantrym/go-firebase-api/internal/apierror"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/hermantrym/go-firebase-api/internal/model"
)

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	CreateUser(ctx context.Context, user model.User) (*model.User, error)
	GetUser(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]model.User, error)
}

// userRepository is the concrete implementation of UserRepository that interacts with Firestore.
type userRepository struct {
	client *firestore.Client
}

// NewUserRepository creates a new instance of the user repository.
func NewUserRepository(client *firestore.Client) UserRepository {
	return &userRepository{client: client}
}

// CreateUser adds a new user document to the "users" collection in Firestore.
func (r *userRepository) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	// Create a new document with a random ID in the "users" collection.
	docRef, _, err := r.client.Collection("users").Add(ctx, map[string]interface{}{
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})

	if err != nil {
		log.Printf("Error creating user in database: %v", err)
		return nil, apierror.NewInternalServerError("Failed to create user in database")
	}

	// Set the auto-generated ID on the user model and return it.
	user.ID = docRef.ID
	return &user, nil
}

// GetUser retrieves a single user document by its ID from Firestore.
func (r *userRepository) GetUser(ctx context.Context, id string) (*model.User, error) {
	docSnap, err := r.client.Collection("users").Doc(id).Get(ctx)

	if err != nil {
		// Specifically handle the case where the document is not found.
		if status.Code(err) == codes.NotFound {
			return nil, apierror.NewNotFoundError("User with ID '" + id + "' not found")
		}

		log.Printf("Error getting user from database: %v", err)
		return nil, apierror.NewInternalServerError("Failed to retrieve user from database")
	}

	var user model.User
	// Map the Firestore document data to the User struct.
	if err := docSnap.DataTo(&user); err != nil {
		log.Printf("Error converting user data: %v", err)
		return nil, apierror.NewInternalServerError("Failed to process user data")
	}

	user.ID = docSnap.Ref.ID
	return &user, nil
}

// GetAllUsers retrieves all user documents from the "users" collection.
func (r *userRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	iter := r.client.Collection("users").Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		// iterator.Done signifies that all documents have been processed.
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			log.Printf("Error iterating users: %v", err)
			return nil, apierror.NewInternalServerError("Failed to retrieve users")
		}

		var user model.User
		if err := doc.DataTo(&user); err != nil {
			log.Printf("Error converting user data for doc ID %s: %v", doc.Ref.ID, err)
			// Continue to the next document if one fails to convert, or return an error.
			// For this implementation, we will return an error to ensure data integrity.
			return nil, apierror.NewInternalServerError("Failed to process user data")
		}

		user.ID = doc.Ref.ID
		users = append(users, user)
	}

	return users, nil
}

// GetUserByEmail retrieves a single user document by their email address.
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	// Query the "users" collection for a document with a matching email field.
	iter := r.client.Collection("users").Where("email", "==", email).Limit(1).Documents(ctx)
	// Ensure the iterator is always closed to release resources.
	defer iter.Stop()

	doc, err := iter.Next()
	if err != nil {
		// The iterator returns a specific error when there are no more documents.
		if err.Error() == "iterator: no more items" {
			return nil, apierror.NewNotFoundError("User with email '" + email + "' not found")
		}

		log.Printf("Error getting user by email from database: %v", err)
		return nil, apierror.NewInternalServerError("Failed to retrieve user from database")
	}

	var user model.User
	// Map the Firestore document data to the User struct.
	if err := doc.DataTo(&user); err != nil {
		log.Printf("Error converting user data: %v", err)
		return nil, apierror.NewInternalServerError("Failed to process user data")
	}

	user.ID = doc.Ref.ID
	return &user, nil
}
