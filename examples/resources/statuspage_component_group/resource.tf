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

resource "statuspage_component_group" "my_group" {
  page_id     = "my_page_id"
  name        = "My Product"
  description = "All components belonging to My Product"
  components  = [
    statuspage_component.api.id,
    statuspage_component.dashboard.id,
  ]
}