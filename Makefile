default: test

lint:
	golangci-lint run ./...

test:
	go test ./...

tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2

protoc: $(shell find ./proto -name '*.proto')
	protoc --go_out=. --go_opt=paths=source_relative proto/engine.proto
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/engine.proto

protoc-examples: $(shell find ./example/client-server -name '*.proto')
	protoc --go_out=. --go_opt=paths=source_relative example/client-server/proto/proto.proto
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative example/client-server/proto/proto.proto

examples: protoc-examples
	set -xe ;\
	cd example/client-server ;\
	cd client ;\
	go build -o terragrunt-engine-client -ldflags "-extldflags '-static'" . ;\
	cd .. ;\
	cd server ;\
	go build -o terragrunt-engine-server -ldflags "-extldflags '-static'" .

pre-commit:
	pre-commit run --all-files

.PHONY: default lint protoc test tools
