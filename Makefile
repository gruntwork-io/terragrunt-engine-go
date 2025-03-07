default: test

lint:
	golangci-lint run ./...

test:
	go test ./... -v

tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.5
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

protoc: $(shell find ./proto -name '*.proto')
	protoc --go_out=. --go_opt=paths=source_relative proto/engine.proto
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/engine.proto

pre-commit:
	pre-commit run --all-files

.PHONY: default lint protoc test tools
