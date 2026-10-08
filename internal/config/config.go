// Package config
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port                         int
	Backends                     []string
	RateLimitBucketCapacity      int
	RateLimitTokenRefilPerSecond int
}

func Load() (*Config, error) {
	backends := strings.Split(os.Getenv("BACKENDS"), ",")

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return nil, fmt.Errorf("could not convert string to int: %v", err)
	}

	bucketCapacity, err := strconv.Atoi(os.Getenv("BUCKET_CAPACITY"))
	if err != nil {
		return nil, fmt.Errorf("could not convert string to int: %v", err)
	}

	refilPerSecond, err := strconv.Atoi(os.Getenv("TOKEN_REFIL_PER_SECOND"))
	if err != nil {
		return nil, fmt.Errorf("could not convert string to int: %v", err)
	}

	return &Config{
		Port:                         port,
		Backends:                     backends,
		RateLimitBucketCapacity:      bucketCapacity,
		RateLimitTokenRefilPerSecond: refilPerSecond,
	}, nil
}
