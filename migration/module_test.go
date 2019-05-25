package migration_test

import (
	"testing"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"

	"github.com/tessig/flamingo-mysql/migration"
)

func TestModule_Configure(t *testing.T) {
	module := new(migration.Module)
	cfgModule := &config.Module{
		Map: module.DefaultConfig(),
	}

	if err := dingo.TryModule(cfgModule, module); err != nil {
		t.Error(err)
	}
}
