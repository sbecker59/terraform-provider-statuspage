package statuspage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccStatuspageSubscriber_Basic(t *testing.T) {
	rid := acctest.RandIntRange(1, 99)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSubscriberConfig(rid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("statuspage_subscriber.default", "id"),
					resource.TestCheckResourceAttr("statuspage_subscriber.default", "email", fmt.Sprintf("email-%d@testacc.tf", rid)),
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
	`, rand, pageId)
}
