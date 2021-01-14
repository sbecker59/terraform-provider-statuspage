---
page_title: "statuspage_components Data Source - terraform-provider-statuspage"
subcategory: ""
description: |-
  
---

# Data Source `statuspage_components`



## Example Usage

```terraform
data "statuspage_components" "default" {
    
    page_id = local.page_id

    filter {
        name = "name"
        values = [ "value_1", "value_2" ]
    }
}
```

## Schema

### Required

- **page_id** (String, Required) the ID of the page this component belongs to

### Optional

- **filter** (Block Set) (see [below for nested schema](#nestedblock--filter))
- **id** (String, Optional) The ID of this resource.

### Read-only

- **components** (List of Object, Read-only) (see [below for nested schema](#nestedatt--components))

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- **name** (String, Required)
- **values** (List of String, Required)

Optional:

- **regex** (Boolean, Optional)


<a id="nestedatt--components"></a>
### Nested Schema for `components`

- **description** (String)
- **group_id** (String)
- **id** (String)
- **name** (String)
- **position** (Number)


