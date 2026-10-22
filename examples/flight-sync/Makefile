default: build

build:
	go build -o terragrunt-engine-flight-sync ./engine
	go build -o flight-sync-auth ./auth

test:
	go test ./...

lint:
	golangci-lint run -c ../../.golangci.yml ./...

.PHONY: default build test lint
