// db.go
package main

import (
    "database/sql"
    "log"
    "fmt"
    _ "github.com/lib/pq" // PostgreSQL driver
)

// checkDatabaseRole connects to the database using the provided DBConfig and checks if it's in recovery mode.
func checkDatabaseRole(dbCfg DBConfig) (string, error) {
    // Use dbCfg to access the database configuration directly
    connStr := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=%s",
        dbCfg.DbName, dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.SslMode)
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Printf("error connecting to database: %v", err)
        return "", err
    }
    defer db.Close()

    var isInRecovery bool
    err = db.QueryRow("SELECT pg_is_in_recovery();").Scan(&isInRecovery)
    if err != nil {
        log.Printf("error querying database: %v", err)
        return "", err
    }

    if isInRecovery {
        return "replica", nil
    }
    return "primary", nil
}
