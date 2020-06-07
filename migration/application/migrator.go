package application

import (
	"os"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"

	"github.com/tessig/flamingo-mysql/db"
)

type (
	// Migrator can migrate the database up and down using the migration scripts
	Migrator struct {
		db                 db.DB
		logger             flamingo.Logger
		databaseName       string
		migrationDirectory string
		migration          *migrate.Migrate
	}
)

func migrationFactory(m *Migrator) *migrate.Migrate {
	conn := m.db.Connection()
	driver, err := mysql.WithInstance(conn.DB, &mysql.Config{})
	if err != nil {
		panic(err)
	}
	migration, err := migrate.NewWithDatabaseInstance(
		"file://"+m.migrationDirectory,
		m.databaseName,
		driver,
	)
	if err != nil {
		panic(err)
	}

	return migration
}

// Inject dependencies
func (m *Migrator) Inject(
	db db.DB,
	logger flamingo.Logger,
	conf *struct {
		DatabaseName       string `inject:"config:mysql.db.databaseName,optional"`
		MigrationDirectory string `inject:"config:mysql.migration.directory"`
	},
) {
	m.db = db
	m.logger = logger
	m.databaseName = conf.DatabaseName
	if dbname, ok := os.LookupEnv("DBNAME"); ok {
		m.databaseName = dbname
	}
	m.migrationDirectory = conf.MigrationDirectory
}

// Up looks at the currently active migration version
// and will migrate step versions up or applying all up migrations if step is not given
func (m *Migrator) Up(steps *int) error {
	m.migration = migrationFactory(m)
	return m.runMigration(m.migration.Up, steps)
}

// Down looks at the currently active migration version
// and will migrate step versions down or applying all down migrations if step is not given
func (m *Migrator) Down(steps *int) error {
	m.migration = migrationFactory(m)
	if steps != nil {
		tmpSteps := -*steps
		steps = &tmpSteps
	}
	return m.runMigration(m.migration.Down, steps)
}

func (m *Migrator) runMigration(migratorFunc func() error, steps *int) error {
	var err error

	logger := m.logger.WithField(flamingo.LogKeyCategory, "migrations")

	logger.Info("Run migrations...")

	if steps == nil {
		err = migratorFunc()
	} else {
		err = m.migration.Steps(*steps)
	}

	if err == migrate.ErrNoChange {
		logger.Info("migrations: No change")
		return nil
	} else if err != nil {
		panic(err)
	}

	logger.Info("Migrations complete")

	return err
}
