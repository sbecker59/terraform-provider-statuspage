---
page_title: "statuspage_component Resource - terraform-provider-statuspage"
subcategory: ""
description: |-
  
---

# Resource `statuspage_component`



## Example Usage

```terraform
resource "statuspage_component" "my_component" {
  page_id     = "my_page_id"
  name        = "My Website"
  description = "Status of my Website"
  status      = "operational"
}
```

## Schema

### Required

- **name** (String, Required) Display Name for the component
- **page_id** (String, Required) the ID of the page this component belongs to

### Optional

- **description** (String, Optional) More detailed description for the component
- **id** (String, Optional) The ID of this resource.
- **only_show_if_degraded** (Boolean, Optional)
- **showcase** (Boolean, Optional) Should this component be shown component only if in degraded state
- **start_date** (String, Optional) Should this component be showcased
- **status** (String, Optional)


