package models

import "errors"

// ErrUserNotFound is returned when a user cannot be found in the database.
// This error is useful when attempting to retrieve a user by their ID or username,
// but the user does not exist in the database.
var ErrUserNotFound = errors.New("models: User not found")

// ErrDuplicatedUsername is returned when a username already exists in the database.
// This error occurs when a new user attempts to register with a username
// that is already taken by another user in the system.
var ErrDuplicatedUsername = errors.New("models: Username already exists")

// ErrAPIKeyNotFound is returned when an API key cannot be found.
// This error occurs when an API request is made with an invalid or missing API key,
// and the application cannot locate a valid API key for the user.
var ErrAPIKeyNotFound = errors.New("models: API Key not found")
