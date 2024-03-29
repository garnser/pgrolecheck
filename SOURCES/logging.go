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
    var logOutput io.Writer
    var logFile *os.File

    if logFilePath == "syslog" {
        logwriter, err := syslog.New(syslog.LOG_NOTICE, "pgrolecheck")
        if err != nil {
            fmt.Println("Failed to initialize syslog:", err)
            os.Exit(1)
        }
        logOutput = logwriter
    } else {
        var err error
        logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            fmt.Println("Failed to open log file:", err)
            os.Exit(1)
        }

        if foreground {
            logOutput = io.MultiWriter(os.Stdout, logFile)
        } else {
            logOutput = logFile
        }
    }

    log.SetOutput(logOutput)
    return logFile
}

