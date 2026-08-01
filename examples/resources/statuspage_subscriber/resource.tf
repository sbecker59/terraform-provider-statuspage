# Email subscriber
resource "statuspage_subscriber" "email" {
  page_id = "my_page_id"
  email   = "subscriber@example.com"
}

# Webhook subscriber
resource "statuspage_subscriber" "webhook" {
  page_id  = "my_page_id"
  endpoint = "https://my-app.example.com/hooks/statuspage"
}