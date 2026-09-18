run-proxy:
	go run cmd/rproxie/main.go

build-proxy:
	go build -o ./bin/rproxie ./cmd/rproxie

start:
	./bin/rproxie/rproxie

test-backends:
	PORT=9001 BACKEND_ID=backend-1 go run ./cmd/testbackend &
	PORT=9002 BACKEND_ID=backend-2 go run ./cmd/testbackend &
	PORT=9003 BACKEND_ID=backend-3 go run ./cmd/testbackend &
	wait
