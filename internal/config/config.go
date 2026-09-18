// Package config
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port     int
	Backends []string
}

func Load() (*Config, error) {
	backends := strings.Split(os.Getenv("BACKENDS"), ",")

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return nil, fmt.Errorf("could not convert string to int: %v", err)
	}

	return &Config{
		Port:     port,
		Backends: backends,
	}, nil
}
