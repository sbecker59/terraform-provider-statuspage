terraform {
    required_providers {
        statuspage = {
            version = "0.1"
            source = "hashicorp.com/sbecker59/statuspage"
        }
    }
}

provider "statuspage" {}