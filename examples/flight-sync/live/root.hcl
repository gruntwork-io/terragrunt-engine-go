engine {
  source = "${get_parent_terragrunt_dir()}/../terragrunt-engine-flight-sync"
}

terraform {
  source = "${get_parent_terragrunt_dir()}/../module"
}
