package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	var cfg Config

	cfgPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	cfgFile, err := os.ReadFile(cfgPath)
	if err != nil {
		return Config{}, fmt.Errorf("I/O error: %w", err)
	}

	if err := json.Unmarshal(cfgFile, &cfg); err != nil {
		return Config{}, fmt.Errorf("json unmarshalling error: %w", err)
	}

	return cfg, nil
}

func InitIfNotExists() error {
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	_, err = os.Stat(cfgPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("I/O error: %w", err)
	} else if err == nil {
		return nil
	}

	fmt.Println("WARNING: writing default config file")

	cfg := Config{
		DBURL: "",
		CurrentUserName: "",
	}

	return cfg.write()
}

func getConfigFilePath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("environment error: %w", err)
	}
	return dir + "/" + configFileName, nil
}

func (c *Config) write() error {
	cfg, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("json marshalling error: %w", err)
	}

	cfgPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(cfgPath, cfg, 0644)
	if err != nil {
		return fmt.Errorf("I/O error: %w", err)
	}

	return nil
}

func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name

	return c.write()
}

func (c *Config) SetDBURL(dbURL string) error {
	c.DBURL = dbURL

	return c.write()
}
