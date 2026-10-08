package ratelimiter

import (
	"math"
	"sync"
	"testing"
	"time"
)

func TestStartsWithFullBucket(t *testing.T) {
	r := New(5, 0)

	expected := []bool{true, true, true, true, true, false}
	cnt := len(expected)

	results := make([]bool, 0, cnt)

	for range cnt {
		r := r.Allow()
		results = append(results, r)
	}

	for i := range len(results) {
		if results[i] != expected[i] {
			t.Errorf("index %d: expected=%v, got=%v", i, expected[i], results[i])
		}
	}
}

func TestRequestConsumesOneToken(t *testing.T) {
	r := New(3, 0)

	for i := range 3 {
		if !r.Allow() {
			t.Errorf("request %d should have been allowed", i+1)
		}
	}

	if r.Allow() {
		t.Error("request should have been rejected after all tokens were consumed")
	}
}

func TestTokensRefillOverTime(t *testing.T) {
	r := New(1, 10)

	if !r.Allow() {
		t.Error("request should have been allowed")
	}

	if r.Allow() {
		t.Error("request should not have been allowed")
	}

	time.Sleep(150 * time.Millisecond)

	if !r.Allow() {
		t.Error("request should have been allowed")
	}
}

func TestCalculateTokenCount(t *testing.T) {
	expected := 3.0
	tolerance := 0.1

	r := &RateLimiter{
		tokenCount:          2,
		lastRefillTime:      time.Now().Add(-100 * time.Millisecond),
		refillRatePerSecond: 10,
	}

	tokenCount := r.calculateTokenCount()

	if math.Abs(tokenCount-expected) > tolerance {
		t.Errorf("expected≈%v, got=%v", expected, tokenCount)
	}
}

func TestUpdateTokenCount(t *testing.T) {
	r := &RateLimiter{
		tokenCount:          0,
		bucketCapacity:      50,
		lastRefillTime:      time.Now(),
		refillRatePerSecond: 0,
	}

	newTokenCount := 100.69
	r.updateTokenCount(newTokenCount)

	if r.GetCurrentTokenCount() != 50 {
		t.Errorf("expected token count = %d, got = %f", 50, r.GetCurrentTokenCount())
	}
}

func TestConcurrentRequests(t *testing.T) {
	r := New(50, 0)
	successCount := 0
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	for range 50 {
		wg.Go(func() {
			if r.Allow() {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		})
	}

	wg.Wait()

	if successCount != 50 {
		t.Errorf("expected 50 successful requests, got %d", successCount)
	}
}
