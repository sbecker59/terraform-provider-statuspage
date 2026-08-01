data "statuspage_components" "web" {
  page_id = "my_page_id"

  filter {
    name   = "name"
    values = ["API", "Dashboard"]
  }
}