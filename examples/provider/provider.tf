terraform {
  required_providers {
    statuspage = {
      source  = "sbecker59/statuspage"
      version = "~> 1.0"
    }
  }
}

# The API key can also be provided via the STATUSPAGE_API_KEY or SP_API_KEY environment variables.
provider "statuspage" {
  api_key = var.statuspage_api_key
}

variable "statuspage_api_key" {
  description = "Atlassian Statuspage API key"
  type        = string
  sensitive   = true
}
