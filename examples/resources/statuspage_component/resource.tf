resource "statuspage_component" "my_component" {
  page_id     = "my_page_id"
  name        = "My Website"
  description = "Tracks the availability of the main website"
  status      = "operational"

  # Show this component even when it is not degraded
  showcase              = true
  only_show_if_degraded = false

  # Optionally pin an SLA start date (YYYY-MM-DD)
  # start_date = "2024-01-01"

  lifecycle {
    # Let Statuspage control the live status; only manage it via incidents.
    ignore_changes = [status]
  }
}