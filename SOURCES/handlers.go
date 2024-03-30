// handlers.go
package main

import (
    "encoding/csv"
    "encoding/json"
    "net/http"
    "log"
    // other imports as necessary
)

type DBStatus struct {
    Status  string `json:"status"`
    Message string `json:"message,omitempty"` // Include only if there's an error
}

func CheckRoleHandler(w http.ResponseWriter, r *http.Request) {
    // Parse the query parameters from the request URL
    query := r.URL.Query()
    dbname := query.Get("dbname")

    // Initialize a map to store the database statuses
    results := make(map[string]DBStatus)

    // Iterate through each database configuration
    for _, dbCfg := range cfg.Databases {
        // If a dbname parameter is specified and it doesn't match the current database name, skip it
        if dbname != "" && dbname != dbCfg.Name {
            continue
        }

        // Check the role of the current database
        status, err := checkDatabaseRole(dbCfg)

        // Create a DBStatus struct for the current database
        dbStatus := DBStatus{
            Status: status,
        }

        // Handle any errors that occurred during the role check
        if err != nil {
            log.Printf("Error checking database role for %s: %v", dbCfg.Name, err)
            dbStatus.Status = "error"
            dbStatus.Message = err.Error()
        }

        // Add the database status to the results map
        results[dbCfg.Name] = dbStatus
    }

    // Generate the appropriate response based on the output format specified in the configuration
    if cfg.OutputFormat == "csv" {
        w.Header().Set("Content-Type", "text/csv")
        writer := csv.NewWriter(w)

        // Write CSV header
        if err := writer.Write([]string{"Name", "Status", "Message"}); err != nil {
            log.Println("Error writing CSV header:", err)
            return
        }

        // Write data rows
        for name, status := range results {
            if err := writer.Write([]string{name, status.Status, status.Message}); err != nil {
                log.Println("Error writing CSV record:", err)
                return
            }
        }

        writer.Flush() // Make sure all buffered data is sent to the client

        if err := writer.Error(); err != nil {
            log.Println("Error flushing CSV data:", err)
            http.Error(w, "Failed to generate CSV output", http.StatusInternalServerError)
        }
    } else {
        // Generate and respond with JSON
        w.Header().Set("Content-Type", "application/json")
        jsonResponse, err := json.Marshal(results)
        if err != nil {
            log.Println("Error creating JSON response:", err)
            http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
            return
        }
        w.Write(jsonResponse)
    }
}
