package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"bizarre-vpn-api/internal/storage"
	"bizarre-vpn-api/internal/storage/models"

	"github.com/jmoiron/sqlx"
)

const (
	userBasicRole string = "basic"
)

type UserStorage struct {
	db Database
}

func (u *UserStorage) MustInit() {
	query := `CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		role VARCHAR(255),
		refresh_token TEXT(500)
	);`

	//CREATE INDEX IF NOT EXISTS idx_id ON users(id)

	_, err := u.db.Exec(query)

	if err != nil {
		panic(fmt.Errorf("failed to init users table: %w", err))
	}
}

// GetUserByTelegramID gets the user by Telegram ID
func (u *UserStorage) GetUserById(ID int64, executor storage.Executor) (*models.BaseUser, error) {
	if executor == nil {
		executor = u.db
	}

	var user models.FullUser

	fmt.Println("-----ID", ID)

	query := "SELECT * FROM users WHERE id = ?"
	err := sqlx.Get(executor, &user, query, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user.BaseUser, nil
}

// CreateUser adds a new user to the database
func (u *UserStorage) CreateUser(username string, executor storage.Executor) (userId int64, Err error) {
	if executor == nil {
		executor = u.db
	}

	user := &models.FullUser{
		BaseUser: models.BaseUser{
			Username: username,
			Role:     userBasicRole,
		},
		RefreshToken: "",
	}

	query := `
    INSERT INTO users (username, refresh_token, role)
    VALUES (:username, :refresh_token, :role)`

	result, err := sqlx.NamedExec(executor, query, user)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	return userID, nil
}

func (u *UserStorage) GetUserRefreshToken(ID int64) (string, error) {
	var refreshToken string

	query := "SELECT refresh_token FROM users WHERE id = ?"
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
