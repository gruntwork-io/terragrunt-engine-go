# Example Engine Client - Server implementation

To build locally the client and server, you can use the following commands in the repository root:
```bash

make examples
```

## Example HCL Configuration

Here is an example of how you can configure the IaC engine client in your Terragrunt configuration:

* update `docker-compose.yml` with path to the project where IaC tool should be invoked
* run `docker-compose up` to start the server
* prepare the client configuration in `terragrunt.hcl` file
```hcl
# terragrunt.hcl
engine {
  source = "https://github.com/gruntwork-io/terragrunt-engine-go/releases/download/v0.0.3-rc2024081902/terragrunt-iac-engine-client_rpc_v0.0.3_linux_amd64.zip"
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
