package statuspage

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccStatuspageSubscriber_Basic(t *testing.T) {

	time.Sleep(10 * time.Second)
	rid := acctest.RandIntRange(1, 99)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStatuspageSubscriberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSubscriberConfig(rid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("statuspage_subscriber.default", "id"),
					resource.TestCheckResourceAttr("statuspage_subscriber.default", "email", fmt.Sprintf("email-%d@testacc.tf", rid)),
				),
			},
			{
				Config: testAccCheckSubscriberConfigUpdated(rid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("statuspage_subscriber.default", "id"),
					resource.TestCheckResourceAttr("statuspage_subscriber.default", "email", fmt.Sprintf("email-updated-%d@testacc.tf", rid)),
				),
			},
		},
	})
}

func testAccCheckSubscriberConfig(rand int) string {
	return fmt.Sprintf(`
	variable "email" {
		default = "email-%d@testacc.tf"
	}
	variable "pageid" {
		default = "%s"
	}
	resource "statuspage_subscriber" "default" {
		page_id = var.pageid
		email = var.email
	}
	`, rand, pageID)
}
func testAccCheckSubscriberConfigUpdated(rand int) string {
	return fmt.Sprintf(`
	variable "email" {
		default = "email-updated-%d@testacc.tf"
	}
	variable "pageid" {
		default = "%s"
	}
	resource "statuspage_subscriber" "default" {
		page_id = var.pageid
		email = var.email
	}
	`, rand, pageID)
}

func testAccCheckStatuspageSubscriberDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ProviderConfiguration)
	statuspageClientV1 := conn.StatuspageClientV1
	authV1 := conn.AuthV1

	conn.Ratelimiter.Wait(authV1)

	for _, r := range s.RootModule().Resources {
		_, httpresp, err := statuspageClientV1.SubscribersApi.GetPagesPageIdSubscribersSubscriberId(authV1, pageID, r.Primary.ID).Execute()
		if err.Error() != "" {
			if httpresp != nil && httpresp.StatusCode == 404 {
				continue
			}
			return translateClientError(err, "error retrieving subscriber")
		}
		return fmt.Errorf("component subscriber still exists")
	}
	return nil
}
