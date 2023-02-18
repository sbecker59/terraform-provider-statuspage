package statuspage

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sp "github.com/sbecker59/statuspage-api-client-go/api/v1/statuspage"
)

func resourcePageAccessUserRead(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	email := d.Get("email").(string)
	log.Printf("[INFO] Looking up user by email '%s'", email)

	pageAccessUsers, _, err := statuspageClientV1.PageAccessUsersApi.GetPagesPageIdPageAccessUsers(authV1, d.Get("page_id").(string)).Page(1).PerPage(100).Execute()

	if err.Error() != "" {
		return TranslateClientErrorDiag(err, "failed to get page access users using Status Page API")
	}

	if &pageAccessUsers == nil {
		d.SetId("")
		return nil
	}

	for _, u := range pageAccessUsers {
		if email == *u.Email {
			d.SetId(*u.Id)
			break
		} else {
			d.SetId("")
		}
	}

	return nil
}

func resourcePageAccessUserCreate(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	email := d.Get("email").(string)

	var pageAccessUser sp.PostPagesPageIdPageAccessUsersPageAccessUser

	pageAccessUser.SetEmail(email)

	o := *sp.NewPostPagesPageIdPageAccessUsers()
	o.SetPageAccessUser(pageAccessUser)

	log.Printf("[INFO] Creating Status Page access user '%s'", email)
	resp, _, err := statuspageClientV1.PageAccessUsersApi.PostPagesPageIdPageAccessUsers(authV1, d.Get("page_id").(string)).PostPagesPageIdPageAccessUsers(o).Execute()

	if err.Error() != "" {
		return TranslateClientErrorDiag(err, "failed to create page access user using Status Page API")
	}

	d.SetId(resp.GetId())

	return resourcePageAccessUserRead(d, m)

}

func resourcePageAccessUserDelete(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	_, err := statuspageClientV1.PageAccessUsersApi.DeletePagesPageIdPageAccessUsersPageAccessUserId(authV1, d.Get("page_id").(string), d.Id()).Execute()

	if err.Error() != "" {
		return TranslateClientErrorDiag(err, "failed to delete page access user using Status Page API")
	}

	return nil
}

func resourcePageAccessUserImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if len(strings.Split(d.Id(), "/")) != 2 {
		return []*schema.ResourceData{}, fmt.Errorf("[ERROR] Invalid resource format: %s. Please use 'page-id/email-address'", d.Id())
	}

	pageID := strings.Split(d.Id(), "/")[0]
	email := strings.Split(d.Id(), "/")[1]

	log.Printf("[INFO] Importing Page Access User %s from Page %s", email, pageID)

	d.Set("page_id", pageID)
	d.Set("email", email)

	err := resourcePageAccessUserRead(d, m)

	return []*schema.ResourceData{d}, err
}

func resourcePageAccessUser() *schema.Resource {
	return &schema.Resource{
		Create: resourcePageAccessUserCreate,
		Read:   resourcePageAccessUserRead,
		Delete: resourcePageAccessUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePageAccessUserImport,
		},
		Schema: map[string]*schema.Schema{
			"page_id": {
				Type:        schema.TypeString,
				Description: "the ID of the page this user belongs to",
				Required:    true,
				ForceNew:    true,
			},
			"email": {
				Type:        schema.TypeString,
				Description: "The email of the user",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}
