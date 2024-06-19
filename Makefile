default: test

lint:
	golangci-lint run ./...

test:
	go test ./...

tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2

protoc: $(shell find . -name '*.proto')
	protoc --go_out=. --go_opt=paths=source_relative engine/engine.proto
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative engine/engine.proto

.PHONY: default lint protoc test tools
