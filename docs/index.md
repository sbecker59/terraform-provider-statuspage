---
page_title: "statuspage Provider"
subcategory: ""
description: |-
  
---

# statuspage Provider



## Example Usage

```terraform
terraform {
    required_providers {
        statuspage = {
            version = "0.1"
            source = "hashicorp.com/sbecker59/statuspage"
        }
    }
}

provider "statuspage" {}
```

## Schema

### Optional

- **api_key** (String, Optional)
