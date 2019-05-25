package application

import (
	"context"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// StartUpMigrations subscribes to the AppStartup
	StartUpMigrations struct {
		migrator *Migrator
		logger   flamingo.Logger
	}
)

// Inject dependencies
func (s *StartUpMigrations) Inject(
	m *Migrator,
	l flamingo.Logger,
) {
	s.migrator = m
	s.logger = l
}

// Notify handles the automigration if configured on the AppStartupEvent
func (s *StartUpMigrations) Notify(_ context.Context, event flamingo.Event) {
	if _, ok := event.(*flamingo.StartupEvent); ok {
		s.logger.Info("Run auto migrations...")
		err := s.migrator.Up(nil)
		if err != nil {
			panic(err)
		}
	}
}
