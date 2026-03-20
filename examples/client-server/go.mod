module github.com/gruntwork-io/terragrunt-engine-go/examples/client-server

go 1.26

require (
	github.com/gruntwork-io/terragrunt-engine-go v0.1.0
	github.com/hashicorp/go-plugin v1.7.0
	github.com/sirupsen/logrus v1.9.4
	google.golang.org/grpc v1.79.3
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/fatih/color v1.19.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/yamux v0.1.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/oklog/run v1.2.0 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260319201613-d00831a3d3e7 // indirect
)

replace github.com/gruntwork-io/terragrunt-engine-go => ../..
