package engine

import (
	"os"

	"gopkg.in/yaml.v3"
)

type DatabaseProps struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type ConfigModel struct {
	Database  DatabaseProps  `yaml:"database"`
	Variables map[string]any `yaml:"variables"`
}

func ReadConfig(filename string) (*ConfigModel, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := new(ConfigModel)
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
