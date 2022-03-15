package statuspage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccStatuspagePagesDatasource(t *testing.T) {

	

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceStatuspagePagesConfig(),
				// Because of the `depends_on` in the datasource, the plan cannot be empty.
				// See https://www.terraform.io/docs/configuration/data-sources.html#data-resource-dependencies
				Check: checkDatasourceStatuspagePagesAttrs(testAccProvider),
			},
		},
	})
}

func checkDatasourceStatuspagePagesAttrs(accProvider *schema.Provider) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("data.statuspage_pages.default", "page_name", pageName),
	)
}

func testAccStatuspagePagesConfig() string {
	return fmt.Sprintf(`
	variable "pageName" {
		default = "%s"
	}
	`, pageName)
}

func testAccDatasourceStatuspagePagesConfig() string {
	return fmt.Sprintf(`
	%s
	data "statuspage_pages" "default" {
		page_name = "${var.pageName}"
	}`, testAccStatuspagePagesConfig())
}
