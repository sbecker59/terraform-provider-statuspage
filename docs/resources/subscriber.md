---
page_title: "statuspage_subscriber Resource - terraform-provider-statuspage"
subcategory: ""
description: |-
  
---

# Resource `statuspage_subscriber`



## Example Usage

```terraform
resource "statuspage_subscriber" "my_subscriber" {
  page_id     = "my_page_id"
  email        = "my_email@example.com"
}
```

## Schema

### Required

- **page_id** (String, Required) the ID of the page this component belongs to

### Optional

- **email** (String, Optional) the email address for creating Email and Webhook subscribers
- **endpoint** (String, Optional) The endpoint URI for creating Webhook subscribers
- **id** (String, Optional) The ID of this resource.


