package model

import "errors"

// ErrUserNotFound is returned when a user cannot be found in the database.
var ErrUserNotFound = errors.New("models: User not found")

// ErrDuplicatedUsername is returned when a username already exists in the database.
var ErrDuplicatedUsername = errors.New("models: Username already exists")

var ErrAPIKeyNotFound = errors.New("models: API Key not found")

