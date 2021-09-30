resource "statuspage_page_access_group" "my_user_group" {
  page_id    = "my_page_id"
  name       = "My Page Access User Group"
  users      = ["${statuspage_page_access_user.my_user.id}"]
  components = ["${statuspage_component.my_component.id}"]
}
