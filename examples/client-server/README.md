# Example Engine Client - Server implementation

Example implementation of the Terragrunt IaC engine to offload the running of OpenTofu commands on a remote server.
Use it only for educational purposes since it allows execution of arbitrary bash commands on the server.

To build the client and server locally, run the `make` command:

```bash
make
```

This will build the `terragrunt-engine-client` and `terragrunt-engine-server` binaries.

## Example Usage

This example includes a `live` directory with sample Terragrunt and Terraform configurations.

### Starting the Server

1. Build and start the server using Docker Compose:

```bash
docker compose up
```

### Running Terragrunt

1. Navigate to the `live` directory:

```bash
cd live
```

2. Set the required environment variables:

```bash
export TG_EXPERIMENTAL_ENGINE=1
export TG_SERVER_TOKEN=secret-token
```

3. Run Terragrunt commands:

```bash
terragrunt apply --auto-approve
```

### Example Configuration

The `live/terragrunt.hcl` file contains the engine configuration:

```hcl
# terragrunt.hcl
engine {
  source = "../terragrunt-engine-client"
  meta = {
    # server endpoint
    endpoint = "127.0.0.1:50051"
    # authentication token
    token    = get_env("TG_SERVER_TOKEN", "secret-token")
  }
}
```

End to end example:
[![asciicast](https://asciinema.org/a/672387.svg)](https://asciinema.org/a/672387)
