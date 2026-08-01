# Datadog metrics provider
resource "statuspage_metric_provider" "datadog" {
  page_id = "my_page_id"
  type    = "Datadog"
  api_key = var.datadog_api_key
}

# Self-hosted / custom metrics provider
resource "statuspage_metric_provider" "self_hosted" {
  page_id         = "my_page_id"
  type            = "Self"
  metric_base_uri = "https://metrics.example.com"
}

variable "datadog_api_key" {
  description = "Datadog API key used to pull metrics into Statuspage"
  type        = string
  sensitive   = true
}
