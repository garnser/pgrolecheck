// main.go
package main

import (
    "flag"
    "log"
    "fmt"
    "net/http"
    "os"
    "encoding/json"
    "github.com/coreos/go-systemd/daemon"
)

func main() {

    flag.Parse()
    fmt.Printf("Running in foreground: %v\n", Foreground)

    if err := LoadConfigurationFromFile(*ConfigPath); err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // After loading configurations from the file, also parse dbConfigsJSON and update cfg.Databases
    for _, dbConfigJSON := range dbConfigsJSON {
        var dbConfig DBConfig
        if err := json.Unmarshal([]byte(dbConfigJSON), &dbConfig); err != nil {
            log.Fatalf("Failed to parse database configuration from JSON: %v", err)
        }
        cfg.Databases = append(cfg.Databases, dbConfig)
    }

    logFile := SetupLogging(cfg.LogFilePath, Foreground)
    if logFile != nil {
        defer logFile.Close()
    }

    listenAddr := cfg.ListenIP
    if listenAddr == "*" || listenAddr == "0.0.0.0" {
        listenAddr = "0.0.0.0"
    }

    fullListenAddr := listenAddr + ":" + cfg.HttpPort

    http.Handle("/", IPWhitelistMiddleware(TokenAuthMiddleware(http.HandlerFunc(CheckRoleHandler))))

    // Adjust log message to reflect multiple databases
    log.Printf("Monitoring %d databases for role status", len(cfg.Databases))

    invocationID := os.Getenv("INVOCATION_ID")
    if invocationID != "" {
       _, err := daemon.SdNotify(false, daemon.SdNotifyReady)
       if err != nil {
          log.Fatalf("Failed to notify systemd: %v\n", err)
       }
    }

    if cfg.UseSSL == "true" {
        _, certErr := os.Stat(cfg.CertFile)
        _, keyErr := os.Stat(cfg.KeyFile)
        if os.IsNotExist(certErr) || os.IsNotExist(keyErr) {
           log.Fatalf("Certificate or key file does not exist.")
       	} else if certErr != nil || keyErr != nil {
       	   log.Fatalf("Error accessing certificate or key file: %v %v", certErr, keyErr)
   	    }

        log.Printf("Starting HTTPS server on %s", fullListenAddr)
        if err := http.ListenAndServeTLS(fullListenAddr, cfg.CertFile, cfg.KeyFile, nil); err != nil {
            log.Fatalf("Failed to start HTTPS server: %v", err)
        }
    } else {
        log.Printf("Starting HTTP server on %s", fullListenAddr)
        if err := http.ListenAndServe(fullListenAddr, nil); err != nil {
            log.Fatalf("Failed to start HTTP server: %v", err)
        }
    }
}
