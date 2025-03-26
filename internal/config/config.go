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

type state struct {
	conf *Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("missing argument for login")
	}
	if len(cmd.args) > 1 {
		return fmt.Errorf("too many arguments for login")
	}
	s.conf.SetUser(cmd.args[0])
	fmt.Printf("Username set to '%s'\n", cmd.args[0])
	return nil
}

type commands struct {
	funcs map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.funcs[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	funcToRun, ok := c.funcs[cmd.name]
	if !ok {
		return fmt.Errorf("function for command '%s' does not exist", cmd.name)
	}
	funcToRun(s, cmd)
	return nil
}
