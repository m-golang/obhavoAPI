package users

import (
	"errors"
	"fmt"
	"havoAPI/internal/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UsersServiceInterface defines the methods that a user service should implement.
type UsersServiceInterface interface {
	InsertNewUser(name, surname, username, password string) error
	UserAuthentication(username, password string) (int, error)
	APIKeyAuthorization(apiKey string) (bool, error)
	FetchUserAPIKey(userID int) (string, error)
}

// UsersService is a concrete implementation of the UsersServiceInterface.
type UsersService struct {
	db model.DBContractUsers // The DBContractUsers interface represents the database operations for users.
}

// NewUsersService is a constructor function that returns a new instance of UsersService.
func NewUsersService(db model.DBContractUsers) *UsersService {
	return &UsersService{db: db}
}

// InsertNewUser inserts a new user into the database after hashing the password.
// Returns an error if there's an issue with the database or password hashing.
func (s *UsersService) InsertNewUser(name, surname, username, password string) error {
	// Hash the user's password using bcrypt
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// Return an error if password hashing fails
		return fmt.Errorf("error occurred while hashing password in the service section: %w", err)
	}

	// Insert the new user into the database
	userID, err := s.db.InsertUser(name, surname, username, hashed_password)
	if err != nil {
		// Check if the error is due to a duplicated username
		if errors.Is(err, model.ErrDuplicatedUsername) {
			return ErrUsernameExists
		}
		// Return the error if the insertion fails
		return fmt.Errorf("error occured while intersing user: %w", err)
	}

	err = s.GenerateNewApiKey(userID)
	if err != nil {
		return err
	}

	// Return nil if the user is successfully inserted and API key genereted successfully
	return nil
}

// UserAuthentication authenticates a user by verifying their username and password.
// It returns the user ID if the credentials are valid, or an error if invalid.
func (s *UsersService) UserAuthentication(username, password string) (int, error) {
	// Retrieve the stored credentials for the provided username
	userID, passwordHash, err := s.db.RetrieveUserCredentials(username)
	if err != nil {
		// Check if the error indicates the user does not exist
		if errors.Is(err, model.ErrUserNotFound) {
			return 0, ErrUserNotFound
		}
		// Return an error if retrieving credentials fails
		return 0, fmt.Errorf("error occured while retrieving user credentials: %w", err)
	}

	// Compare the provided password with the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		// Return an error if the passwords do not match
		return 0, ErrInvalidUserCredentials
	}

	// Return the user ID if authentication is successful
	return userID, nil
}

// GenerateNewApiKey generates a new API key for the user.
// It creates a unique API key using UUID and inserts it into the database.
// If an error occurs during the insertion, it is returned.
func (s *UsersService) GenerateNewApiKey(userID int) error {
	// Generate a new unique API key using UUID
	newAPIKey := uuid.New().String()

	// Insert the generated API key into the database for the user
	err := s.db.InsertUserAPIKey(userID, newAPIKey)
	if err != nil {
		// Return an error if inserting the API key into the database fails
		return fmt.Errorf("error occured while inserting new api key: %w", err)
	}

	// Return nil indicating successful generation and insertion of the API key
	return nil
}

func (s *UsersService) FetchUserAPIKey(userID int) (string, error) {
	apiKey, err := s.db.RetriveUserAPIKey(userID)
	if err != nil {
		return "", fmt.Errorf("error occured while fetching user API key: %w", err)
	}

	return apiKey, nil
}

func (s *UsersService) APIKeyAuthorization(apiKey string) (bool, error) {
	isKeyTrue, err := s.db.CheckUserAPIKey(apiKey)
	if err != nil {
		if errors.Is(err, model.ErrAPIKeyNotFound) {
			return false, ErrAPIKeyNotFound
		}
		return false, fmt.Errorf("error occured while checking user API key: %w", err)
	}

	return isKeyTrue, nil
}
