package application

import (
	"fmt"
	"os"
	"path/filepath"

	"flamingo.me/flamingo/v3/framework/flamingo"

	"github.com/tessig/flamingo-mysql/db"
)

type (
	// Seeder can import the test data into the database
	Seeder struct {
		db             db.DB
		logger         flamingo.Logger
		seedsDirectory string
	}
)

// Inject dependencies
func (s *Seeder) Inject(
	db db.DB,
	logger flamingo.Logger,
	conf *struct {
		SeedsDirectory string `inject:"config:mysql.seed.directory"`
	},
) {
	s.db = db
	s.logger = logger
	s.seedsDirectory = conf.SeedsDirectory
}

// Seed runs all sql files in the seeding directory
func (s *Seeder) Seed() error {
	logger := s.logger.WithField(flamingo.LogKeyCategory, "seeds")

	err := filepath.Walk(s.seedsDirectory, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(info.Name()) == ".sql" {
			logger.Info("Seed file %s ...", info.Name())

			data, err := os.ReadFile(path)
			if err != nil {
				logger.Fatal("Problem while reading file content of %s:", info.Name())
				panic(err)
			}

			logger.Info("Seeding file contents...")

			s.db.Connection().MustExec(string(data))
			logger.Info("Seeding complete...")
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("seeding failed: %w", err)
	}

	return nil
}
