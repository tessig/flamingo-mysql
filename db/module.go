package db

import (
	"fmt"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/flamingo"

	// import the mysql driver here to have it when module is activated
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type (
	// Module basic struct
	Module struct{}

	// DB interface for database connections
	DB interface {
		Connection() *sqlx.DB
	}

	db struct {
		connection *sqlx.DB
	}

	dbConfig struct {
		Host         string `inject:"config:db.host,optional"`
		Port         string `inject:"config:db.port,optional"`
		DatabaseName string `inject:"config:db.databaseName,optional"`
		Username     string `inject:"config:db.user,optional"`
		Password     string `inject:"config:db.password,optional"`
	}
)

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind((*DB)(nil)).ToProvider(dbProvider).AsEagerSingleton()
	flamingo.BindEventSubscriber(injector).To(&ShutdownSubscriber{})
}

func dbProvider(cfg *dbConfig, logger flamingo.Logger) DB {
	dbConnection, err := sqlx.Connect("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DatabaseName,
	))
	if err != nil {
		logger.Fatal(err)
	}
	return &db{connection: dbConnection}
}

// Connection to the database
func (db *db) Connection() *sqlx.DB {
	return db.connection
}
