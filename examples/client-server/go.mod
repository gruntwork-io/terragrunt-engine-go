module github.com/gruntwork-io/terragrunt-engine-go/examples/client-server

go 1.23.0

toolchain go1.23.1

require (
	github.com/gruntwork-io/terragrunt-engine-go v0.0.8
	github.com/hashicorp/go-plugin v1.6.3
	github.com/sirupsen/logrus v1.9.3
	google.golang.org/grpc v1.71.0
	google.golang.org/protobuf v1.36.5
)

require (
	github.com/fatih/color v1.18.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/yamux v0.1.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/oklog/run v1.1.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250303144028-a0af3efb3deb // indirect
)

replace github.com/gruntwork-io/terragrunt-engine-go => ../..
