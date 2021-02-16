package statuspage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccStatuspageComponentGroup_Basic(t *testing.T) {
	rid := acctest.RandIntRange(1, 99)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStatuspageComponentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckComponentGroupConfig(rid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("statuspage_component_group.default", "id"),
					resource.TestCheckResourceAttr("statuspage_component_group.default", "description", "Acc. Tests"),
					resource.TestCheckResourceAttr("statuspage_component_group.default", "components.#", "1"),
				),
			},
			{
				Config: testAccCheckComponentGroupConfigUpdated(rid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("statuspage_component_group.default", "id"),
					resource.TestCheckResourceAttr("statuspage_component_group.default", "description", "Acc. Tests"),
					resource.TestCheckResourceAttr("statuspage_component_group.default", "components.#", "2"),
				),
			},
		},
	})
}

func testAccCheckComponentGroupConfig(rand int) string {
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

func testAccCheckComponentGroupConfigUpdated(rand int) string {
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
	resource "statuspage_component_group" "default" {
		page_id     = var.pageid
		name        = var.component_name
		description = "Acc. Tests"
		components  = ["${statuspage_component.component_1.id}", "${statuspage_component.component_2.id}"]
	}
	`, rand, pageID)
}

func testAccCheckStatuspageComponentGroupDestroy(s *terraform.State) error {

	conn := testAccProvider.Meta().(*ProviderConfiguration)
	statuspageClientV1 := conn.StatuspageClientV1
	authV1 := conn.AuthV1

	for _, r := range s.RootModule().Resources {
		_, httpresp, err := statuspageClientV1.ComponentGroupsApi.GetPagesPageIdComponentGroupsId(authV1, pageID, r.Primary.ID).Execute()
		if err.Error() != "" {
			if httpresp != nil && httpresp.StatusCode == 404 {
				continue
			}
			return translateClientError(err, "error retrieving component group")
		}
		return fmt.Errorf("component group still exists")
	}
	return nil

}
