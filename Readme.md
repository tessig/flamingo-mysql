# Flamingo MySQL module [![Build Status](https://travis-ci.org/tessig/flamingo-mysql.svg?branch=master)](https://travis-ci.org/tessig/flamingo-mysql)

This flamingo module provides a simple MySQL implementation by wrapping 
github.com/jmoiron/sqlx and a migration tool by wrapping github.com/golang-migrate/migrate 
as [Flamingo](https://www.flamingo.me/) Modules.

## DB Module

The DB Flamingo module provides the interface `DB` and binds an sqlx connection 
as singleton to it. The module will panic on startup when the connection can't be established.

### Configuration

```yaml
db:
  host: "host"
  port: "3306" # must be a string!
  databaseName: "databaseName"
  user: "user"
  password: "password"
  maxConnectionLifetime: 0 # in seconds, 0 means to set nothing, negative values mean unlimited
  # a set of additional connection options which are added as parameters to the DB URL
  connectionOptions: 
    myOption1: "myValue1" # all option values must be strings
    myOption2: "false"
```

## Migration Module

The migration module relies on the db module and can handle schema migration and data seeding scripts. Both must be
provided as simple SQL scripts.

The module provides additional Flamingo commands as entrypoints:

* `migrate [up|down] (-s[number of steps])`
* `seed`

### Configuration

```yaml
migrations:
  automigrate: false,
  directory:   "sql/migrations/",
seeds:
  directory:   "sql/seeds/",
```

The Migration Module also adds `"db.connectionOptions.multiStatements": "true"` to the db configuration 
to handle migration and seed scripts.

### Migration

Migration scripts must be placed into the configured directory. For each migration,
there must be an "up" and a "down" script. Please refer to [github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate)
for more detailed documentation.

If you set the `automigrate` config to `true`, flamingo will run a `migrate up` on each application start (`flamingo.StartupEvent`).

### Seeding

Seeding scripts must be placed into the configured directory. The seed command runs all
scripts in lexical order (see [filepath.Walk](https://godoc.org/path/filepath#Walk)).

### Example directory structure:

```
sql
├── migrations
│   ├── 1_usertable.up.sql
│   ├── 1_usertable.down.sql
│   ├── 2_other-table.up.sql
│   ├── 2_other-table.down.sql
│   ├── 3_usertable-addColumn.up.sql
│   └── 3_usertable-addColumn.down.sql
└── seeds
    ├── users.sql
    └── other-data.sql
```
