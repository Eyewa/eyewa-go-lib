# eyewa-go-lib
Shared Go Lib for Eyewa's microservices.

# db
This package provides an abstraction layer for a couple of database drivers by means of creating clients. Currently this package only supports MySQL, SQLite & Postgres. In the future more drivers can be added if the need arises.

The clients implement exponential backoffs in the following instance:
- when on the start of a service, the connection to the db cannot be established.

# How to use
The following env variables can be injected in order to use this pkg:

```go
"DB_DRIVER" // required - available options include mysql, postgres, sqlite

// required if driver is either mysql, postgres
"DB_HOST"
"DB_USER"
"DB_PASSWORD"
"DB_DATABASE"

// required - if sqlite is the chosen driver
// https://www.sqlite.org/inmemorydb.html
"DB_PATH" // options are ":memory:" or file path.

// optional and only for postgres. option is  "disable" or left blank
// if blank ssl mode will be in used.
"DB_SSL_MODE"
```
### Connecting to a Database

```go
	// JUST FOR DEMO - these should be injected
	vars = map[string]string{
		"DB_DRIVER":   "postgres",
		"DB_HOST":     "localhost",
		"DB_USER":     "admin",
		"DB_PORT":     "5432",
		"DB_PASSWORD": "secret",
		"DB_DATABASE": "catalogconsumer",
		"DB_SSL_MODE": "disable",
	}
	for e, v := range vars {
		os.Setenv(e, v)
	}

	dbclient, err := OpenConnection()
	if err != nil {
		log.Error("Failed to connect to DB", zap.Error(err))
	}
```

### Requiring multiple DB clients in one application/service.

```go
	// connect to postgres
	pCfg := Config{
		Database: RDMS{
			Name:     "catalogindexer",
			Host:     "localhost",
			User:     "admin",
			Port:     "5432",
			Password: "secret",
		},
	}

	postgresClient := NewPostgresClientFromConfig(pCfg)
	client, err := postgresClient.OpenConnection()
	if err != nil {
		log.Error("Failed to connect to DB", zap.Error(err))
	}

	// connect to mysql - can be any other db client like mongo/redis etc.
	mCfg := Config{
		Database: RDMS{
			Name:     "catalogconsumer",
			Host:     "localhost",
			User:     "admin007",
			Port:     "3306",
			Password: "mystic",
		},
	}

	mysqlClient := NewMySQLClientFromConfig(mCfg)
	client, err := mysqlClient.OpenConnection()
	if err != nil {
		log.Error("Failed to connect to DB", zap.Error(err))
	}

```

### CRUD Operations
TBD