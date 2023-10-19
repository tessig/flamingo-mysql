package db_test

import (
	"testing"

	"flamingo.me/flamingo/v3/framework/config"

	"github.com/tessig/flamingo-mysql/db"
)

func TestModule_Configure(t *testing.T) {
	t.Parallel()

	if err := config.TryModules(nil, new(db.Module)); err != nil {
		t.Error(err)
	}
}
