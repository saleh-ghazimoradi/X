package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var migrationsFS embed.FS

type Migration struct {
	db        *sql.DB
	dbname    string
	migration *migrate.Migrate
}

func (m *Migration) Up() error {
	if err := m.migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%w: failed to apply migration", err)
	}
	fmt.Println("Migration applied successfully")
	return nil
}

func (m *Migration) Rollback() error {
	if err := m.migration.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%w: failed to rollback migration", err)
	}
	fmt.Println("Migration rollback applied successfully")
	return nil
}

func (m *Migration) Close() error {
	if m.migration != nil {
		source, driver := m.migration.Close()
		if source != nil {
			return fmt.Errorf("%w: failed to close migration source", source)
		}
		if driver != nil {
			return fmt.Errorf("%w: failed to close migration driver", driver)
		}
	}
	return nil
}

func NewMigration(db *sql.DB, dbname string) (*Migration, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create migration driver", err)
	}

	source, err := iofs.New(migrationsFS, ".")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create migration source", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, dbname, driver)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create migration", err)
	}
	return &Migration{
		db:        db,
		dbname:    dbname,
		migration: m,
	}, nil
}
