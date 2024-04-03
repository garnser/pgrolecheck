// config.go
package main

import (
    "flag"
    "fmt"
    "gopkg.in/ini.v1"
    "reflect"
    "strings"
)

var (
    cfg          Configuration
    ConfigPath   = flag.String("config", "/etc/pgrolecheck/pgrolecheck.conf", "Path to configuration file")
    dbConfigsJSON JSONConfigs
    Foreground   bool
)

type Configuration struct {
    ListenIP     string `config:"listen_ip" section:"server" default:""`
    UseSSL       bool `config:"use_ssl" section:"server" default:"false"`
    HttpPort     string `config:"http_port" section:"server" default:"8080"`
    CertFile     string `config:"cert_file" section:"server" default:""`
    KeyFile      string `config:"key_file" section:"server" default:""`
    LogFilePath  string `config:"log_file" section:"logging" default:"/var/log/pgrolecheck.log"`
    OutputFormat string `config:"output_format" section:"server" default:"json"`
    EnableAccessLog bool   `config:"enable_access_log" section:"logging" default:"true"`
    Databases    []DBConfig
    IPWhitelist  []string `config:"ip_whitelist" section:"security" default:""`
    AuthToken    string `config:"auth_token" section:"security" default:""`
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
        defaultValue, hasDefault := field.Tag.Lookup("default")

        if !hasKey || !hasDefault {
            continue
        }

        description := fmt.Sprintf("Description for %s", key)
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
