package common

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Configuration stores setting values for tvshow service
type Configuration struct {
	Port string `json:"port"`

	LogFilename   string `json:"log_filename"`
	LogMaxSize    int    `json:"log_max_size"`
	LogMaxBackups int    `json:"log_max_backups"`
	LogMaxAge     int    `json:"log_max_age"`

	DbAddr     string `json:"db_addr"`
	DbName     string `json:"db_name"`
	DbUsername string `json:"db_username"`
	DbPassword string `json:"db_password"`
}

// Config shares the global configuration
var (
	Config *Configuration
)

// Collections names if user wants to use mongo
const (
	CollectionTVShow = "tvshows"
)

// LoadConfig loads configuration from the config file
func LoadConfig() error {
	file, err := os.Open("config/config.json")
	if err != nil {
		return err
	}

	Config = new(Configuration)
	if err := json.NewDecoder(file).Decode(Config); err != nil {
		return err
	}

	// set up logging
	log.SetOutput(&lumberjack.Logger{
		Filename:   Config.LogFilename,
		MaxSize:    Config.LogMaxSize,
		MaxBackups: Config.LogMaxBackups,
		MaxAge:     Config.LogMaxAge,
	})
	log.SetLevel(log.DebugLevel)

	// output in JSON format
	log.SetFormatter(&log.JSONFormatter{})

	return nil
}
