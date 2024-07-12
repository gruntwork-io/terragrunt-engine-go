# terragrunt-engine-go

This repository contains the implementation of the Terragrunt Engine written in Go.
It uses gRPC for communication and Protocol Buffers for data serialization, ensuring high performance and scalability.

Make commands:
- `make tools`: Install tools required for development.
- `make proto`: Generate Go code from Protocol Buffers definitions.
- `make lint`: Run linters.
- `make test`: Run tests.

## Engine Methods

* `Init(InitRequest) returns (stream InitResponse)`: Initializes the engine.
* `Run(RunRequest) returns (stream RunResponse)`: Runs a command.
* `Shutdown(ShutdownRequest) returns (stream ShutdownResponse)`: Shuts down the engine.

## Example engine implementation

Example engine implementation can be found in the `test/engine_test.go` file.

## License

[Mozilla Public License v2.0](./LICENSE)

