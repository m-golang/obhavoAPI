package services

import "errors"

// ErrUserNotFound is returned when the requested user cannot be found in the system.
var ErrUserNotFound = errors.New("services: User not found")

// ErrUsernameExists is returned when an attempt is made to create a user with a username that already exists.
var ErrUsernameExists = errors.New("services: Username already exists")

// ErrInvalidUserCredentials is returned when the provided user credentials (username/password) are invalid.
var ErrInvalidUserCredentials = errors.New("services: Invalid user credentials")

var ErrAPIKeyNotFound = errors.New("models: API Key not found")
var ErrNoLocationFound = errors.New("no matching location found")

var ErrUnexpectedEndOfJSONInput = errors.New("unexpected end of JSON input")

var ErrNoDataCache = errors.New("no data in cache for location")
