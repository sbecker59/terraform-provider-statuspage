resource "statuspage_component" "internal_api" {
  page_id     = "my_page_id"
  name        = "Internal API"
  description = "Private API available only to authorized users"
  status      = "operational"
  showcase    = false

  lifecycle {
    ignore_changes = [status]
  }
}

resource "statuspage_page_access_user" "alice" {
  page_id = "my_page_id"
  email   = "alice@example.com"
}

resource "statuspage_page_access_user" "bob" {
  page_id = "my_page_id"
  email   = "bob@example.com"
}

resource "statuspage_page_access_group" "internal_team" {
  page_id = "my_page_id"
  name    = "Internal Team"
  users   = [
    statuspage_page_access_user.alice.id,
    statuspage_page_access_user.bob.id,
  ]
  components = [
    statuspage_component.internal_api.id,
  ]
}
