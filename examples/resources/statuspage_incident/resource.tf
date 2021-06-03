resource "statuspage_component" "my_component" {
  page_id     = "my_page_id"

  name        = "My Website"
  description = "Status of my Website"
  status      = "operational"
}

resource "statuspage_component" "my_component_2" {
  page_id     = "my_page_id"

  name        = "My Website 2"
  description = "Status of my Website 2"
  status      = "operational"
}

resource "statuspage_incident" "my_incident" {
  page_id     = "my_page_id"

  name    = "Incident name"
  impact_override = "none"
  status = "investigating"
  body   = "We are currently investigating the issue."

  component {
    id = statuspage_component.my_component.id
    name = statuspage_component.my_component.name
    status = "under_maintenance"
  }

  component {
    id = statuspage_component.my_component_2.id
    name = statuspage_component.my_component_2.name
    status = "under_maintenance"
  }

}