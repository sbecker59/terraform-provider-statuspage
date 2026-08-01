
# Terraform Provider for Atlassian Statuspage

[![Terraform Registry](https://img.shields.io/badge/registry-sbecker59%2Fstatuspage-623CE4?logo=terraform)](https://registry.terraform.io/providers/sbecker59/statuspage/latest/docs)
![release](https://github.com/sbecker59/terraform-provider-statuspage/workflows/release/badge.svg)
[![codecov](https://codecov.io/gh/sbecker59/terraform-provider-statuspage/branch/main/graph/badge.svg?token=OalDkaUlvu)](https://codecov.io/gh/sbecker59/terraform-provider-statuspage)
[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/sbecker59/terraform-provider-statuspage)

Manage your [Atlassian Statuspage](https://www.atlassian.com/software/statuspage) resources with Terraform.  
This provider lets you automate components, incidents, metrics, subscribers, and access control for any Statuspage page — fully declarative, version-controlled, and compatible with modern Terraform workflows.

📖 **Full documentation:** [registry.terraform.io/providers/sbecker59/statuspage](https://registry.terraform.io/providers/sbecker59/statuspage/latest/docs)

---

## Features

- **Terraform Plugin SDK v2** — compatible with Terraform ≥ 1.0
- **Actively maintained** with an up-to-date Go Statuspage API client
- Covers resources that competing providers do not: `metric_provider`, `page_access_group`, `page_access_user`
- Automatic retry logic for transient API errors (via `go-retryablehttp`)
- Import support for all major resources (`terraform import`)

---

## Available Resources

| Resource | Description |
|---|---|
| `statuspage_component` | Create and manage status components (API, website, database, …) |
| `statuspage_component_group` | Group components into logical sections on your status page |
| `statuspage_incident` | Declare realtime incidents and scheduled maintenance windows |
| `statuspage_metric_provider` | Connect a metrics source (Datadog, NewRelic, Librato, Pingdom, Self) |
| `statuspage_subscriber` | Add email or webhook subscribers to your status page |
| `statuspage_page_access_group` | Manage access groups on audience-restricted pages |
| `statuspage_page_access_user` | Grant individual users access to a restricted page |

## Available Data Sources

| Data Source | Description |
|---|---|
| `statuspage_pages` | Look up a Statuspage page by name |
| `statuspage_components` | List and filter components on a page |
| `statuspage_component_groups` | List and filter component groups on a page |

---

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) ≥ 1.0
- [Go](https://golang.org/doc/install) ≥ 1.26 (only to build from source)
- A [Statuspage API key](https://developer.statuspage.io/#section/Authentication)

---

## Quick Start

### 1 — Configure the provider

```hcl
terraform {
  required_providers {
    statuspage = {
      source  = "sbecker59/statuspage"
      version = "~> 1.0"
    }
  }
}

provider "statuspage" {
  # Can also be set via STATUSPAGE_API_KEY or SP_API_KEY environment variables
  api_key = var.statuspage_api_key
}
```

### 2 — Create a component

```hcl
resource "statuspage_component" "api" {
  page_id     = "your_page_id"
  name        = "API"
  description = "Core API availability"
  status      = "operational"
  showcase    = true

  lifecycle {
    ignore_changes = [status] # status is managed via incidents
  }
}
```

### 3 — Declare an incident

```hcl
resource "statuspage_incident" "outage" {
  page_id         = "your_page_id"
  name            = "API degraded performance"
  status          = "investigating"
  impact_override = "major"
  body            = "We are investigating reports of degraded API performance."

  component {
    id     = statuspage_component.api.id
    name   = statuspage_component.api.name
    status = "degraded_performance"
  }
}
```

### 4 — Query existing data

```hcl
data "statuspage_pages" "my_page" {
  page_name = "My Status Page"
}

data "statuspage_components" "all" {
  page_id = data.statuspage_pages.my_page.id
}
```

---

## Installation

The provider is published on the [Terraform Registry](https://registry.terraform.io/providers/sbecker59/statuspage/latest).  
Add the `required_providers` block above to any Terraform configuration and run:

```shell
terraform init
```

Terraform will download the provider automatically — no manual installation needed.

---

## Authentication

The Statuspage API key can be supplied in three ways (in order of precedence):

1. `api_key` argument in the `provider` block
2. `STATUSPAGE_API_KEY` environment variable
3. `SP_API_KEY` environment variable

Generate an API key in your Statuspage account under **Profile → API info**.

---

## Import

Most resources support `terraform import`. The import ID format is `<page_id>/<resource_id>`:

```shell
terraform import statuspage_component.api your_page_id/your_component_id
terraform import statuspage_incident.outage your_page_id/your_incident_id
terraform import statuspage_page_access_user.alice your_page_id/alice@example.com
```

---

## Development

### Build

```shell
go build -o terraform-provider-statuspage
```

### Install locally

```shell
make install
```

### Run acceptance tests

```shell
export STATUSPAGE_API_KEY=your_api_key
export STATUSPAGE_PAGE_ID=your_page_id
make testacc
```

### Generate documentation

```shell
make docs
```

---

## Contributing

Contributions are welcome! Please open an issue or a pull request on [GitHub](https://github.com/sbecker59/terraform-provider-statuspage).

---

## License

[Mozilla Public License 2.0](LICENSE)
