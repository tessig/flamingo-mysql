package migration_test

import (
	"testing"

	"flamingo.me/flamingo/v3/framework/config"

	"github.com/tessig/flamingo-mysql/migration"
)

func TestModule_Configure(t *testing.T) {
	module := new(migration.Module)

	if err := config.TryModules(module.DefaultConfig(), module); err != nil {
		t.Error(err)
	}
}
