package statuspage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccStatuspageComponentGroup_Basic(t *testing.T) {
	rid := acctest.RandIntRange(1, 99)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
	`, rand, pageId)
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
	`, rand, pageId)
}
