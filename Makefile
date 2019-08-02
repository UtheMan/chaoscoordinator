.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -o bin/chaos -v cmd/chaos.go
modules:
	go mod download