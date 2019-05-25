package db_test

import (
	"testing"

	"flamingo.me/dingo"

	"github.com/tessig/flamingo-mysql/db"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(db.Module)); err != nil {
		t.Error(err)
	}
}
