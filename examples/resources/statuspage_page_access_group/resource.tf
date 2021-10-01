resource "statuspage_component" "my_component" {
  page_id     = "my_page_id"
  name        = "My Website"
  description = "Status of my Website"
  status      = "operational"

  lifecycle {
      ignore_changes = [
          status
      ]
  }
}

resource "statuspage_page_access_user" "my_user" {
  page_id = "my_page_id"
  email   = "my_user@example.com"
}

resource "statuspage_page_access_group" "my_user_group" {
  page_id    = "my_page_id"
  name       = "My Page Access User Group"
  users      = [ statuspage_page_access_user.my_user.id ]
  components = [ statuspage_component.my_component.id ]
}
