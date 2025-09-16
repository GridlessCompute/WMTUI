package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Sites []Site `json:"sites"`
}

type Site struct {
	Name    string `json:"name"`
	IPRange string `json:"iprange"`
}

func NewConfig(dir string) (Config, error) {
	conf := Config{
		Sites: []Site{{Name: "Example Site", IPRange: "192.168.10.0/24"}},
	}

	file, err := json.MarshalIndent(conf, "", " ")
	if err != nil {
		return Config{}, fmt.Errorf("issue marshalling default config: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(dir), 0755); err != nil {
		return Config{}, fmt.Errorf("issue creating config directory: %v", err)
	}

	if err := os.WriteFile(dir, file, 0644); err != nil {
		return Config{}, fmt.Errorf("issue writing default config file: %v", err)
	}

	return conf, nil
}

func GetConf() (Config, error) {
	var conf Config

	confdir, err := os.UserConfigDir()
	if err != nil {
		return Config{}, fmt.Errorf("issue getting user config directory: %v", err)
	}

	configPath := filepath.Join(confdir, "wmtui", "config.json")

	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		return NewConfig(configPath)
	}

	confFile, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("issue reading config file from %s: %v", configPath, err)
	}

	if err = json.Unmarshal(confFile, &conf); err != nil {
		return Config{}, fmt.Errorf("issue trying to unmarshal config: %v", err)
	}

	return conf, nil
}
