package statuspage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccStatuspageIncident_Basic(t *testing.T) {

	rid := acctest.RandIntRange(1, 99)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStatuspageIncidentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncidentConfig(rid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("statuspage_incident.default", "id"),
					resource.TestCheckResourceAttr("statuspage_incident.default", "impact_override", "maintenance"),
					resource.TestCheckResourceAttr("statuspage_incident.default", "status", "investigating"),
					resource.TestCheckResourceAttr("statuspage_incident.default", "body", "-"),
				),
			},
			{
				Config: testAccCheckIncidentConfigUpdated(rid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("statuspage_incident.default", "id"),
					resource.TestCheckResourceAttr("statuspage_incident.default", "impact_override", "critical"),
					resource.TestCheckResourceAttr("statuspage_incident.default", "status", "identified"),
					resource.TestCheckResourceAttr("statuspage_incident.default", "body", "-"),
				),
			},
		},
	})
}

func testAccCheckIncidentConfig(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testacc-Incident-%d"
	}
	variable "pageid" {
		default = "%s"
	}
	resource "statuspage_incident" "default" {
		page_id = var.pageid
		name = var.name
		impact_override = "maintenance"
		status = "investigating"
		body = "-"
	}
	`, rand, pageID)
}

func testAccCheckIncidentConfigUpdated(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testacc-Incident-%d"
	}
	variable "pageid" {
		default = "%s"
	}

	resource "statuspage_component" "my_component" {
		page_id     = var.pageid
		name        = var.name
		status      = "operational"
	}

	resource "statuspage_incident" "default" {
		page_id = var.pageid
		name = var.name
		impact_override = "critical"
		status = "identified"
		body = "-"

		component {
			id = statuspage_component.my_component.id
			name = statuspage_component.my_component.name
			status = "operational"
		}
	}
	`, rand, pageID)
}

func testAccCheckStatuspageIncidentDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ProviderConfiguration)
	statuspageClientV1 := conn.StatuspageClientV1
	authV1 := conn.AuthV1

	for _, r := range s.RootModule().Resources {
		_, httpresp, err := statuspageClientV1.IncidentsApi.GetPagesPageIdIncidentsIncidentId(authV1, pageID, r.Primary.ID).Execute()
		if err.Error() != "" {
			if httpresp != nil && httpresp.StatusCode == 404 {
				continue
			}
			return TranslateClientErrorDiag(err, "error retrieving Incident")
		}
		return fmt.Errorf("Incident still exists")
	}
	return nil
}
