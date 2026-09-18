package loadbalancer

import (
	"fmt"
	"sync"
)

type RoundRobin struct {
	backendCount        int
	currentBackendIndex int
	mu                  sync.Mutex
}

func New(backendCount int) *RoundRobin {
	return &RoundRobin{
		backendCount:        backendCount,
		currentBackendIndex: 0,
	}
}

func (lb *RoundRobin) NextBackendIndex() (int, error) {
	if lb.backendCount <= 0 {
		return -1, fmt.Errorf("no backends configured, backendCount: %d", lb.backendCount)
	}

	lb.mu.Lock()
	defer lb.mu.Unlock()

	currentIndex := lb.currentBackendIndex

	nextBackendIndex := (lb.currentBackendIndex + 1) % lb.backendCount
	lb.currentBackendIndex = nextBackendIndex

	return currentIndex, nil
}
