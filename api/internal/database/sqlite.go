package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"

	"github.com/edgejay/pify-player/api/internal/utils"
)

type SQLiteDB struct {
	SQL *sql.DB
	Bun *bun.DB
}

var db *SQLiteDB = nil

func initSQLiteDB() {
	if db != nil {
		log.Println("Database already initialised, skipping setup")
	}

	db = &SQLiteDB{}

	dbFile := utils.GetDBFilename()
	if dbFile == "" {
		dbFile = "./database/db.sqlite3"
	}

	log.Printf("Initiailising database at %s\n", dbFile)

	// Creates a SQLite database in file (determined by DB_FILE env variable)
	sqldb, err := sql.Open(
		sqliteshim.DriverName(),
		fmt.Sprintf("file:%s?cache=shared", dbFile),
	)

	if err != nil {
		panic(err)
	}

	db.SQL = sqldb
	db.Bun = bun.NewDB(db.SQL, sqlitedialect.New())
	db.Bun.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(false),
		bundebug.FromEnv(),
	))

	log.Printf("Database setup complete. file created in %s\n", dbFile)
}

func GetSQLiteDB() *SQLiteDB {
	initSQLiteDB()
	return db
}

func (db *SQLiteDB) Ping() error {
	return db.Bun.Ping()
}
