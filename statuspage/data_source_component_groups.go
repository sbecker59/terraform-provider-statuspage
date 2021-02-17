package statuspage

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceComponentGroups() *schema.Resource {
	return &schema.Resource{
		Description: "",
		Read:        dataSourceComponentGroupsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"page_id": {
				Description:  "the ID of the page this component belongs to",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			// Computed values
			"component_groups": {
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
						"components": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceComponentGroupsRead(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	err := providerConf.Ratelimiter.Wait(authV1) // This is a blocking call. Honors the rate limit
	if err != nil {
		return translateClientError(err, "error Ratelimiter")
	}

	res, _, err := statuspageClientV1.ComponentGroupsApi.GetPagesPageIdComponentGroups(authV1, d.Get("page_id").(string)).Execute()

	if err.Error() != "" {
		return translateClientError(err, "error querying component group list")
	}

	d.SetId(GenerateDataSourceHashID("DataSourceComponentGroups-", dataSourceComponentGroups(), d))
	resources := []map[string]interface{}{}

	for _, r := range res {
		componentGroup := map[string]interface{}{}

		if _, ok := r.GetNameOk(); ok {
			componentGroup["id"] = r.GetId()
			componentGroup["name"] = r.GetName()
			componentGroup["description"] = r.GetDescription()
			componentGroup["position"] = r.GetPosition()
			componentGroup["components"] = r.GetComponents()
		}

		resources = append(resources, componentGroup)
	}

	if f, fOk := d.GetOkExists("filter"); fOk {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceComponentGroups().Schema["component_groups"].Elem.(*schema.Resource).Schema)
	}

	if err := d.Set("component_groups", resources); err != nil {
		return err
	}

	return nil
}
