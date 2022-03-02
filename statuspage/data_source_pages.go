package statuspage

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourcePages() *schema.Resource {
	return &schema.Resource{
		Description: "",
		Read:        dataSourcePagesRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"page_name": {
				Description:  "the name of the page to",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			// Computed values
			// "pages": {
			// 	Type:     schema.TypeList,
			// 	Computed: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			// Required

			// 			// Optional

			// 			// Computed
			// 			"id": {
			// 				Type:     schema.TypeString,
			// 				Computed: true,
			// 			},
			// 			"name": {
			// 				Type:     schema.TypeString,
			// 				Computed: true,
			// 			},
						
			// 		},
			// 	},
			// },
		},
	}
}

func dataSourcePagesRead(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1
	page_name := d.Get("page_name").(string)

	req := statuspageClientV1.PagesApi.GetPages(authV1)

	res,_,_  := req.Execute()

	for _, r := range res  {
		if _, ok := r.GetNameOk(); ok  {
			if r.GetName() == page_name{
				d.SetId(r.GetId())
				break
			}
		}
	}
	return nil
}
