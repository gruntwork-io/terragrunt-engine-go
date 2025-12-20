default: test

lint: SHELL:=/bin/bash
lint:
	golangci-lint run -c <(curl -s https://raw.githubusercontent.com/gruntwork-io/terragrunt/main/.golangci.yml) ./...

update-local-lint: SHELL:=/bin/bash
update-local-lint:
	curl -s https://raw.githubusercontent.com/gruntwork-io/terragrunt/main/.golangci.yml --output .golangci.yml
	tmpfile=$$(mktemp) ;\
	echo '# This file is generated using `make update-local-lint` to track the linting used in Terragrunt. Do not edit manually.' | cat - .golangci.yml > $${tmpfile} && mv $${tmpfile} .golangci.yml


test:
	go test ./... -v

tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.5
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6

protoc: $(shell find ./proto -name '*.proto')
	protoc --go_out=. --go_opt=paths=source_relative proto/engine.proto
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/engine.proto

pre-commit:
	pre-commit run --all-files

.PHONY: default lint protoc test tools
