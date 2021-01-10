
# Terraform Provider StatusPage

![release](https://github.com/sbecker59/terraform-provider-statuspage/workflows/release/badge.svg)
[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/sbecker59/terraform-provider-statuspage)

Run the following command to build the provider

```shell
go build -o terraform-provider-statuspage
```

## Test sample configuration

First, build and install the provider.

```shell
make install
```

Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
terraform init && terraform apply
```
