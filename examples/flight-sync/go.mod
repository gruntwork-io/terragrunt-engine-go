module github.com/gruntwork-io/terragrunt-engine-go/examples/flight-sync

go 1.27

require (
	github.com/gruntwork-io/terragrunt-engine-go v0.1.0
	github.com/hashicorp/go-hclog v1.6.3
	github.com/hashicorp/go-plugin v1.8.0
	github.com/stretchr/testify v1.12.1
	golang.org/x/sync v0.24.0
	golang.org/x/sys v0.48.0
)

require (
	github.com/fatih/color v1.19.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/hashicorp/yamux v0.1.2 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/oklog/run v1.2.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/net v0.59.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260921155816-b14227669459 // indirect
	google.golang.org/grpc v1.84.0 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace github.com/gruntwork-io/terragrunt-engine-go => ../..
