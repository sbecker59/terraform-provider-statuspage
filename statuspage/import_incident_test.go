package statuspage

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestStatuspageIncident_import(t *testing.T) {
	resourceName := "statuspage_incident.default"
	time.Sleep(10 * time.Second)
	rid := acctest.RandIntRange(1, 99)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStatuspageIncidentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckStatuspageIncidentConfigImported(rid),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(ts *terraform.State) (string, error) {
					rs := ts.RootModule().Resources["statuspage_incident.default"]
					return fmt.Sprintf("%s/%s", pageID, rs.Primary.ID), nil
				},
			},
		},
	})
}

func testAccCheckStatuspageIncidentConfigImported(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testacc-component-%d"
	}
	variable "pageid" {
		default = "%s"
	}
	resource "statuspage_incident" "default" {
		page_id = "${var.pageid}"
		name = "${var.name}"
		impact_override = "critical"
		status = "identified"
	}
	`, rand, pageID)
}
