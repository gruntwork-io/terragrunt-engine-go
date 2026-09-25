default: test

lint:
	golangci-lint run ./...

update-local-lint: SHELL:=/bin/bash
update-local-lint:
	curl -s https://raw.githubusercontent.com/gruntwork-io/terragrunt/main/.golangci.yml --output .golangci.yml
	tmpfile=$$(mktemp) ;\
	echo '# This file is generated using `make update-local-lint` to track the linting used in Terragrunt. Do not edit manually.' | cat - .golangci.yml > $${tmpfile} && mv $${tmpfile} .golangci.yml


test:
	go test ./... -v

protoc: $(shell find ./proto -name '*.proto')
	protoc --go_out=. --go_opt=paths=source_relative proto/engine.proto
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/engine.proto

pre-commit:
	pre-commit run --all-files

.PHONY: default lint protoc test
