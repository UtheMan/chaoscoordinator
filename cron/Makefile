.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -o bin/chaoscoordinatorservice -v cron.go
modules:
	go mod download