package statuspage

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sp "github.com/sbecker59/statuspage-api-client-go/api/v1/statuspage"
)

func resourcePageAccessGroupRead(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	name := d.Get("name").(string)
	log.Printf("[INFO] Reading Status Page component '%s'", name)

	pageAccessGroups, _, err := statuspageClientV1.PageAccessGroupsApi.GetPagesPageIdPageAccessGroupsPageAccessGroupId(authV1, d.Get("page_id").(string), d.Id()).Execute()

	if err.Error() != "" {
		return TranslateClientErrorDiag(err, "failed to get component groups using Status Page API")
	}

	if &pageAccessGroups == nil {
		d.SetId("")
		return nil
	}

	d.Set("external_identifier", pageAccessGroups.ExternalIdentifier)
	d.Set("name", pageAccessGroups.Name)
	d.Set("components", pageAccessGroups.ComponentIds)
	d.Set("metrics", pageAccessGroups.MetricIds)
	d.Set("users", pageAccessGroups.PageAccessUserIds)

	return nil
}

func resourcePageAccessGroupCreate(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	name := d.Get("name").(string)
	externalIdentifier := d.Get("external_identifier").(string)

	var pageAccessGroup sp.PostPagesPageIdPageAccessGroupsPageAccessGroup

	pageAccessGroup.SetName(name)
	pageAccessGroup.SetExternalIdentifier(externalIdentifier)
	pageAccessGroup.SetComponentIds(StringListFromSchemaKey(d, "components"))
	pageAccessGroup.SetMetricIds(StringListFromSchemaKey(d, "metrics"))
	pageAccessGroup.SetPageAccessUserIds(StringListFromSchemaKey(d, "users"))

	o := *sp.NewPostPagesPageIdPageAccessGroups()
	o.SetPageAccessGroup(pageAccessGroup)

	log.Printf("[INFO] Creating Status Page componant groups '%s'", name)
	resp, _, err := statuspageClientV1.PageAccessGroupsApi.PostPagesPageIdPageAccessGroups(authV1, d.Get("page_id").(string)).PostPagesPageIdPageAccessGroups(o).Execute()

	if err.Error() != "" {
		return TranslateClientErrorDiag(err, "failed to create component groups using Status Page API")
	}

	d.SetId(resp.GetId())

	return resourcePageAccessGroupRead(d, m)

}

func resourcePageAccessGroupUpdate(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	name := d.Get("name").(string)
	externalIdentifier := d.Get("external_identifier").(string)

	var pageAccessGroup sp.PostPagesPageIdPageAccessGroupsPageAccessGroup

	pageAccessGroup.SetName(name)
	pageAccessGroup.SetExternalIdentifier(externalIdentifier)
	pageAccessGroup.SetComponentIds(StringListFromSchemaKey(d, "components"))
	pageAccessGroup.SetMetricIds(StringListFromSchemaKey(d, "metrics"))
	pageAccessGroup.SetPageAccessUserIds(StringListFromSchemaKey(d, "users"))

	o := *sp.NewPatchPagesPageIdPageAccessGroups()
	o.SetPageAccessGroup(pageAccessGroup)

	log.Printf("[INFO] Update Status Page componant group '%s'", name)
	resp, _, err := statuspageClientV1.PageAccessGroupsApi.PatchPagesPageIdPageAccessGroupsPageAccessGroupId(authV1, d.Get("page_id").(string), d.Id()).PatchPagesPageIdPageAccessGroups(o).Execute()

	if err.Error() != "" {
		return TranslateClientErrorDiag(err, "failed to update component group using Status Page API")
	}

	d.SetId(resp.GetId())

	return resourcePageAccessGroupRead(d, m)
}

func resourcePageAccessGroupDelete(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	_, _, err := statuspageClientV1.PageAccessGroupsApi.DeletePagesPageIdPageAccessGroupsPageAccessGroupId(authV1, d.Get("page_id").(string), d.Id()).Execute()

	if err.Error() != "" {
		return TranslateClientErrorDiag(err, "failed to delete page access group using Status Page API")
	}

	return nil
}

func resourcePageAccessGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if len(strings.Split(d.Id(), "/")) != 2 {
		return []*schema.ResourceData{}, fmt.Errorf("[ERROR] Invalid resource format: %s. Please use 'page-id/component-group-id'", d.Id())
	}

	pageID := strings.Split(d.Id(), "/")[0]
	pageAccessGroupID := strings.Split(d.Id(), "/")[1]

	log.Printf("[INFO] Importing Page Access Group %s from Page %s", pageAccessGroupID, pageID)

	d.Set("page_id", pageID)
	d.SetId(pageAccessGroupID)

	err := resourcePageAccessGroupRead(d, m)

	return []*schema.ResourceData{d}, err
}

func resourcePageAccessGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourcePageAccessGroupCreate,
		Read:   resourcePageAccessGroupRead,
		Update: resourcePageAccessGroupUpdate,
		Delete: resourcePageAccessGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePageAccessGroupImport,
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
				Description: "Name for this Group.",
				Required:    true,
			},
			"external_identifier": {
				Type:        schema.TypeString,
				Description: "Associates group with external group",
				Optional:    true,
			},
			"components": {
				Type:        schema.TypeSet,
				Description: "An array with the IDs of the components in this group",
				Optional:    true,
				Set:         schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"metrics": {
				Type:        schema.TypeSet,
				Description: "An array with the IDs of the metrics in this group",
				Optional:    true,
				Set:         schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"users": {
				Type:        schema.TypeSet,
				Description: "An array with the Page Access User IDs that are in this group",
				Optional:    true,
				Set:         schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}
