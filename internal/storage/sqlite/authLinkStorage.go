package sqlite

import "fmt"

type AuthLinksStorage struct {
	db Database
}

func (s *AuthLinksStorage) MustInit() {
	query := `CREATE TABLE IF NOT EXISTS auth_links(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
		code TEXT NOT NULL UNIQUE,
		status TEXT NOT NULL
	)`

	_, err := s.db.Exec(query)

	if err != nil {
		panic(fmt.Errorf("failed to init auth_links table: %w", err))
	}
}
