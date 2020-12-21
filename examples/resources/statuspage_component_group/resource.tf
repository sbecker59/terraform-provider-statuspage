resource "statuspage_component_group" "my_group" {
  page_id     = "my_page_id"
  name        = "Terraform"
  description = "Created by terraform"
  components  = ["${statuspage_component.my_component.id}"]
}