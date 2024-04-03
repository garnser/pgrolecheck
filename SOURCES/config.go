// config.go
package main

import (
    "flag"
    "fmt"
    "gopkg.in/ini.v1"
    "reflect"
    "strings"
    "strconv"
)

var (
    cfg          Configuration
    ConfigPath   = flag.String("config", "/etc/pgrolecheck/pgrolecheck.conf", "Path to configuration file")
    dbConfigsJSON JSONConfigs
    Foreground   bool
)

type Configuration struct {
    ListenIP     string `config:"listen_ip" section:"server" default:"" description:"The IP address on which the server will listen for incoming requests. Leave blank to listen on all interfaces."`
    UseSSL       bool `config:"use_ssl" section:"server" default:"false" description:"Determines whether SSL is enabled. Set to true to enable SSL."`
    HttpPort     string `config:"http_port" section:"server" default:"8080" description:"The port on which the server will listen for HTTP requests."`
    CertFile     string `config:"cert_file" section:"server" default:"" description:"The file path to the SSL certificate. Required if SSL is enabled."`
    KeyFile      string `config:"key_file" section:"server" default:"" description:"The file path to the SSL certificate key. Required if SSL is enabled."`
    LogFilePath  string `config:"log_file" section:"logging" default:"/var/log/pgrolecheck.log" description:"The file path where logs will be written. Use 'syslog' for system log or leave blank for stdout."`
    OutputFormat string `config:"output_format" section:"server" default:"json" description:"The format of the response returned by the server. Options are 'json', 'csv', or 'simple' text."`
    EnableAccessLog bool   `config:"enable_access_log" section:"logging" default:"true" description:"Enables or disables HTTP access logging. Useful for monitoring and debugging."`
    Databases    []DBConfig // This field doesn't directly correspond to a command-line argument or ini setting.
    IPWhitelist  []string `config:"ip_whitelist" section:"security" default:"" description:"A comma-separated list of IP addresses or CIDR ranges that are allowed to access the server. Leave blank to allow all."`
    AuthToken    string `config:"auth_token" section:"security" default:"" description:"The token that clients must provide for authentication. Leave blank to disable token authentication."`
}

type DBConfig struct {
    Name     string `json:"name"`
    DbName   string `json:"dbname"`
    User     string `json:"user"`
    Password string `json:"password"`
    Host     string `json:"host"`
    Port     string `json:"port"`
    SslMode  string `json:"sslmode"`
}

type JSONConfigs []string

func (j *JSONConfigs) String() string {
    return strings.Join(*j, ",")
}

func (j *JSONConfigs) Set(value string) error {
    *j = append(*j, value)
    return nil
}

func init() {
    flag.BoolVar(&Foreground, "f", false, "Run in foreground and log to STDOUT")
    defineConfigurationFlags()
    flag.Var(&dbConfigsJSON, "dbconfig", "Database configuration as a JSON string")
}

func defineConfigurationFlags() {
    t := reflect.TypeOf(cfg)
    v := reflect.ValueOf(&cfg).Elem()

    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        key, hasKey := field.Tag.Lookup("config")
        description, hasDescription := field.Tag.Lookup("description")
        defaultValue, hasDefault := field.Tag.Lookup("default")

        if !hasKey || !hasDefault {
            continue
        }

        if !hasDescription {
            description = fmt.Sprintf("No description for %s", key) // Fallback description
        }

        switch field.Type.Kind() {
        case reflect.String:
            flag.StringVar(v.Field(i).Addr().Interface().(*string), key, defaultValue, description)
        case reflect.Bool:
            defaultValueBool, _ := strconv.ParseBool(defaultValue) // safely ignore error, default is provided
            flag.BoolVar(v.Field(i).Addr().Interface().(*bool), key, defaultValueBool, description)
        }
    }
}

func LoadConfigurationFromFile(path string) error {
    configFile, err := ini.Load(path)
    if err != nil {
        return fmt.Errorf("failed to load config file: %w", err)
    }

    // First, load server and logging configurations as before
    t := reflect.TypeOf(cfg)
    v := reflect.ValueOf(&cfg).Elem()
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        key, hasKey := field.Tag.Lookup("config")
        section, hasSection := field.Tag.Lookup("section")
        if !hasKey || !hasSection {
            continue
        }

        if sectionValue := configFile.Section(section).Key(key).String(); sectionValue != "" {
            if field.Type.Kind() == reflect.String {
                v.Field(i).SetString(sectionValue)
            }
        }
    }

    // Dynamically handle database configurations
    cfg.Databases = []DBConfig{}
    dbConfigType := reflect.TypeOf(DBConfig{})

    for _, section := range configFile.Sections() {
        if strings.HasPrefix(section.Name(), "database_") {
            dbCfg := DBConfig{}
            dbCfgValue := reflect.ValueOf(&dbCfg).Elem()

            dbName := section.Name()[len("database_"):]

            for i := 0; i < dbConfigType.NumField(); i++ {
                field := dbConfigType.Field(i)
                key := field.Tag.Get("json")

                if key == "" {
                    key = strings.ToLower(field.Name) // Fallback to using the struct field name
                }

                iniValue := section.Key(key).String()
                if iniValue != "" {
                    dbCfgValue.Field(i).SetString(iniValue)
                }
            }

            // Assign the extracted database name to the Name field of DBConfig
            dbCfg.Name = dbName

            cfg.Databases = append(cfg.Databases, dbCfg)
        }
    }

    ipWhitelistStr := configFile.Section("security").Key("ip_whitelist").String()
    if ipWhitelistStr != "" {
        cfg.IPWhitelist = strings.Split(ipWhitelistStr, ",")
    }

    return nil
}
