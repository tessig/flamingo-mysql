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
		ConnectionOptions     config.Map `inject:"config:mysql.db.connectionOptions,optional"`
		Host                  string     `inject:"config:mysql.db.host,optional"`
		Port                  string     `inject:"config:mysql.db.port,optional"`
		DatabaseName          string     `inject:"config:mysql.db.databaseName,optional"`
		Username              string     `inject:"config:mysql.db.user,optional"`
		Password              string     `inject:"config:mysql.db.password,optional"`
		MaxConnectionLifetime float64    `inject:"config:mysql.db.maxConnectionLifetime,optional"`
	}
)

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind((*DB)(nil)).ToProvider(dbProvider).AsEagerSingleton().In(dingo.ChildSingleton)
	flamingo.BindEventSubscriber(injector).To(&ShutdownSubscriber{})
}

// FlamingoLegacyConfigAlias maps legacy config entries to new ones
func (m *Module) FlamingoLegacyConfigAlias() map[string]string {
	return map[string]string{
		"db": "mysql.db",
	}
}

// CueConfig for the module
func (m *Module) CueConfig() string {
	return `
mysql: {
	DefaultConnectionOptions:: {
		parseTime: "true"
	}
	db: {
		host: string | *"localhost"
		port: string | *"3306"
		databaseName: string | *""
		user: string | *""
		password: string | *""
		maxConnectionLifetime: float | *0
		connectionOptions: DefaultConnectionOptions & {
			[string]: string
		}
	}
}`
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
			if v, ok := value.(string); ok {
				options.Set(key, v)
			}
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
