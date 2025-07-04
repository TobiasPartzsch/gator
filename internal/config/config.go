package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	var err error
	var cfgPath string
	if cfgPath, err = getConfigFilePath(); err != nil {
		return Config{}, fmt.Errorf(
			"error trying to calculate path to config file: %w",
			err,
		)
	}

	var cfgData []byte
	if cfgData, err = os.ReadFile(cfgPath); err != nil {
		return Config{}, fmt.Errorf(
			"error trying to read config file %s: %w",
			cfgPath,
			err,
		)
	}

	var cfg Config
	if err = json.Unmarshal(cfgData, &cfg); err != nil {
		return Config{}, fmt.Errorf(
			"error trying to unmarshal content of config file %s: %w",
			cfgPath,
			err,
		)
	}
	return cfg, nil
}

func (c *Config) SetUser(userName string) error {
	var err error

	c.CurrentUserName = userName
	err = write(*c)
	if err != nil {
		return fmt.Errorf(
			"error trying to set a new user: %w",
			err,
		)
	}
	return nil
}

func getConfigFilePath() (string, error) {
	var err error
	var homeDir string

	if homeDir, err = os.UserHomeDir(); err != nil {
		return "", fmt.Errorf(
			"error trying to figure out home directory: %w",
			err,
		)
	}
	cfgPath := homeDir + "/" + configFileName
	return cfgPath, nil
}

func write(cfg Config) error {
	var err error
	var cfgPath string

	if cfgPath, err = getConfigFilePath(); err != nil {
		return fmt.Errorf(
			"error trying to calculate path to config file: %w",
			err,
		)
	}

	var cfgData []byte
	if cfgData, err = json.Marshal(cfg); err != nil {
		return fmt.Errorf(
			"error trying to marshal config data: %w",
			err,
		)
	}

	if err = os.WriteFile(cfgPath, cfgData, 0644); err != nil {
		return fmt.Errorf(
			"error trying to write config data to config file: %w",
			err,
		)
	}

	return nil
}
