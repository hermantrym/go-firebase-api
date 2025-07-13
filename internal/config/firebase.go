package config

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// InitializeFirebase sets up the connection to Google Firestore and returns a client instance.
// It relies on the FIREBASE_SERVICE_ACCOUNT_KEY_PATH environment variable for the credentials file path.
// The application will exit with a fatal error if initialization fails at any step.
func InitializeFirebase() *firestore.Client {
	ctx := context.Background()

	// Get the service account key file path from an environment variable.
	serviceAccountKeyPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY_PATH")

	// Ensure the environment variable is set.
	if serviceAccountKeyPath == "" {
		log.Fatalf("FIREBASE_SERVICE_ACCOUNT_KEY_PATH environment variable not set.")
	}

	// Create a client option with the credentials file.
	opt := option.WithCredentialsFile(serviceAccountKeyPath)
	// Initialize the Firebase app.
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Get a Firestore client from the initialized app.
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing Firestore client: %v\n", err)
	}

	return client
}
