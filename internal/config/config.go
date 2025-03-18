package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	conf := Config{}
	fp, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("error: %v", err)
	}
	content, err := os.ReadFile(fp) //filename hardcoded for now
	if err != nil {
		return conf, fmt.Errorf("error: %v", err)
	}
	err = json.Unmarshal(content, &conf)
	if err != nil {
		return Config{}, err
	}
	return conf, nil
}

func (c Config) write() error {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	fp, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	return os.WriteFile(fp, jsonData, 0666)
}

// untested
func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	return c.write()
}

func getConfigFilePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homedir + "/" + configFileName, nil
}
