package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
)

const configFileName = ".rss-aggregator.config.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUsername string `json:"current_username"`
}

func Read() Config {
	configPath, err := getConfigFilePath()
	if err != nil {
		log.Fatalf("getting config file path failed: %v", err)
	}
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("opening config file failed: %v", err)
	}
	defer file.Close()
	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		log.Fatalf("decoding config file failed: %v", err)
	}
	return cfg
}

func (cfg *Config) SetUser(username string) error {
	if username == "" {
		return fmt.Errorf("username must not be empty")
	}
	cfg.CurrentUsername = username
	return write(cfg)
}

func write(cfg *Config) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		log.Fatalf("getting config file path failed: %v", err)
	}
	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := json.NewEncoder(file).Encode(&cfg); err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home dir failed: %v", err)
	}
	configPath := path.Join(home, configFileName)
	return configPath, nil
}
