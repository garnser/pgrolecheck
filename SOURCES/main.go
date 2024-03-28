// main.go
package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"log/syslog"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/ini.v1"
)

var foreground = flag.Bool("f", false, "Run in foreground and log to STDOUT")
var logFile *os.File

// Configuration holds all configurations for the application
type Configuration struct {
	DbName      string
	User        string
	Password    string
	Host        string
	Port        string
	SslMode     string
	ListenIP    string
	HttpsPort   string
	CertFile    string
	KeyFile     string
	LogFilePath string
}

// LoadConfiguration loads configuration from the given file path
func LoadConfiguration(path string) (Configuration, error) {
	var cfg Configuration
	conf, err := ini.Load(path)
	if err != nil {
		return cfg, err
	}

	cfg.DbName = conf.Section("database").Key("dbname").String()
	cfg.User = conf.Section("database").Key("user").String()
	cfg.Password = conf.Section("database").Key("password").String()
	cfg.Host = conf.Section("database").Key("host").String()
	cfg.Port = conf.Section("database").Key("port").String()
	cfg.SslMode = conf.Section("database").Key("sslmode").String()
	cfg.ListenIP = conf.Section("server").Key("listen_ip").String()
	cfg.HttpsPort = conf.Section("server").Key("https_port").String()
	cfg.CertFile = conf.Section("server").Key("cert_file").String()
	cfg.KeyFile = conf.Section("server").Key("key_file").String()
	cfg.LogFilePath = conf.Section("logging").Key("log_file").String()

	return cfg, nil
}

// SetupLogging configures logging based on the configuration and flags
func SetupLogging(logFilePath string, foreground bool) {
    var logOutput io.Writer
    
    if logFilePath == "syslog" {
        logwriter, err := syslog.New(syslog.LOG_NOTICE, "pgrolecheck")
        if err != nil {
            fmt.Println("Failed to initialize syslog:", err)
            os.Exit(1) // Exiting if syslog cannot be initialized
        }
        logOutput = logwriter
    } else {
        var err error
        logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            fmt.Println("Failed to open log file:", err)
            os.Exit(1) // Exiting if the log file cannot be opened
        }
        
        if foreground {
            // In foreground mode, log to both STDOUT and the log file
            logOutput = io.MultiWriter(os.Stdout, logFile)
        } else {
            // Otherwise, log only to the file
            logOutput = logFile
        }
    }
    
    log.SetOutput(logOutput)
}

func checkRoleHandler(w http.ResponseWriter, cfg Configuration) {
	connStr := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=%s", cfg.DbName, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.SslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
	       log.Println("Error connecting to database:", err)
	       respondWithError(w, "notok")
	       return
	}
	defer db.Close()

	var isInRecovery bool
	err = db.QueryRow("SELECT pg_is_in_recovery();").Scan(&isInRecovery)
	if err != nil {
	       log.Println("Error querying database:", err)
	       respondWithError(w, "notok")
	       return
	}

	status := "primary"
	if isInRecovery {
		status = "notprimary"
	}

	respondWithStatus(w, status)
	log.Printf("Web-server accessed, HTTP return-code: %d, status: %s\n", http.StatusOK, status)
}

func respondWithError(w http.ResponseWriter, status string) {
   respondWithStatus(w, status)
}

func respondWithStatus(w http.ResponseWriter, status string) {
   response := map[string]string{"status": status}
   jsonResponse, err := json.Marshal(response)
   if err != nil {
       log.Println("Error creating JSON response:", err)
       http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
       return
   }
   
   w.Header().Set("Content-Type", "application/json")
   w.WriteHeader(http.StatusOK) // Consider adjusting the status code based on the error
   w.Write(jsonResponse)
}

func main() {
	flag.Parse()

	cfg, err := LoadConfiguration("/etc/pgrolecheck/pgrolecheck.conf")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	SetupLogging(cfg.LogFilePath, *foreground)

	defer func() {
       	      if logFile != nil {
              	 logFile.Close()
	      }
   	}()

	// Determine the IP address to listen on
	listenAddr := cfg.ListenIP
	if listenAddr == "*" || listenAddr == "0.0.0.0" {
	   	listenAddr = "0.0.0.0"
	}

	// Combine ListenIP and HttpsPort to form the full listen address
	fullListenAddr := listenAddr + ":" + cfg.HttpsPort

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		checkRoleHandler(w, cfg)
	})

	log.Printf("Starting server on %s\n", fullListenAddr)
	err = http.ListenAndServeTLS(fullListenAddr, cfg.CertFile, cfg.KeyFile, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
