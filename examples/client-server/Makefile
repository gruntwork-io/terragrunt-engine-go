default: build

tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.5
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6

protoc: $(shell find ./examples/client-server -name '*.proto')
	protoc --go_out=. --go_opt=paths=source_relative proto/proto.proto
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/proto.proto

build:
	set -xe ;\
	cd client ;\
	go build -o ../terragrunt-engine-client -ldflags "-extldflags '-static'" . ;\
	cd .. ;\
	cd server ;\
	go build -o ../terragrunt-engine-server -ldflags "-extldflags '-static'" .

docker:
	cd ../.. && docker build -f examples/client-server/Dockerfile -t terragrunt-engine-server .

pre-commit:
	pre-commit run --all-files

.PHONY: default lint protoc test tools
