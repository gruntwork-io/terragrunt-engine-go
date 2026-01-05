# Example Terragrunt configuration using the client-server engine
engine {
  source = "../terragrunt-engine-client"
  meta = {
    # server endpoint
    endpoint = "127.0.0.1:50051"
    # authentication token
    token    = get_env("TG_SERVER_TOKEN", "secret-token")
  }
}


