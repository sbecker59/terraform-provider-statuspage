package statuspage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccStatuspagePageAccessGroup_Basic(t *testing.T) {

	rid := acctest.RandIntRange(1, 99)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStatuspagePageAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPageAccessGroupConfig(rid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("statuspage_page_access_group.default", "id"),
					resource.TestCheckResourceAttr("statuspage_page_access_group.default", "name", "Test Access Group"),
					resource.TestCheckResourceAttr("statuspage_page_access_group.default", "users.#", "1"),
					resource.TestCheckResourceAttr("statuspage_page_access_group.default", "components.#", "1"),
				),
			},
			{
				Config: testAccCheckPageAccessGroupConfigUpdated(rid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("statuspage_page_access_group.default", "id"),
					resource.TestCheckResourceAttr("statuspage_page_access_group.default", "name", "Test Access Group"),
					resource.TestCheckResourceAttr("statuspage_page_access_group.default", "users.#", "2"),
					resource.TestCheckResourceAttr("statuspage_page_access_group.default", "components.#", "2"),
				),
			},
		},
	})
}

func testAccCheckPageAccessGroupConfig(rand int) string {
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

	resource "statuspage_page_access_user" "user_1" {
		page_id     = "${var.pageid}"
		email       = "${var.component_name}@example.com"
    }

	resource "statuspage_page_access_group" "default" {
		page_id     = "${var.pageid}"
		name        = "Test Access Group"
		users       = ["${statuspage_page_access_user.user_1.id}"]
		components  = ["${statuspage_component.component_1.id}"]
	}
	`, rand, pageID)
}

func testAccCheckPageAccessGroupConfigUpdated(rand int) string {
	return fmt.Sprintf(`
	variable "component_name" {
		default = "tf-testacc-component-group-%d"
	}
	variable "pageid" {
		default = "%s"
	}
	resource "statuspage_component" "component_1" {
		page_id = var.pageid
		name = "${var.component_name}_component"
		description = "Test component 1"
		status = "operational"
	}
	resource "statuspage_component" "component_2" {
		page_id = var.pageid
		name = "${var.component_name}_component"
		description = "Test component 2"
		status = "operational"
	}
	resource "statuspage_page_access_user" "user_1" {
		page_id     = "${var.pageid}"
		email       = "${var.component_name}@example.com"
    }
	resource "statuspage_page_access_user" "user_2" {
		page_id     = "${var.pageid}"
		email       = "${var.component_name}-two@example.com"
    }
	resource "statuspage_page_access_group" "default" {
		page_id     = "${var.pageid}"
		name        = "Test Access Group"
		users       = ["${statuspage_page_access_user.user_1.id}", "${statuspage_page_access_user.user_2.id}"]
		components  = ["${statuspage_component.component_1.id}", "${statuspage_component.component_2.id}"]
	}
	`, rand, pageID)
}

func testAccCheckStatuspagePageAccessGroupDestroy(s *terraform.State) error {

	conn := testAccProvider.Meta().(*ProviderConfiguration)
	statuspageClientV1 := conn.StatuspageClientV1
	authV1 := conn.AuthV1

	for _, r := range s.RootModule().Resources {

		_, httpresp, err := statuspageClientV1.PageAccessGroupsApi.GetPagesPageIdPageAccessGroupsPageAccessGroupId(authV1, pageID, r.Primary.ID).Execute()
		if err.Error() != "" {
			if httpresp != nil && httpresp.StatusCode == 404 {
				continue
			}
			return TranslateClientErrorDiag(err, "error retrieving page access group")
		}
		return fmt.Errorf("page access group still exists")
	}
	return nil

}
