module github.com/gruntwork-io/terragrunt-engine-go/examples/client-server

go 1.24

require (
	github.com/gruntwork-io/terragrunt-engine-go v0.0.11
	github.com/hashicorp/go-plugin v1.6.3
	github.com/sirupsen/logrus v1.9.3
	google.golang.org/grpc v1.73.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/fatih/color v1.18.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/yamux v0.1.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/oklog/run v1.1.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822 // indirect
)

replace github.com/gruntwork-io/terragrunt-engine-go => ../..
