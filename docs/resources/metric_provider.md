---
page_title: "statuspage_metric_provider Resource - terraform-provider-statuspage"
subcategory: ""
description: |-
  
---

# Resource `statuspage_metric_provider`





## Schema

### Required

- **page_id** (String, Required) the ID of the page this component group belongs to
- **type** (String, Required) One of 'Pingdom', 'NewRelic', 'Librato', 'Datadog', or 'Self'

### Optional

- **api_key** (String, Optional) Required by the Datadog and NewRelic type metrics providers
- **api_token** (String, Optional) Required by the Librato, Pingdom and Datadog type metrics providers
- **application_key** (String, Optional) Required by the Pingdom-type metrics provider
- **email** (String, Optional) Required by the Librato and Pingdom type metrics providers
- **id** (String, Optional) The ID of this resource.
- **metric_base_uri** (String, Optional) Required by the NewRelic-type metrics provider
- **password** (String, Optional) Required by the Pingdom-type metrics provider


