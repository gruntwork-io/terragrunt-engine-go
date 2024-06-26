# terragrunt-engine-go

This repository contains the implementation of the Terragrunt Engine written in Go.
It uses gRPC for communication and Protocol Buffers for data serialization, ensuring high performance and scalability.

## Engine Methods

* `Init(InitRequest) returns (stream InitResponse)`: Initializes the engine.
* `Run(RunRequest) returns (stream RunResponse)`: Runs a command.
* `Shutdown(ShutdownRequest) returns (stream ShutdownResponse)`: Shuts down the engine.

