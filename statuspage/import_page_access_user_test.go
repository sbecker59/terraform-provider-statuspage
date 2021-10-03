package statuspage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccStatuspagePageAccessUser_import(t *testing.T) {
	resourceName := "statuspage_page_access_user.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStatuspageComponentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStatuspagePageAccessUserConfigImported(),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(ts *terraform.State) (string, error) {
					rs := ts.RootModule().Resources["statuspage_page_access_user.default"]
					return fmt.Sprintf("%s/%s", audienceSpecificPageID, rs.Primary.ID), nil
				},
			},
		},
	})
}

func testAccStatuspagePageAccessUserConfigImported() string {
	return fmt.Sprintf(`
	variable "name" {
		default = "example@example.com"
	}
	variable "pageid" {
		default = "%s"
	}
	resource "statuspage_page_access_user" "default" {
		page_id = "${var.pageid}"
		email = "${var.name}"
	}
	`, audienceSpecificPageID)
}
