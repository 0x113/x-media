package common

import (
	"encoding/json"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Configuration stores setting values for user service
type Configuration struct {
	Port string `json:"port"`

	LogFilename   string `json:"log_filename"`
	LogMaxSize    int    `json:"log_max_size"`
	LogMaxBackups int    `json:"log_max_backups"`
	LogMaxAge     int    `json:"log_max_age"`

	AccessSecret  string `json:"access_secret"`
	RefreshSecret string `json:"refresh_secret"`
}

// Config shares the global configuration
var (
	Config *Configuration
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

	// set up logging TODO: different foramatter for Stdout
	multiWriter := io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   Config.LogFilename,
		MaxSize:    Config.LogMaxSize,
		MaxBackups: Config.LogMaxBackups,
		MaxAge:     Config.LogMaxAge,
	})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(multiWriter)

	// output in JSON format
	log.SetFormatter(&log.JSONFormatter{})

	return nil
}
