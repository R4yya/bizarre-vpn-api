package sqlite

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"bizarre-vpn-api/internal/models"
	"bizarre-vpn-api/internal/shared/fs"
	"bizarre-vpn-api/internal/shared/logger/sl"
	"bizarre-vpn-api/internal/storage"
)

type Database = storage.Database
type Executor = storage.Executor

type Storage struct {
	db                              Database
	SubscriptionPlanStorage         *subscriptionPlanStorage
	UserStorage                     *UserStorage
	LnkUserProviderStorage          *LnkUserProviderStorage
	BackendTypesStorage             *LibraryItemStorage
	ProtocolsStorage                *LibraryItemStorage
	LnkProtocolsBackendTypesStorage *LnkProtocolsBackendTypesStorage
	VpnServersStorage               *VpnServersStorage
	AuthLinksStorage                *AuthLinksStorage
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

	db.MustExec("PRAGMA foreign_keys = ON")

	log.Info(fmt.Sprintf("Connected to SQLite database at %s", dbPath))

	userStorage := &UserStorage{db}
	userStorage.MustInit()

	lnkUserProviderStorage := &LnkUserProviderStorage{
		db,
	}
	lnkUserProviderStorage.MustInit()

	subscriptionPlanStorage := &subscriptionPlanStorage{db}
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

	lnkProtocolsBackendTypesStorage := &LnkProtocolsBackendTypesStorage{db}

	lnkProtocolsBackendTypesStorage.MustInit()

	vpnServersStorage := &VpnServersStorage{db}

	vpnServersStorage.MustInit()

	authLinksStorage := &AuthLinksStorage{db}

	authLinksStorage.MustInit()

	storage := &Storage{
		db:                              db,
		SubscriptionPlanStorage:         subscriptionPlanStorage,
		UserStorage:                     userStorage,
		LnkUserProviderStorage:          lnkUserProviderStorage,
		BackendTypesStorage:             backendTypesStorage,
		ProtocolsStorage:                protocolsStorage,
		LnkProtocolsBackendTypesStorage: lnkProtocolsBackendTypesStorage,
		VpnServersStorage:               vpnServersStorage,
		AuthLinksStorage:                authLinksStorage,
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
