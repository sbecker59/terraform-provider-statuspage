resource "statuspage_component" "api" {
  page_id     = "my_page_id"
  name        = "API"
  description = "Core API availability"
  status      = "operational"
}

resource "statuspage_component" "dashboard" {
  page_id     = "my_page_id"
  name        = "Dashboard"
  description = "Web dashboard availability"
  status      = "operational"
}

# Realtime incident affecting two components
resource "statuspage_incident" "outage" {
  page_id = "my_page_id"

  name            = "API and Dashboard degraded performance"
  status          = "investigating"
  impact_override = "major"
  body            = "We are currently investigating reports of degraded performance. Our team is actively working on a fix."

  component {
    id     = statuspage_component.api.id
    name   = statuspage_component.api.name
    status = "degraded_performance"
  }

  component {
    id     = statuspage_component.dashboard.id
    name   = statuspage_component.dashboard.name
    status = "degraded_performance"
  }
}

# Scheduled maintenance window
resource "statuspage_incident" "maintenance" {
  page_id = "my_page_id"

  name            = "Scheduled database maintenance"
  status          = "scheduled"
  impact_override = "maintenance"
  body            = "We will be performing routine database maintenance during this window. Expect brief interruptions."

  scheduled_remind_prior    = true
  scheduled_auto_in_progress = true
  scheduled_auto_completed  = true

  component {
    id     = statuspage_component.api.id
    name   = statuspage_component.api.name
    status = "under_maintenance"
  }
}