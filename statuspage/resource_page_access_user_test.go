package statuspage

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestUnitResourcePageAccessUserRead_FiltersByEmail(t *testing.T) {
	const email = "user@example.com"
	client, authV1 := newTestPagesClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("email") != email {
			w.Write([]byte(`[]`))
			return
		}
		w.Write([]byte(`[{"id":"user-id","email":"user@example.com"}]`))
	})

	d := schema.TestResourceDataRaw(t, resourcePageAccessUser().Schema, map[string]interface{}{
		"page_id": "page-id",
		"email":   email,
	})
	m := &ProviderConfiguration{StatuspageClientV1: client, AuthV1: authV1}

	if err := resourcePageAccessUserRead(d, m); err != nil {
		t.Fatalf("resourcePageAccessUserRead() error = %v, want nil", err)
	}
	if got, want := d.Id(), "user-id"; got != want {
		t.Errorf("Id() = %q, want %q", got, want)
	}
}

func TestAccStatuspagePageAccessUser_Basic(t *testing.T) {

	rid := acctest.RandIntRange(1, 99)
	paEmail := fmt.Sprintf("tf-testacc-page-access-user-%d@example.com", rid)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStatuspagePageAccessUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPageAccessUserConfig(rid, paEmail),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("statuspage_page_access_user.default", "id"),
					resource.TestCheckResourceAttr("statuspage_page_access_user.default", "email", paEmail),
				),
			},
			{
				Config: testAccCheckPageAccessUserConfigUpdated(rid, paEmail),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("statuspage_page_access_user.default", "id"),
					resource.TestCheckResourceAttr("statuspage_page_access_user.default", "email", fmt.Sprintf("new_%s", paEmail)),
				),
			},
		},
	})
}

func testAccCheckPageAccessUserConfig(rand int, email string) string {
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
	`, rand, audienceSpecificPageID, email)
}

func testAccCheckPageAccessUserConfigUpdated(rand int, email string) string {
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
	`, rand, audienceSpecificPageID, email)
}

func testAccCheckStatuspagePageAccessUserDestroy(s *terraform.State) error {

	conn := testAccProvider.Meta().(*ProviderConfiguration)
	statuspageClientV1 := conn.StatuspageClientV1
	authV1 := conn.AuthV1

	for _, r := range s.RootModule().Resources {

		_, httpresp, err := statuspageClientV1.PageAccessUsersAPI.GetPagesPageIdPageAccessUsersPageAccessUserId(authV1, audienceSpecificPageID, r.Primary.ID).Execute()
		if err != nil {
			if httpresp != nil && httpresp.StatusCode == 404 {
				continue
			}
			return TranslateClientErrorDiag(err, "error retrieving component group")
		}
		return fmt.Errorf("component group still exists")
	}
	return nil

}
