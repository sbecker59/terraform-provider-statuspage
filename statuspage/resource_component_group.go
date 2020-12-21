package statuspage

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sp "github.com/sbecker59/statuspage-api-client-go/api/v1/statuspage"
)

func resourceComponentGroupRead(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	name := d.Get("name").(string)
	log.Printf("[INFO] Reading OpsGenie component '%s'", name)

	componentGroups, _, err := statuspageClientV1.ComponentGroupsApi.GetPagesPageIdComponentGroupsId(authV1, d.Get("page_id").(string), d.Id())
	if err != nil {
		return translateClientError(err, "failed to get component groups using Status Page API")
	}

	if &componentGroups == nil {
		d.SetId("")
		return nil
	}

	d.Set("description", componentGroups.Description)
	d.Set("name", componentGroups.Name)
	d.Set("components", componentGroups.Components)

	return nil
}

func resourceComponentGroupCreate(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	name := d.Get("name").(string)
	description := d.Get("description").(string)

	components := d.Get("components").([]interface{})
	c := make([]string, len(components))
	for i, v := range components {
		c[i] = fmt.Sprint(v)
	}

	p := sp.PostPagesPageIdComponentGroups{
		Description: description,
		ComponentGroup: &sp.PostPagesPageIdComponentGroupsComponentGroup{
			Name:       name,
			Components: c,
		},
	}

	log.Printf("[INFO] Creating Status Page componant groups '%s'", name)
	result, _, err := statuspageClientV1.ComponentGroupsApi.PostPagesPageIdComponentGroups(authV1, d.Get("page_id").(string), p)

	if err != nil {
		return translateClientError(err, "failed to create component groups using Status Page API")
	}

	d.SetId(result.Id)

	return resourceComponentGroupRead(d, m)

}

func resourceComponentGroupUpdate(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	name := d.Get("name").(string)
	description := d.Get("description").(string)

	components := d.Get("components").([]interface{})
	c := make([]string, len(components))
	for i, v := range components {
		c[i] = fmt.Sprint(v)
	}

	p := sp.PatchPagesPageIdComponentGroups{
		Description: description,
		ComponentGroup: &sp.PostPagesPageIdComponentGroupsComponentGroup{
			Name:       name,
			Components: c,
		},
	}

	log.Printf("[INFO] Update Status Page componant group '%s'", name)
	result, _, err := statuspageClientV1.ComponentGroupsApi.PatchPagesPageIdComponentGroupsId(authV1, d.Get("page_id").(string), d.Id(), p)

	if err != nil {
		return translateClientError(err, "failed to update component group using Status Page API")
	}

	d.SetId(result.Id)

	return resourceComponentGroupRead(d, m)
}

func resourceComponentGroupDelete(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	_, _, err := statuspageClientV1.ComponentGroupsApi.DeletePagesPageIdComponentGroupsId(authV1, d.Get("page_id").(string), d.Id())

	if err != nil {
		return translateClientError(err, "failed to delete component using Status Page API")
	}

	return nil
}

func resourceComponentGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceComponentGroupCreate,
		Read:   resourceComponentGroupRead,
		Update: resourceComponentGroupUpdate,
		Delete: resourceComponentGroupDelete,
		Schema: map[string]*schema.Schema{
			"page_id": {
				Type:        schema.TypeString,
				Description: "the ID of the page this component group belongs to",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "An array with the IDs of the components in this group",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "More detailed description for this component group",
				Optional:    true,
			},
			"components": {
				Type:        schema.TypeList,
				Description: "An array with the IDs of the components in this group",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}
