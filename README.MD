# Modular REST API with Go & Firebase

[![Go Version](https://img.shields.io/badge/go-1.18+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A robust and modular RESTful API built with Golang, using Firebase (Firestore) for data persistence. This project demonstrates professional backend development practices including a clean, layered architecture, JWT authentication, role-based authorization, input validation, and structured error handling.

---

## Features

-   **Modular Architecture**: Clean separation of concerns using a layered structure (Handler, Service, Repository).
-   **JWT Authentication**: Secure endpoints using a JWT-based authentication middleware.
-   **Role-Based Authorization (RBAC)**: Securely restricts access based on user roles. Features separate endpoints for public registration and admin-level user management.
-   **Configuration Management**: Securely manages configuration and secrets using environment variables (`.env` file).
-   **Input Validation**: Strong server-side validation of request data using `go-playground/validator`.
-   **Structured Error Handling**: A custom error handling system to provide clear, consistent error responses for different scenarios.
-   **Firebase Integration**: Uses the Firebase Admin SDK for Go to interact with Cloud Firestore.

---

## Project Structure

The project follows a standard layered architecture to ensure maintainability and scalability.

```
.
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/
│   ├── apierror/
│   │   └── apierror.go       # Custom error types
│   ├── auth/
│   │   └── auth.go           # JWT generation and middleware
│   ├── config/
│   │   └── firebase.go       # Firebase initialization
│   ├── handler/
│   │   ├── auth_handler.go   # HTTP handler for authentication
│   │   └── user_handler.go   # HTTP handler for user resources
│   ├── model/
│   │   └── user.go           # User data structure
│   ├── repository/
│   │   └── user_repository.go# Data access layer (Firestore)
│   └── service/
│       └── user_service.go   # Business logic layer
├── .env                        # Local environment variables (gitignored)
├── .gitignore
├── go.mod
├── go.sum
└── serviceAccountKey.json      # Firebase credentials (gitignored)
```

---

## Prerequisites

-   [Go](https://golang.org/dl/) version 1.18 or higher.
-   A [Google Firebase](https://console.firebase.google.com/) project with Firestore enabled.

---

## Installation & Setup

Follow these steps to get the project running on your local machine.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/hermantrym/go-firebase-api.git
    cd go-firebase-api
    ```

2.  **Set up Firebase Credentials:**
    -   Go to your Firebase project settings > Service Accounts.
    -   Click **"Generate new private key"** to download a JSON file.
    -   Rename the downloaded file to `serviceAccountKey.json` and place it in the root directory of the project.

3.  **Configure Environment Variables:**
    -   Create a new file named `.env` in the root directory. You can copy the `.env.example` file if it exists.
    -   Open the `.env` file and set the required variables. See the [Environment Variables](#environment-variables) section below for details.

4.  **Install Dependencies:**
    ```bash
    go mod tidy
    ```

5.  **Run the Application:**
    ```bash
    go run ./cmd/api/main.go
    ```
    The server will start on `http://localhost:8080`.

---

## API Endpoints

### Authentication

#### 1. Login to Get a Token

-   **Method**: `POST`
-   **Path**: `/login`
-   **Description**: Authenticates a user based on their email and returns a JWT if successful.
-   **Access**: Public

**Request Body:**
```json
{
    "email": "user@example.com"
}
```

**Success Response (200 OK):**
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### User Management

#### 1. Register a New User

-   **Method**: `POST`
-   **Path**: `/users`
-   **Description**: Creates a new user with a default "user" role. Any `role` field provided in the request body will be ignored for security reasons.
-   **Access**: Public

**Request Body:**
```json
{
  "name": "Budi Santoso",
  "email": "budi.santoso@example.com"
}
```

**Success Response (201 Created):**
```json
{
  "id": "some-generated-id",
  "name": "Budi Santoso",
  "email": "budi.santoso@example.com",
  "role": "user"
}
```

#### 2. Get User Details by ID

-   **Method**: `GET`
-   **Path**: `/users/:id`
-   **Description**: Retrieves the details of a specific user.
-   **Access**: **Protected** (Requires a valid JWT for any authenticated user)

**Example Request:**
```bash
TOKEN="<your-jwt>"
USER_ID="<user-id-to-fetch>"

curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/users/$USER_ID
```

**Success Response (200 OK):**
```json
{
    "id": "some-user-id",
    "name": "Budi Santoso",
    "email": "budi.santoso@example.com",
    "role": "user"
}
```

### Admin Endpoints

#### 1. Get All Users

-   **Method**: `GET`
-   **Path**: `/admin/users`
-   **Description**: Retrieves a list of all users in the system.
-   **Access**: **Protected (Admin Only)**

**Example Request:**
```bash
# Ensure this token belongs to a user with the 'admin' role
ADMIN_TOKEN="<your-admin-jwt>"

curl -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:8080/admin/users
```

**Success Response (200 OK):**
```json
[
    {
        "id": "user-id-1",
        "name": "Admin User",
        "email": "admin@example.com",
        "role": "admin"
    },
    {
        "id": "user-id-2",
        "name": "Budi Santoso",
        "email": "budi.santoso@example.com",
        "role": "user"
    }
]
```

#### 2. Create a New User (Admin)

-   **Method**: `POST`
-   **Path**: `/admin/users`
-   **Description**: Allows an admin to create a new user with a specific role. If the `role` is omitted, it defaults to "user".
-   **Access**: **Protected (Admin Only)**

**Request Body:**
```json
{
    "name": "Admin Baru",
    "email": "admin.baru@example.com",
    "role": "admin"
}
```

**Example Request:**
```bash
ADMIN_TOKEN="<your-admin-jwt>"

curl -X POST -H "Authorization: Bearer $ADMIN_TOKEN" \
-H "Content-Type: application/json" \
-d '{"name": "Admin Baru", "email": "admin.baru@example.com", "role": "admin"}' \
http://localhost:8080/admin/users
```

**Success Response (200 OK):**
```json
{
    "id": "another-generated-id",
    "name": "Admin Baru",
    "email": "admin.baru@example.com",
    "role": "admin"
}
```

---

## Environment Variables

These variables must be defined in a `.env` file in the project root.

| Variable                            | Description                                                      | Example                               |
|-------------------------------------|------------------------------------------------------------------|---------------------------------------|
| `FIREBASE_SERVICE_ACCOUNT_KEY_PATH` | The file path to your Firebase service account JSON credentials. | `./serviceAccountKey.json`            |
| `JWT_SECRET_KEY`                    | A long, random, and secret string used to sign and verify JWTs.  | `a-very-strong-and-random-secret-key` |

---

## Key Technologies Used

-   **Go**: The core programming language.
-   **Gin**: A high-performance HTTP web framework.
-   **Firebase Admin SDK**: For connecting to Cloud Firestore.
-   **JWT-Go**: For generating and validating JSON Web Tokens.
-   **Go-Playground Validator**: For request data validation.

---

## Future Improvements

-   [x] **Role-Based Access Control (RBAC)**: Restricts access to specific endpoints based on user roles.
-   [ ] **Password Hashing**: Implement `bcrypt` for secure password storage and authentication.
-   [ ] **Unit & Integration Tests**: Write comprehensive tests for all layers of the application.
-   [ ] **Structured Logging**: Integrate a logging library like `Logrus` or `Zap` for better log management.
-   [ ] **Dockerize the Application**: Create a `Dockerfile` to containerize the application for easier deployment.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE.md) file for details.