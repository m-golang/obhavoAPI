package services

import "errors"

// ErrUserNotFound is returned when the requested user cannot be found in the system.
// This is typically used when a user attempts to log in with a non-existent account.
var ErrUserNotFound = errors.New("services: User not found")

// ErrUsernameExists is returned when an attempt is made to create a user with a username
// that already exists in the database. This helps in enforcing unique usernames.
var ErrUsernameExists = errors.New("services: Username already exists")

// ErrInvalidUserCredentials is returned when the provided user credentials (username/password)
// do not match any existing records in the system. It indicates failed authentication.
var ErrInvalidUserCredentials = errors.New("services: Invalid user credentials")

// ErrAPIKeyNotFound is returned when the provided API key does not exist in the database.
// This can occur when a user provides an invalid or expired API key during authentication.
var ErrAPIKeyNotFound = errors.New("models: API Key not found")

// ErrNoLocationFound is returned when no matching location is found for a weather query.
// This helps indicate that the location provided by the user does not exist or is not recognized.
var ErrNoLocationFound = errors.New("no matching location found")

// ErrUnexpectedEndOfJSONInput is returned when the JSON input provided is incomplete or malformed.
// It signals an error while parsing a request or response body due to missing data.
var ErrUnexpectedEndOfJSONInput = errors.New("unexpected end of JSON input")

// ErrNoDataCache is returned when a request for cached weather data cannot find any available data
// for the specified location. This may happen if the data has expired or hasn't been cached yet.
var ErrNoDataCache = errors.New("no data in cache for location")
