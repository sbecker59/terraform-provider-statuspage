---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "statuspage_subscriber Resource - terraform-provider-statuspage"
subcategory: ""
description: |-
  
---

# statuspage_subscriber (Resource)



## Example Usage

```terraform
resource "statuspage_subscriber" "my_subscriber" {
  page_id     = "my_page_id"
  email        = "my_email@example.com"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `page_id` (String) the ID of the page this component belongs to

### Optional

- `email` (String) the email address for creating Email and Webhook subscribers
- `endpoint` (String) The endpoint URI for creating Webhook subscribers

### Read-Only

- `id` (String) The ID of this resource.
