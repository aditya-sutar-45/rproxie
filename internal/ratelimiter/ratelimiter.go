// Package ratelimiter uses token bucket algorithm to limit the rate of requests.
package ratelimiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	bucketCapacity      int
	tokenCount          float64
	refillRatePerSecond float64 // in tokens per secound
	lastRefillTime      time.Time
	mu                  sync.RWMutex
}

func New(bucketCapacity int, refillRatePerSecond float64) *RateLimiter {
	return &RateLimiter{
		bucketCapacity:      bucketCapacity,
		tokenCount:          float64(bucketCapacity),
		refillRatePerSecond: refillRatePerSecond,
		lastRefillTime:      time.Now(),
	}
}

func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Calculate elapsed time and  tokens earned
	tokenCount := r.calculateTokenCount()

	// Add them to current tokens and cap at capacity also Update lastRefillTime
	r.updateTokenCount(tokenCount)

	// Check whether a token is available
	if r.tokenCount < 1 {
		return false
	}

	// Consume one
	r.consumeToken()

	return true
}

func (r *RateLimiter) calculateTokenCount() float64 {
	elapsedTime := time.Since(r.lastRefillTime).Seconds()
	tokenCount := r.tokenCount + (float64(elapsedTime) * r.refillRatePerSecond)

	return tokenCount
}

func (r *RateLimiter) updateTokenCount(newTokenCount float64) {
	if newTokenCount > float64(r.bucketCapacity) {
		newTokenCount = float64(r.bucketCapacity)
	}

	r.tokenCount = newTokenCount
	r.lastRefillTime = time.Now()
}

func (r *RateLimiter) consumeToken() {
	r.tokenCount -= 1
}

func (r *RateLimiter) GetCurrentTokenCount() float64 { return r.tokenCount }
