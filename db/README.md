# eyewa-go-lib
Shared Go Lib for Eyewa's microservices.

# db
This package provides an abstraction layer for `https://github.com/go-gorm/gorm` ORM pkg in having DB clients for several database drivers. Currently this package only supports MySQL, SQLite & Postgres. In the future more drivers can be added if need arises.

# How to use
The following env variables can be injected in order to use this pkg:

```go
"DB_DRIVER" // required - available options include mysql, postgres

// required if driver is either mysql, postgres
"DB_HOST"
"DB_USER"
"DB_PASSWORD"
"DB_DATABASE"

// required - if sqlite is the chosen driver
// https://www.sqlite.org/inmemorydb.html
"DB_PATH" // options are ":memory:" or file path.

// optional and only for postgres. option is  "disabled" or left blank
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
		"DB_SSL_MODE": "disabled",
	}
	for e, v := range vars {
		os.Setenv(e, v)
	}

	dbclient, err := OpenConnection()
	if err != nil {
		log.Error("Failed to connect to DB", zap.Error(err))
	}
```

### CRUD operations
TBD