package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Spirit struct {
		Token  string `yaml:"token"`
		ClubID string `yaml:"club_id"`
	} `yaml:"spirit"`
	GSheets struct {
		SheetID   string `yaml:"sheet_id"`
		DataRange string `yaml:"data_range"`
	} `yaml:"gsheets"`
	OpenWeather struct {
		APIKey string `yaml:"apikey"`
	} `yaml:"openweather"`
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("can't open config file: %w", err)
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, fmt.Errorf("can't parse config file: %w", err)
	}

	return config, nil
}
