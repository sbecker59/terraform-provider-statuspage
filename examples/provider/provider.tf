terraform {
  required_providers {
    statuspage = {
      version = "0.1.0"
      source  = "sbecker59/statuspage"
    }
  }
}

provider "statuspage" {}
