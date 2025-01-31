package services

import (
	"errors"
	"fmt"
	"havoAPI/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UsersServiceInterface defines the methods that a user service should implement.
// This interface is used for managing user-related operations, including user creation,
// authentication, and API key management.
type UsersServiceInterface interface {
	// InsertNewUser inserts a new user into the system with the provided details.
	// It returns an error if there is an issue with the database or password hashing.
	InsertNewUser(name, surname, username, password string) error

	// UserAuthentication authenticates a user by verifying their username and password.
	// It returns the user ID if authentication is successful, or an error if the credentials are invalid.
	UserAuthentication(username, password string) (int, error)

	// FetchUserAPIKey retrieves the API key for a given user by user ID.
	// It returns the API key or an error if the retrieval fails.
	FetchUserAPIKey(userID int) (string, error)
}

// UsersService is a concrete implementation of the UsersServiceInterface.
// It provides user management functionalities, including user insertion, authentication, and API key generation.
type UsersService struct {
	// db is an instance of the DBContractUsers interface which handles user-related database operations.
	db models.DBContractUsers
}

// NewUsersService initializes and returns a new instance of the UsersService struct.
// This function is used to create a new UsersService instance with the provided database interface.
func NewUsersService(db models.DBContractUsers) *UsersService {
	return &UsersService{db: db}
}

// InsertNewUser inserts a new user into the database after hashing the password.
// It returns an error if there's an issue with the password hashing or database insertion.
// This function also generates a new API key for the user after successful insertion.
func (s *UsersService) InsertNewUser(name, surname, username, password string) error {
	// Hash the user's password using bcrypt to ensure secure storage.
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// Return an error if password hashing fails
		return fmt.Errorf("error occurred while hashing password in the service section: %w", err)
	}

	// Insert the new user into the database, and get the generated user ID.
	userID, err := s.db.InsertUser(name, surname, username, hashed_password)
	if err != nil {
		// Check if the error is due to a duplicated username.
		if errors.Is(err, models.ErrDuplicatedUsername) {
			return ErrUsernameExists
		}
		// Return any other error that occurred during user insertion.
		return fmt.Errorf("error occurred while inserting user: %w", err)
	}

	// Generate a new API key for the user after successfully inserting the user.
	err = s.GenerateNewApiKey(userID)
	if err != nil {
		return err
	}

	// Return nil if the user is successfully inserted and API key generated.
	return nil
}

// UserAuthentication authenticates a user by checking the provided username and password.
// It returns the user ID if the credentials are valid, or an error if the credentials are invalid.
func (s *UsersService) UserAuthentication(username, password string) (int, error) {
	// Retrieve the stored credentials for the provided username.
	userID, passwordHash, err := s.db.RetrieveUserCredentials(username)
	if err != nil {
		// Check if the error indicates the user does not exist.
		if errors.Is(err, models.ErrUserNotFound) {
			return 0, ErrUserNotFound
		}
		// Return any other error that occurred while retrieving user credentials.
		return 0, fmt.Errorf("error occurred while retrieving user credentials: %w", err)
	}

	// Compare the provided password with the stored password hash.
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		// Return an error if the passwords do not match.
		return 0, ErrInvalidUserCredentials
	}

	// Return the user ID if authentication is successful.
	return userID, nil
}

// GenerateNewApiKey generates a new API key for the user using UUID and inserts it into the database.
// It returns an error if the API key insertion fails.
func (s *UsersService) GenerateNewApiKey(userID int) error {
	// Generate a new unique API key using UUID for the user.
	newAPIKey := uuid.New().String()

	// Insert the generated API key into the database for the user.
	err := s.db.InsertUserAPIKey(userID, newAPIKey)
	if err != nil {
		// Return an error if inserting the API key into the database fails.
		return fmt.Errorf("error occurred while inserting new API key: %w", err)
	}

	// Return nil if the API key is successfully generated and inserted.
	return nil
}

// FetchUserAPIKey retrieves the API key for a specific user by their user ID.
// It returns the API key if found or an error if the retrieval fails.
func (s *UsersService) FetchUserAPIKey(userID int) (string, error) {
	// Retrieve the user's API key from the database using the user ID.
	apiKey, err := s.db.RetriveUserAPIKey(userID)
	if err != nil {
		// Return an error if fetching the API key fails.
		return "", fmt.Errorf("error occurred while fetching user API key: %w", err)
	}

	// Return the retrieved API key.
	return apiKey, nil
}
