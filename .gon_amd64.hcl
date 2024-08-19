# See https://github.com/gruntwork-io/terraform-aws-ci/blob/main/modules/sign-binary-helpers/
# for further instructions on how to sign the binary + submitting for notarization.

source = ["./example/client-server/client/bin/terragrunt-iac-engine-client_darwin_amd64", "example/client-server/server/bin/terragrunt-iac-engine-server_darwin_amd64"]

bundle_id = "io.gruntwork.app.terragrunt"

apple_id {
  username = "machine.apple@gruntwork.io"
}

sign {
  application_identity = "Developer ID Application: Gruntwork, Inc."
}

zip {
  output_path = "terragrunt-iac-engine-client-server_amd64.zip"
}
