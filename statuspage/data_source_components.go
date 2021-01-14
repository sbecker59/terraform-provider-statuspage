package statuspage

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceComponents() *schema.Resource {
	return &schema.Resource{
		Description: "",
		Read:        dataSourceComponentsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"page_id": {
				Description:  "the ID of the page this component belongs to",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			// Computed values
			"components": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Required

						// Optional

						// Computed
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"position": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceComponentsRead(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	res, _, err := statuspageClientV1.ComponentsApi.GetPagesPageIdComponents(authV1, d.Get("page_id").(string)).Execute()

	if err.Error() != "" {
		return translateClientError(err, "error querying component list")
	}

	d.SetId(GenerateDataSourceHashID("DataSourceComponents-", dataSourceComponents(), d))
	resources := []map[string]interface{}{}

	for _, r := range res {
		component := map[string]interface{}{}
		fmt.Errorf("%s", r.GetName())
		if r.Name != nil {
			component["id"] = r.GetId()
			component["name"] = r.GetName()
			component["description"] = r.GetDescription()
			component["position"] = r.GetPosition()
			component["group_id"] = r.GetGroupId()
		}

		resources = append(resources, component)
	}

	if f, fOk := d.GetOkExists("filter"); fOk {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceComponents().Schema["components"].Elem.(*schema.Resource).Schema)
	}

	if err := d.Set("components", resources); err != nil {
		return err
	}

	return nil
}
