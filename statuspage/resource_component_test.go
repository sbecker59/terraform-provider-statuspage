package statuspage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccStatuspageComponent_Basic(t *testing.T) {
	rid := acctest.RandIntRange(1, 99)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckComponentConfig(rid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("statuspage_component.default", "id"),
					resource.TestCheckResourceAttr("statuspage_component.default", "description", "test component"),
					resource.TestCheckResourceAttr("statuspage_component.default", "status", "operational"),
					resource.TestCheckResourceAttr("statuspage_component.default", "showcase", "true"),
				),
			},
			{
				Config: testAccCheckComponentConfigUpdated(rid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("statuspage_component.default", "id"),
					resource.TestCheckResourceAttr("statuspage_component.default", "description", "updated component"),
					resource.TestCheckResourceAttr("statuspage_component.default", "status", "major_outage"),
					resource.TestCheckResourceAttr("statuspage_component.default", "showcase", "false"),
				),
			},
		},
	})
}

func testAccCheckComponentConfig(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testacc-component-%d"
	}
	variable "pageid" {
		default = "%s"
	}
	resource "statuspage_component" "default" {
		page_id = var.pageid
		name = var.name
		description = "test component"
		status = "operational"
		showcase = true
	}
	`, rand, pageId)
}

func testAccCheckComponentConfigUpdated(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testacc-component-%d"
	}
	variable "pageid" {
		default = "%s"
	}
	resource "statuspage_component" "default" {
		page_id = var.pageid
		name = var.name
		description = "updated component"
		status = "major_outage"
		showcase = false
	}
	`, rand, pageId)
}
