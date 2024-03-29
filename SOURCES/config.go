// config.go
package main

import (
	"flag"
	"gopkg.in/ini.v1"
)

var (
    cfg Configuration
    ConfigPath = flag.String("config", "/etc/pgrolecheck/pgrolecheck.conf", "Path to configuration file")
    Foreground = flag.Bool("f", false, "Run in foreground and log to STDOUT")
)

var configMap = map[string]*struct {
    flagValue *string
    cfgField  *string
}{
    "dbname":      {flag.String("dbname", "", "Database name"), &cfg.DbName},
    "user":        {flag.String("user", "", "Database user"), &cfg.User},
    "password":    {flag.String("password", "", "Database password"), &cfg.Password},
    "host":        {flag.String("host", "", "Database host"), &cfg.Host},
    "port":        {flag.String("port", "", "Database port"), &cfg.Port},
    "sslmode":     {flag.String("sslmode", "", "Database SSL mode"), &cfg.SslMode},
    "listen_ip":   {flag.String("listenip", "", "Server listen IP"), &cfg.ListenIP},
    "https_port":  {flag.String("httpsport", "", "HTTPS port"), &cfg.HttpsPort},
    "cert_file":   {flag.String("certfile", "", "SSL certificate file"), &cfg.CertFile},
    "key_file":    {flag.String("keyfile", "", "SSL key file"), &cfg.KeyFile},
    "log_file":    {flag.String("logfile", "", "Log file path"), &cfg.LogFilePath},
}

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

func LoadConfigurationFromFile(path string) error {
    conf, err := ini.Load(path)
    if err != nil {
        return err
    }

    // Correctly load settings based on the actual sections and keys in your INI file
    cfg.DbName = conf.Section("database").Key("dbname").String()
    cfg.User = conf.Section("database").Key("user").String()
    cfg.Password = conf.Section("database").Key("password").String()
    cfg.SslMode = conf.Section("database").Key("sslmode").String()
    cfg.Port = conf.Section("database").Key("port").String()
    // Continue for the rest of your database settings

    cfg.ListenIP = conf.Section("server").Key("listen_ip").String()
    cfg.HttpsPort = conf.Section("server").Key("https_port").String()
    cfg.CertFile = conf.Section("server").Key("cert_file").String()
    cfg.KeyFile = conf.Section("server").Key("key_file").String()
    // Assuming these are the correct keys under the [server] section

    cfg.LogFilePath = conf.Section("logging").Key("log_file").String()
    // Assuming log_file is under the [logging] section

    return nil
}

func OverrideConfigurationWithFlags() {
    // Assume flag.Parse() has already been called
    for _, value := range configMap {
        if *value.flagValue != "" {
            *value.cfgField = *value.flagValue
        }
    }
}
