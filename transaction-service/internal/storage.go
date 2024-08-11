package internal

import (
    "database/sql"
    _ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() error {
    var err error
    connStr := "user=rohit dbname=locodb sslmode=disable"
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        return err
    }
    return db.Ping()
}

func GetDB() *sql.DB {
    return db
}
