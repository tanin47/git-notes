package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Repos []string `json:"repos"`
}

type ConfigReader interface {
	Read(path string) (*Config, error)
}

type JsonConfigReader struct{}

func (c *JsonConfigReader) Read(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)

	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
