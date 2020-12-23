# Terraform Provider StatusPage

![release](https://github.com/sbecker59/terraform-provider-statuspage/workflows/release/badge.svg)

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

export TF_LOG=DEBUG
export SP_API_KEY=6c1efb96-4e5e-4c27-92a8-d4a69f69a7d3
export STATUSPAGE_PAGE_ID=9jt380nsjqh1
make testacc