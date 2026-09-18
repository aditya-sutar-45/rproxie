run-proxy:
	go run cmd/rproxie/main.go

build-proxy:
	go build -o ./bin/rproxie ./cmd/rproxie

start:
	./bin/rproxie/rproxie
