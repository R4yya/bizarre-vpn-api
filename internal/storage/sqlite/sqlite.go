package sqlite

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"bizarre-vpn-api/internal/lib/fs"
	"bizarre-vpn-api/internal/lib/logger/sl"
	"bizarre-vpn-api/internal/storage/models"
)

type Database = *sqlx.DB

type Storage struct {
	db                      Database
	SubscriptionPlanStorage *SubscriptionPlanStorage
	UserStorage             *UserStorage
	LnkUserProviderStorage  *LnkUserProviderStorage
	BackendTypesStorage     *LibraryItemStorage
	ProtocolsStorage        *LibraryItemStorage
}

func MustInit(dbPath string, log *slog.Logger) *Storage {
	const op = "storage.sqlite.New"

	err := fs.CheckOrMakeDir(dbPath)
	if err != nil {
		log.Error("checkOrMakeDir err", sl.Err(fmt.Errorf("%v: %w", op, err)))
	}

	db, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		cErr := fmt.Errorf("%v: failed to connect to database: %w", op, err)

		log.Error("database initialization error", sl.Err(cErr))
		panic(cErr)
	}

	log.Info(fmt.Sprintf("Connected to SQLite database at %s", dbPath))

	userStorage := &UserStorage{db}
	userStorage.MustInit()

	lnkUserProviderStorage := &LnkUserProviderStorage{
		db,
		userStorage,
	}
	lnkUserProviderStorage.MustInit()

	subscriptionPlanStorage := &SubscriptionPlanStorage{db}
	subscriptionPlanStorage.MustInit()

	protocolsStorage := NewLibraryItemStorage(
		db,
		"protocols",
		&[]models.LibraryItem{
			{
				ID:   1,
				Name: "vless",
			},
		},
	)

	backendTypesStorage := NewLibraryItemStorage(
		db,
		"backend_types",
		&[]models.LibraryItem{
			{
				ID:   1,
				Name: "3xUI",
			},
		},
	)

	lnkProtocolsBackendTypesStorage := LnkProtocolsBackendTypesStorage{db}

	lnkProtocolsBackendTypesStorage.MustInit()

	_ = lnkProtocolsBackendTypesStorage

	storage := &Storage{
		db:                      db,
		SubscriptionPlanStorage: subscriptionPlanStorage,
		UserStorage:             userStorage,
		LnkUserProviderStorage:  lnkUserProviderStorage,
		BackendTypesStorage:     backendTypesStorage,
		ProtocolsStorage:        protocolsStorage,
	}

	return storage
}

// CloseDB closes the connection to the database
func (s *Storage) CloseDB() error {
	if s.db == nil {
		return nil
	}

	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	return nil
}
