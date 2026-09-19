package loadbalancer

import (
	"sync"
	"testing"
)

func TestRoundRobinNextBackendIndex(t *testing.T) {
	lb := New(3)

	expected := []int{0, 1, 2, 0, 1, 2}

	for _, want := range expected {
		got, err := lb.NextBackendIndex()
		if err != nil {
			t.Fatal(err)
		}

		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	}
}

func TestRoundRobinWithOneBackend(t *testing.T) {
	lb := New(1)

	expected := []int{0, 0, 0, 0, 0, 0}
	for _, want := range expected {
		got, err := lb.NextBackendIndex()
		if err != nil {
			t.Fatal(err)
		}

		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	}
}

func TestRoundRobinWithZeroBackends(t *testing.T) {
	lb := New(0)

	expected := -1

	got, err := lb.NextBackendIndex()
	if err == nil {
		t.Fatal("no error thrown when 0 backends are passed")
	}

	if got != expected {
		t.Errorf("got %d want %d", got, expected)
	}
}

func TestRoundRobinMutex(t *testing.T) {
	lb := New(3)

	var wg sync.WaitGroup
	results := make(chan int, 300)

	for range 300 {
		wg.Go(func() {
			index, err := lb.NextBackendIndex()
			if err != nil {
				t.Error(err)
				return
			}

			results <- index
		})
	}

	wg.Wait()
	close(results)

	counts := make([]int, 3)

	for index := range results {
		counts[index]++
	}

	for index, count := range counts {
		if count != 100 {
			t.Errorf(
				"backend %d selected %d times, want 100",
				index,
				count,
			)
		}
	}
}
