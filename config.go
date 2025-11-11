package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Spirit struct {
		Token  string `yaml:"token"`
		ClubID string `yaml:"club_id"`
	} `yaml:"spirit"`
	Database struct {
		Address  string `yaml:"address"`
		Port     int    `yaml:"port"`
		Database string `yaml:"database"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"database"`
	CronWithSeconds string `yaml:"cron_with_seconds"`
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("can't open config file: %w", err)
	}
	defer file.Close()

	var root yaml.Node
	d := yaml.NewDecoder(file)
	if err := d.Decode(&root); err != nil {
		return nil, fmt.Errorf("can't parse config file: %w", err)
	}

	replaceEnvVars(&root)

	config := &Config{}
	if err := root.Decode(config); err != nil {
		return nil, fmt.Errorf("can't map config to struct: %w", err)
	}

	fmt.Printf("%+v\n", config)
	return config, nil
}

// replaceEnvVars обходит ноды и подставляет значения из окружения
func replaceEnvVars(node *yaml.Node) {
	if node.Kind == yaml.ScalarNode && node.Tag == "!env_str" {
		node.Value = os.Getenv(node.Value)
		node.Tag = "" // убираем тег после подстановки
	}

	for _, child := range node.Content {
		replaceEnvVars(child)
	}
}
