data "statuspage_component_groups" "all" {
  page_id = "my_page_id"

  filter {
    name   = "name"
    values = ["My Product"]
  }
}