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

    // Combine ListenIP and HttpsPort to form the full listen address
    fullListenAddr := listenAddr + ":" + cfg.HttpsPort

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        CheckRoleHandler(w, cfg)
    })

    log.Printf("Starting server on %s\n", fullListenAddr)
    if err := http.ListenAndServeTLS(fullListenAddr, cfg.CertFile, cfg.KeyFile, nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
