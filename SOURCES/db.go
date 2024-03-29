// db.go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// checkDatabaseRole connects to the database and checks if it's in recovery mode.
func checkDatabaseRole(cfg Configuration) (string, error) {
	connStr := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=%s", cfg.DbName, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.SslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return "", fmt.Errorf("error connecting to database: %v", err)
	}
	defer db.Close()

	var isInRecovery bool
	err = db.QueryRow("SELECT pg_is_in_recovery();").Scan(&isInRecovery)
	if err != nil {
		return "", fmt.Errorf("error querying database: %v", err)
	}

	if isInRecovery {
		return "notprimary", nil
	}
	return "primary", nil
}
