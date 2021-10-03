package statuspage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccStatuspagePageAccessGroup_import(t *testing.T) {
	resourceName := "statuspage_page_access_group.default"
	rid := acctest.RandIntRange(1, 99)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStatuspageComponentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStatuspagePageAccessGroupConfigImported(rid),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(ts *terraform.State) (string, error) {
					rs := ts.RootModule().Resources["statuspage_page_access_group.default"]
					return fmt.Sprintf("%s/%s", audienceSpecificPageID, rs.Primary.ID), nil
				},
			},
		},
	})
}

func testAccStatuspagePageAccessGroupConfigImported(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testacc-access-group-%d"
	}
	variable "pageid" {
		default = "%s"
	}
	resource "statuspage_page_access_group" "default" {
		page_id = "${var.pageid}"
		name = "${var.name}"
	}
	`, rand, audienceSpecificPageID)
}
