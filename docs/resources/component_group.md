---
page_title: "statuspage_component_group Resource - terraform-provider-statuspage"
subcategory: ""
description: |-
  
---

# Resource `statuspage_component_group`



## Example Usage

```terraform
resource "statuspage_component_group" "my_group" {
    page_id     = "my_page_id"
    name        = "Terraform"
    description = "Created by terraform"
    components  = ["${statuspage_component.my_component.id}"]
}
```

## Schema

### Required

- **components** (Set of String, Required) An array with the IDs of the components in this group
- **name** (String, Required) An array with the IDs of the components in this group
- **page_id** (String, Required) the ID of the page this component group belongs to

### Optional

- **description** (String, Optional) More detailed description for this component group
- **id** (String, Optional) The ID of this resource.

