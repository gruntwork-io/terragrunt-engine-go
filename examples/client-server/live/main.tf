# Example Terraform/OpenTofu configuration
terraform {
  required_version = ">= 1.0"
}

# Simple example resource
resource "null_resource" "example" {
  triggers = {
    timestamp = timestamp()
  }

  provisioner "local-exec" {
    command = "echo 'Hello from Terragrunt Engine!'"
  }
}

output "message" {
  value = "Example Terragrunt Engine integration"
}


