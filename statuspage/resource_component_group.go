package statuspage

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sp "github.com/sbecker59/statuspage-api-client-go/api/v1/statuspage"
)

func resourceComponentGroupRead(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	name := d.Get("name").(string)
	log.Printf("[INFO] Reading Status Page component '%s'", name)

	componentGroups, _, err := statuspageClientV1.ComponentGroupsApi.GetPagesPageIdComponentGroupsId(authV1, d.Get("page_id").(string), d.Id()).Execute()
	if err.Error() != "" {
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

	components := d.Get("components").(*schema.Set).List()
	c := make([]string, len(components))
	for i, v := range components {
		c[i] = fmt.Sprint(v)
	}

	var componentGroup sp.PostPagesPageIdComponentGroupsComponentGroup

	componentGroup.SetName(name)
	componentGroup.SetDescription(description)
	componentGroup.SetComponents(c)

	o := *sp.NewPostPagesPageIdComponentGroups()
	o.SetComponentGroup(componentGroup)

	log.Printf("[INFO] Creating Status Page componant groups '%s'", name)
	resp, _, err := statuspageClientV1.ComponentGroupsApi.PostPagesPageIdComponentGroups(authV1, d.Get("page_id").(string)).PostPagesPageIdComponentGroups(o).Execute()

	if err.Error() != "" {
		return translateClientError(err, "failed to create component groups using Status Page API")
	}

	d.SetId(resp.GetId())

	return resourceComponentGroupRead(d, m)

}

func resourceComponentGroupUpdate(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	name := d.Get("name").(string)
	description := d.Get("description").(string)

	components := d.Get("components").(*schema.Set).List()
	c := make([]string, len(components))
	for i, v := range components {
		c[i] = fmt.Sprint(v)
	}

	var componentGroup sp.PostPagesPageIdComponentGroupsComponentGroup

	componentGroup.SetName(name)
	componentGroup.SetDescription(description)
	componentGroup.SetComponents(c)

	o := *sp.NewPatchPagesPageIdComponentGroups()
	o.SetComponentGroup(componentGroup)

	log.Printf("[INFO] Update Status Page componant group '%s'", name)
	resp, _, err := statuspageClientV1.ComponentGroupsApi.PatchPagesPageIdComponentGroupsId(authV1, d.Get("page_id").(string), d.Id()).PatchPagesPageIdComponentGroups(o).Execute()

	if err.Error() != "" {
		return translateClientError(err, "failed to update component group using Status Page API")
	}

	d.SetId(resp.GetId())

	return resourceComponentGroupRead(d, m)
}

func resourceComponentGroupDelete(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	_, _, err := statuspageClientV1.ComponentGroupsApi.DeletePagesPageIdComponentGroupsId(authV1, d.Get("page_id").(string), d.Id()).Execute()

	if err.Error() != "" {
		return translateClientError(err, "failed to delete component using Status Page API")
	}

	return nil
}

func resourceComponentGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if len(strings.Split(d.Id(), "/")) != 2 {
		return []*schema.ResourceData{}, fmt.Errorf("[ERROR] Invalid resource format: %s. Please use 'page-id/component-group-id'", d.Id())
	}

	pageID := strings.Split(d.Id(), "/")[0]
	componentGroupID := strings.Split(d.Id(), "/")[1]

	log.Printf("[INFO] Importing Component Group %s from Page %s", componentGroupID, pageID)

	d.Set("page_id", pageID)
	d.SetId(componentGroupID)

	err := resourceComponentGroupRead(d, m)

	return []*schema.ResourceData{d}, err
}

func resourceComponentGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceComponentGroupCreate,
		Read:   resourceComponentGroupRead,
		Update: resourceComponentGroupUpdate,
		Delete: resourceComponentGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceComponentGroupImport,
		},
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
				Type:        schema.TypeSet,
				Description: "An array with the IDs of the components in this group",
				Required:    true,
				Set:         schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}
