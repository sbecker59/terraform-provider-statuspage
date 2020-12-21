resource "statuspage_component" "my_component" {
  page_id     = "my_page_id"
  name        = "My Website"
  description = "Status of my Website"
  status      = "operational"
}