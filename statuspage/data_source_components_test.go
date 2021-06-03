package statuspage

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccStatuspageComponentsDatasource(t *testing.T) {

	time.Sleep(10 * time.Second)
	rid := acctest.RandIntRange(1, 99)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStatuspageComponentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceStatuspageComponentsConfig(rid),
				// Because of the `depends_on` in the datasource, the plan cannot be empty.
				// See https://www.terraform.io/docs/configuration/data-sources.html#data-resource-dependencies
				Check: checkDatasourceStatuspageComponentsAttrs(testAccProvider, rid),
			},
		},
	})
}

func checkDatasourceStatuspageComponentsAttrs(accProvider *schema.Provider, rand int) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("data.statuspage_components.default", "page_id", pageID),
	)
}

func testAccStatuspageComponentsConfig(rand int) string {
	return fmt.Sprintf(`
	variable "component_name" {
		default = "tf-testacc-component-%d"
	}
	variable "pageid" {
		default = "%s"
	}
	resource "statuspage_component" "default" {
		page_id = "${var.pageid}"
		name = "${var.component_name}_component"
		description = "Test component 1"
		status = "operational"
	}
	`, rand, pageID)
}

func testAccDatasourceStatuspageComponentsConfig(uniq int) string {
	return fmt.Sprintf(`
	%s
	data "statuspage_components" "default" {
	depends_on = [
		statuspage_component.default,
	]

	page_id = "${var.pageid}"
	}`, testAccStatuspageComponentsConfig(uniq))
}
