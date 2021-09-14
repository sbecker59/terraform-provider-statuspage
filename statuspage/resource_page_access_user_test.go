package statuspage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	paEmail = "example@example.com"
)

func TestAccStatuspagePageAccessUser_Basic(t *testing.T) {

	rid := acctest.RandIntRange(1, 99)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStatuspagePageAccessUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPageAccessUserConfig(rid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("statuspage_page_access_user.default", "id"),
					resource.TestCheckResourceAttr("statuspage_page_access_user.default", "email", paEmail),
				),
			},
			{
				Config: testAccCheckPageAccessUserConfigUpdated(rid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("statuspage_page_access_user.default", "id"),
					resource.TestCheckResourceAttr("statuspage_page_access_user.default", "email", fmt.Sprintf("new_%s", paEmail)),
				),
			},
		},
	})
}

func testAccCheckPageAccessUserConfig(rand int) string {
	return fmt.Sprintf(`
	variable "component_name" {
		default = "tf-testacc-page-access-user-%d"
	}
	variable "pageid" {
		default = "%s"
	}
	variable "email" {
		default = "%s"
    }
	resource "statuspage_page_access_user" "default" {
		page_id     = "${var.pageid}"
		email       = "${var.email}"
	}
	`, rand, pageID, paEmail)
}

func testAccCheckPageAccessUserConfigUpdated(rand int) string {
	return fmt.Sprintf(`
	variable "component_name" {
		default = "tf-testacc-page-access-user-%d"
	}
	variable "pageid" {
		default = "%s"
	}
	variable "email" {
		default = "%s"
    }
	resource "statuspage_page_access_user" "default" {
		page_id     = "${var.pageid}"
		email       = "new_${var.email}"
	}
	`, rand, pageID, paEmail)
}

func testAccCheckStatuspagePageAccessUserDestroy(s *terraform.State) error {

	conn := testAccProvider.Meta().(*ProviderConfiguration)
	statuspageClientV1 := conn.StatuspageClientV1
	authV1 := conn.AuthV1

	for _, r := range s.RootModule().Resources {

		_, httpresp, err := statuspageClientV1.PageAccessUsersApi.GetPagesPageIdPageAccessUsersPageAccessUserId(authV1, pageID, r.Primary.ID).Execute()
		if err.Error() != "" {
			if httpresp != nil && httpresp.StatusCode == 404 {
				continue
			}
			return TranslateClientErrorDiag(err, "error retrieving component group")
		}
		return fmt.Errorf("component group still exists")
	}
	return nil

}
