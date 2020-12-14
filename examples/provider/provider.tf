terraform {
    required_providers {
        statuspage = {
            version = "0.1"
            source = "hashicorp.com/sbecker59/statuspage"
        }
    }
}

provider "statuspage" {
    api_key = "da27995a-08ff-4c2a-b513-65b0eecf0e1a"
}