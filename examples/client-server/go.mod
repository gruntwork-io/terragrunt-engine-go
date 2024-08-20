module github.com/gruntwork-io/terragrunt-engine-go/example/client-server

go 1.21

require (
	github.com/gruntwork-io/terragrunt-engine-go v0.0.2
	github.com/hashicorp/go-plugin v1.6.1
	github.com/sirupsen/logrus v1.9.3
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.2
)

require (
	github.com/fatih/color v1.7.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/hashicorp/go-hclog v0.14.1 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/mitchellh/go-testing-interface v0.0.0-20171004221916-a61a99592b77 // indirect
	github.com/oklog/run v1.0.0 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240711142825-46eb208f015d // indirect
)

replace github.com/gruntwork-io/terragrunt-engine-go => ../..
