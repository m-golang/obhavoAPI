package services

import (
	"errors"
	"fmt"
	"havoAPI/internal/model"

	"golang.org/x/crypto/bcrypt"
)

// UsersServiceInterface defines the methods that a user service should implement.
type UsersServiceInterface interface {
	// InsertNewUser handles the creation of a new user, including password hashing and insertion into the DB.
	InsertNewUser(name, surname, username, password string) error
	// UserAuthentication checks user credentials (username and password) and returns user ID if valid.
	UserAuthentication(username, password string) (int, error)
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
	err = s.db.InsertUser(name, surname, username, hashed_password)
	if err != nil {
		// Check if the error is due to a duplicated username
		if errors.Is(err, model.ErrDuplicatedUsername) {
			return ErrUsernameExists
		}
		// Return the error if the insertion fails
		return fmt.Errorf("error occured while intersing user: %w", err)
	}
	// Return nil if the user is successfully inserted
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
