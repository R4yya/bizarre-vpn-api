package custom_errors

import "errors"

// User errors
var (
	// ErrUserNotFound occurs if the user is not found in the database
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists occurs if a user with the specified Telegram ID already exists
	ErrUserAlreadyExists = errors.New("user with this Telegram ID already exists")
)
