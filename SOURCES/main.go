// main.go
package main

import (
    "flag"
    "log"
    "net/http"
)

func main() {

    flag.Parse()

    if err := LoadConfigurationFromFile(*ConfigPath); err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
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

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        CheckRoleHandler(w, r)
    })

    // Adjust log message to reflect multiple databases
    log.Printf("Monitoring %d databases for role status", len(cfg.Databases))

    if cfg.UseSSL == "true" {
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
