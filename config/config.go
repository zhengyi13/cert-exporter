package config

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v3"
)

// HostPort simply separates out the two needed components for getting a cert into labelled fields.
type HostPort struct {
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
}

func (hp HostPort) String() string {
	return fmt.Sprintf("%s:%d", hp.Hostname, hp.Port)
}

type Config []HostPort

// GetConfig takes a filename, and returns either a functional config, or an error.
func GetConfig(fn string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, fmt.Errorf("Cannot read config file: %w", err)
	}
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse config: %w", err)
	}
	return &config, nil
}
