package statuspage

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccStatuspageComponentGroupsDatasource(t *testing.T) {

	time.Sleep(10 * time.Second)
	rid := acctest.RandIntRange(1, 99)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStatuspageComponentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceStatuspageComponentGroupsConfig(rid),
				// Because of the `depends_on` in the datasource, the plan cannot be empty.
				// See https://www.terraform.io/docs/configuration/data-sources.html#data-resource-dependencies
				Check: checkDatasourceStatuspageComponentGroupsAttrs(testAccProvider, rid),
			},
		},
	})
}

func checkDatasourceStatuspageComponentGroupsAttrs(accProvider *schema.Provider, rand int) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("data.statuspage_component_groups.default", "page_id", pageID),
	)
}

func testAccStatuspageComponentGroupConfig(rand int) string {
	return fmt.Sprintf(`
	variable "component_name" {
		default = "tf-testacc-component-group-%d"
	}
	variable "pageid" {
		default = "%s"
	}
	resource "statuspage_component" "component_1" {
		page_id = "${var.pageid}"
		name = "${var.component_name}_component"
		description = "Test component 1"
		status = "operational"
	}
	resource "statuspage_component_group" "default" {
		page_id     = "${var.pageid}"
		name        = "${var.component_name}"
		description = "Acc. Tests"
		components  = ["${statuspage_component.component_1.id}"]
	}
	`, rand, pageID)
}

func testAccDatasourceStatuspageComponentGroupsConfig(uniq int) string {
	return fmt.Sprintf(`
	%s
	data "statuspage_component_groups" "default" {
	depends_on = [
		statuspage_component_group.default,
	]

	page_id = "${var.pageid}"
	}`, testAccStatuspageComponentGroupConfig(uniq))
}
