// logging.go
package main

import (
    "fmt"
    "io"
    "log"
    "log/syslog"
    "os"
)

func SetupLogging(logFilePath string, foreground bool) *os.File {
    // Debug statement to print the log file path
    fmt.Printf("Debug: Attempting to set up logging with LogFilePath: '%s'\n", logFilePath)

    var logOutput io.Writer
    var logFile *os.File

    if logFilePath == "syslog" {
        logwriter, err := syslog.New(syslog.LOG_NOTICE, "pgrolecheck")
        if err != nil {
            fmt.Println("Failed to initialize syslog:", err)
            os.Exit(1)
        }
        logOutput = logwriter
    } else if logFilePath != "" {
        var err error
        logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            fmt.Printf("Failed to open log file: '%s', Error: %v\n", logFilePath, err)
            os.Exit(1)
        }
        if foreground {
            logOutput = io.MultiWriter(os.Stdout, logFile)
        } else {
            logOutput = logFile
        }
    } else {
        // If logFilePath is not defined, fallback to stdout
        fmt.Println("Debug: No LogFilePath specified, defaulting to stdout")
        logOutput = os.Stdout
    }

    log.SetOutput(logOutput)
    return logFile
}
