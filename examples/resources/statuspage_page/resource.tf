# Pages cannot be created or deleted through the Statuspage API — only read
# and updated. This resource can only attach to a page that already exists;
# see import.sh in this directory. `terraform destroy` (or removing the
# block) only forgets the page in Terraform state; the page itself is left
# untouched.
resource "statuspage_page" "my_page" {
  page_id = "my_page_id"

  name      = "My Status Page"
  time_zone = "Berlin" # Rails ActiveSupport name, not IANA (not "Europe/Berlin")

  # Changing domain makes Statuspage start certificate provisioning for the
  # new host, which this API cannot observe — the CNAME must already point
  # at Statuspage's edge before this is set.
  domain = "status.example.com"

  hidden_from_search        = false
  allow_email_subscribers   = true
  allow_sms_subscribers     = false
  allow_webhook_subscribers = false

  css_link_color = "3498DB"
}
