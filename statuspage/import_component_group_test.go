package statuspage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestStatuspageComponentGroup_import(t *testing.T) {
	resourceName := "statuspage_component_group.default"
	rid := acctest.RandIntRange(1, 99)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStatuspageComponentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckStatuspageComponentGroupConfigImported(rid),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(ts *terraform.State) (string, error) {
					rs := ts.RootModule().Resources["statuspage_component_group.default"]
					return fmt.Sprintf("%s/%s", pageID, rs.Primary.ID), nil
				},
			},
		},
	})
}

func testAccCheckStatuspageComponentGroupConfigImported(rand int) string {
	return fmt.Sprintf(`
	variable "component_name" {
		default = "tf-testacc-component-group-%d"
	}
	variable "pageid" {
		default = "%s"
	}
	resource "statuspage_component" "comp1" {
		page_id = "${var.pageid}"
		name = "${var.component_name}_component"
		description = "test component"
		status = "operational"
	}
	resource "statuspage_component_group" "default" {
		page_id     = "${var.pageid}"
		name        = "${var.component_name}"
		description = "Acc. Tests"
		components  = ["${statuspage_component.comp1.id}"]
	}
	`, rand, pageID)
}
