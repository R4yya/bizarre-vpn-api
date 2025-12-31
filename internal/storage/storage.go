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
	ErrLoginOccupied     = errors.New("user with this login already exist")
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

// AuthLinks

var (
	ErrAuthLinkNotFound = errors.New("auth link not found")
)
