package storage

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type Executor = sqlx.Ext

// User errors
var (
	// ErrUserNotFound occurs if the user is not found in the database
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists occurs if a user with the specified Telegram ID already exists
	ErrUserAlreadyExists = errors.New("user with this Telegram ID already exists")
)

// Subscription plans errors
var (
	// ErrPlanNotFound occurs if the plan is not found in the database
	ErrPlanNotFound = errors.New("subscription plan not found")
)

// VpnServers Error

var (
	ErrVpnServerNotFound = errors.New("vpn server not found")
)
