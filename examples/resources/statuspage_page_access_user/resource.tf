resource "statuspage_page_access_user" "alice" {
  page_id = "my_page_id"
  email   = "alice@example.com"
}

resource "statuspage_page_access_user" "bob" {
  page_id = "my_page_id"
  email   = "bob@example.com"
}
