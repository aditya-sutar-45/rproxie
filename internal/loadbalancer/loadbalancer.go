// Package loadbalancer
package loadbalancer

type LoadBalancer interface {
	NextBackendIndex() (int, error)
}
