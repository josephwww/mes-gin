package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type Config struct {
	AppName    string `json:"app_name"`
	ServerPort string `json:"server_port"`
	Mode       string `json:"mode"` // Running mode: prod, test, dev
	Database   struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
	} `json:"database"`
	Redis struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"redis"`
}

var AppConfig Config

// Load function reads the configuration from settings.json and loads it into AppConfig
func Load() error {
	configFile, err := os.Open("config/settings.json")
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
		return err
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		log.Fatalf("Error decoding config file: %v", err)
		return err
	}

	// Set Gin mode based on the configuration
	switch AppConfig.Mode {
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	log.Printf("Configuration loaded successfully: %+v", AppConfig)
	return nil
}
