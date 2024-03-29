package main

import (
    "flag"
    "log"
    "net/http"
)

func main() {
    flag.Parse()

    // Assuming LoadConfigurationFromFile and OverrideConfigurationWithFlags
    // are correctly defined in config.go and modify the cfg global variable.
    if err := LoadConfigurationFromFile(*ConfigPath); err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    OverrideConfigurationWithFlags()

    logFile := SetupLogging(cfg.LogFilePath, *Foreground)
    if logFile != nil {
        defer logFile.Close()
    }

    // Determine the IP address to listen on
    listenAddr := cfg.ListenIP
    if listenAddr == "*" || listenAddr == "0.0.0.0" {
        listenAddr = "0.0.0.0"
    }

    // Combine ListenIP and HttpPort to form the full listen address
    fullListenAddr := listenAddr + ":" + cfg.HttpPort

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        CheckRoleHandler(w, cfg)
    })

    if cfg.UseSSL == "true" {
        // Start HTTPS server
        log.Printf("Starting HTTPS server on %s\n", fullListenAddr)
        if err := http.ListenAndServeTLS(fullListenAddr, cfg.CertFile, cfg.KeyFile, nil); err != nil {
            log.Fatalf("Failed to start HTTPS server: %v", err)
        }
    } else {
        // Start HTTP server
        log.Printf("Starting HTTP server on %s\n", fullListenAddr)
        if err := http.ListenAndServe(fullListenAddr, nil); err != nil {
            log.Fatalf("Failed to start HTTP server: %v", err)
        }
    }
}
