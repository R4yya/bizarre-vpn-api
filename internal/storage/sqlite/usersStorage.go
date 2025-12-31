package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"bizarre-vpn-api/internal/storage"
	"bizarre-vpn-api/internal/storage/models"

	"github.com/jmoiron/sqlx"
)

type UserStorage struct {
	db Database
}

func (u *UserStorage) MustInit() {
	query := `CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		login TEXT UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		role VARCHAR(255),
		refresh_token TEXT(500),
		password TEXT
	);`

	//CREATE INDEX IF NOT EXISTS idx_id ON users(id)

	_, err := u.db.Exec(query)

	if err != nil {
		panic(fmt.Errorf("failed to init users table: %w", err))
	}
}

func (u *UserStorage) GetUsersList() (*[]models.BaseUser, error) {
	var usersList []models.FullUser

	query := `SELECT 
	id,
	username,
	login,
	role,
	created_at,
	updated_at 
	FROM users`
	err := u.db.Select(&usersList, query)

	if err != nil {
		return nil, fmt.Errorf("getting users list error: %w", err)
	}

	var basicUsersList []models.BaseUser

	for _, fullUser := range usersList {
		basicUsersList = append(basicUsersList, fullUser.BaseUser)
	}

	return &basicUsersList, nil
}

// GetUserByTelegramID gets the user by Telegram ID
func (u *UserStorage) GetUserById(ID int64, executor storage.Executor) (*models.BaseUser, error) {
	if executor == nil {
		executor = u.db
	}

	var user models.FullUser

	query := `SELECT 
	id,
	username,
	login,
	role,
	created_at,
	updated_at 
	FROM users WHERE id = ?`
	err := sqlx.Get(executor, &user, query, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user.BaseUser, nil
}

func (u *UserStorage) GetUserByLogin(login string) (*models.BaseUser, error) {
	var user models.FullUser

	query := `SELECT 
	id,
	username,
	login,
	role,
	created_at,
	updated_at 
	FROM users WHERE login = ?`
	err := u.db.Get(&user, query, login)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user.BaseUser, nil
}

func (u *UserStorage) GetUserPasswordHash(userId int64) (string, error) {
	const query = `SELECT password FROM users WHERE id = ?`

	var hashedPassword string

	err := u.db.Get(&hashedPassword, query, userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrUserNotFound
		}

		return "", fmt.Errorf("failed to get user password hash: %w", err)
	}

	return hashedPassword, nil
}

// CreateUser add a new user to the database
func (u *UserStorage) CreateUser(
	payload *models.CreateUserPayload,
	executor storage.Executor,
) (*models.BaseUser, error) {
	if executor == nil {
		executor = u.db
	}

	user := &models.FullUser{
		BaseUser: models.BaseUser{
			Username: payload.Username,
			Role:     payload.Role,
		},
		RefreshToken: "",
	}

	query := `
    INSERT INTO users (login, username, refresh_token, role, password)
    VALUES (:login, :username, :refresh_token, :role, :password)`

	rows, err := sqlx.NamedQuery(
		executor,
		query,
		map[string]interface{}{
			"login":         payload.Login,
			"username":      user.Username,
			"refresh_token": user.RefreshToken,
			"role":          user.Role,
			"password":      payload.Password,
		},
	)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.login") {
			return nil, storage.ErrLoginOccupied
		}

		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	defer rows.Close()

	var createdUser models.BaseUser

	if rows.Next() {
		err = rows.StructScan(&user)

		if err != nil {
			return nil, fmt.Errorf("failed to scan inserted user: %w", err)
		}
	}

	return &createdUser, nil
}

func (u *UserStorage) GetUserRefreshToken(ID int64) (string, error) {
	var refreshToken string

	const query = "SELECT refresh_token FROM users WHERE id = ?"
	err := u.db.Get(&refreshToken, query, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrUserNotFound
		}
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	return refreshToken, nil
}

func (u *UserStorage) UpdateUserRefreshToken(ID int64, refreshToken string) error {
	const query = `UPDATE users SET refresh_token = :refresh_token WHERE id = :id`

	params := map[string]interface{}{
		"id":            ID,
		"refresh_token": refreshToken,
	}

	_, err := u.db.NamedExec(query, params)
	if err != nil {
		return fmt.Errorf("failed to update user token for userId %d: %w", ID, err)
	}

	return nil
}
