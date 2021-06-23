package statuspage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestStatuspageComponent_import(t *testing.T) {
	resourceName := "statuspage_component.default"
	rid := acctest.RandIntRange(1, 99)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStatuspageComponentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckStatuspageComponentConfigImported(rid),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(ts *terraform.State) (string, error) {
					rs := ts.RootModule().Resources["statuspage_component.default"]
					return fmt.Sprintf("%s/%s", pageID, rs.Primary.ID), nil
				},
			},
		},
	})
}

func testAccCheckStatuspageComponentConfigImported(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testacc-component-%d"
	}
	variable "pageid" {
		default = "%s"
	}
	resource "statuspage_component" "default" {
		page_id = "${var.pageid}"
		name = "${var.name}"
		description = "test component"
		status = "operational"
		showcase = true
	}
	`, rand, pageID)
}
