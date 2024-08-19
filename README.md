# terragrunt-engine-go

This repository contains the implementation of the [Terragrunt](https://github.com/gruntwork-io/terragrunt) Engine written in Go.
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

## Engines

Implementations of the Terragrunt Engine include:

- [terragrunt-engine-opentofu](https://github.com/gruntwork-io/terragrunt-engine-opentofu)
- [terragrunt-engine-client-server](./engine/client-server)

## Example Engine Implementation

Example engine implementation can be found in the [test/engine_test.go](./test/engine_test.go) file.

## Contributions

Contributions are welcome! Check out the [Contributing Guidelines](./CONTRIBUTING.md) for more information.

## License

[Mozilla Public License v2.0](./LICENSE)

