package db

import (
	"database/sql"
	"errors"
	"log"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/ndrscodes/yasp/internal/db/migrations"
)

type DbLogger struct {
	log     *log.Logger
	verbose bool
}

func (l *DbLogger) Verbose() bool {
	return l.verbose
}

func (l *DbLogger) Printf(format string, args ...any) {
	log.Printf(format, args...)
}

type Migrator struct {
	migrate *migrate.Migrate
	verbose bool
}

func NewMigrator(db *sql.DB, verbose bool) (*Migrator, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	if db.Ping() != nil {
		return nil, errors.New("unable to ping database, are we connected?")
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	source, err := iofs.New(migrations.Migrations, ".")
	if err != nil {
		return nil, err
	}

	mig, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return nil, err
	}

	mig.Log = &DbLogger{
		log:     log.Default(),
		verbose: verbose,
	}

	return &Migrator{
		migrate: mig,
		verbose: verbose,
	}, nil
}

func (m *Migrator) Down() error {
	err := m.migrate.Down()
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
		log.Println("No migrations to roll back.")
	}
	return nil
}

func (m *Migrator) Up() error {
	v, d, err := m.migrate.Version()
	if err != nil {
		if err.Error() == "no migration" {
			slog.Info("no previous migration detected")
		} else {
			return err
		}
	}
	m.migrate.Log.Printf("Applying database updates (current version: %d, dirty: %v)", v, d)

	err = m.migrate.Up()
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
		log.Println("Database is already up to date!")
	}

	return nil
}
