// Package config
package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

type RProxieConfig struct {
	Server      ServerConfig       `yaml:"server"`
	Backends    []*BackendConfig   `yaml:"backends"`
	RateLimiter *RateLimiterConfig `yaml:"rate_limiter"`
	Logging     *LoggingConfig     `yaml:"logging"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type BackendConfig struct {
	ID  string `yaml:"id"`
	URL string `yaml:"url"`
}

type RateLimiterConfig struct {
	Enabled             bool `yaml:"enabled"`
	BucketCapacity      int  `yaml:"bucket_capacity"`
	RefillRatePerSecond int  `yaml:"refill_rate_per_second"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}

func LoadYAML(path string) (*RProxieConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %v", err)
	}

	var config RProxieConfig

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return &config, nil
}
