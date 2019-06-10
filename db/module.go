package db

import (
	"fmt"
	"net/url"
	"time"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
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
		Host                  string     `inject:"config:db.host,optional"`
		Port                  string     `inject:"config:db.port,optional"`
		DatabaseName          string     `inject:"config:db.databaseName,optional"`
		Username              string     `inject:"config:db.user,optional"`
		Password              string     `inject:"config:db.password,optional"`
		MaxConnectionLifetime float64    `inject:"config:db.maxConnectionLifetime,optional"`
		ConnectionOptions     config.Map `inject:"config:db.connectionOptions,optional"`
	}
)

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind((*DB)(nil)).ToProvider(dbProvider).AsEagerSingleton()
	flamingo.BindEventSubscriber(injector).To(&ShutdownSubscriber{})
}

// DefaultConfig for the module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"db.connectionOptions": config.Map{
			"parseTime": "true", // required for correct handling of datetime in Scan
		},
	}
}

func dbProvider(cfg *dbConfig, logger flamingo.Logger) DB {
	dataSourceURL := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DatabaseName,
	)

	if len(cfg.ConnectionOptions) > 0 {
		options := url.Values{}
		for key, value := range cfg.ConnectionOptions {
			options.Set(key, value.(string))
		}
		dataSourceURL += "?" + options.Encode()
	}

	dbConnection, err := sqlx.Connect("mysql", dataSourceURL)
	if err != nil {
		logger.Fatal(err)
	}

	if cfg.MaxConnectionLifetime != 0 {
		dbConnection.SetConnMaxLifetime(time.Second * time.Duration(cfg.MaxConnectionLifetime))
	}

	return &db{connection: dbConnection}
}

// Connection to the database
func (db *db) Connection() *sqlx.DB {
	return db.connection
}
