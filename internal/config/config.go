// Package config
package config

import (
	"fmt"
	"net/url"
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

func Load(path string) (*RProxieConfig, error) {
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

func (c *RProxieConfig) Validate() error {
	if c.RateLimiter == nil {
		return fmt.Errorf("rate limiter config is required")
	}

	if c.Logging == nil {
		c.Logging = &LoggingConfig{
			Level: "INFO",
		}
	}

	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid port for rproxie")
	}

	if len(c.Backends) <= 0 {
		return fmt.Errorf("backends cannot be empty")
	}

	seenBackendIDs := make(map[string]bool)
	for i, b := range c.Backends {
		if b == nil {
			return fmt.Errorf("backend is null")
		}

		if b.ID == "" {
			b.ID = fmt.Sprintf("backend-%d", i+1)
		}

		if _, exists := seenBackendIDs[b.ID]; exists {
			return fmt.Errorf("duplicate backend ID's")
		}

		seenBackendIDs[b.ID] = true

		parsedURL, err := url.Parse(b.URL)
		if err != nil {
			return fmt.Errorf("invalid URL for backend %q: %w", b.ID, err)
		}

		if (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") ||
			parsedURL.Hostname() == "" {
			return fmt.Errorf("invalid URL for backend %q: must be an HTTP/HTTPS URL with a host", b.ID)
		}
	}

	if c.RateLimiter.Enabled {
		if c.RateLimiter.BucketCapacity <= 0 {
			return fmt.Errorf("bucket capacity cannot be negative or zero")
		}
		if c.RateLimiter.RefillRatePerSecond <= 0 {
			return fmt.Errorf("refill rate per second cannot be negative or zero")
		}
	}

	if c.Logging.Level == "" {
		c.Logging.Level = "INFO"
	}

	return nil
}
