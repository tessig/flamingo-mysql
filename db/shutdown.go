package db

import (
	"context"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// ShutdownSubscriber handles the graceful app shutdown for the db connection
	ShutdownSubscriber struct {
		db DB
	}
)

// Inject dependencies
func (s *ShutdownSubscriber) Inject(db DB) {
	s.db = db
}

// Notify handles the incoming event if it is an AppShutdownEvent and closes the db connection
func (s *ShutdownSubscriber) Notify(_ context.Context, event flamingo.Event) {
	if _, ok := event.(*flamingo.ShutdownEvent); ok {
		// no error handling here because it's app shutdown. the connection will be hard aborted anyways
		_ = s.db.Connection().Close()
	}
}
