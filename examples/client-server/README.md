# Example Engine Client - Server implementation

Example implementation of the Terragrunt IaC engine client and server.
Use it only for testing purposes since it is allowing execution of arbitrary bash commands on the server.

To build the client and server locally, run the `make` command:
```bash
make
```
This will build the `terragrunt-engine-client` and `terragrunt-engine-server` binaries.

## Example HCL Configuration

Here is an example of how you can configure the IaC engine client in your Terragrunt configuration for AMD64 Linux:
* run `docker compose up` to start the server
* prepare the client configuration in `terragrunt.hcl` file
```hcl
# terragrunt.hcl
engine {
  source = "./terragrunt-iac-engine-client"
  meta = {
    # server endpoint
    endpoint = "127.0.0.1:50051"
    # authentication token
    token    = get_env("TG_SERVER_TOKEN")
  }
}
```

Terragrunt execution:
```bash
export TG_EXPERIMENTAL_ENGINE=1
export TG_SERVER_TOKEN=secret-token

terragrunt apply --auto-approve
```

End to end example:
[![asciicast](https://asciinema.org/a/672387.svg)](https://asciinema.org/a/672387)