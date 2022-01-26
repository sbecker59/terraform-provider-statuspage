package statuspage

import (
	"fmt"

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

	fmt.Printf(page_name)
	// res, _, err := statuspageClientV1.ComponentsApi.GetPagesPageIdComponents(authV1, d.Get("page_id").(string)).Execute()
	req := statuspageClientV1.PagesApi.GetPages(authV1)

	res,_,_  := req.Execute()

	

	// if err.Error() != "" {
	// 	return TranslateClientErrorDiag(err, "error querying component list")
	// }

	// d.SetId(GenerateDataSourceHashID("DataSourcePages-", dataSourcePages(), d))
	resources := []map[string]interface{}{}

	for _, r := range res  {
		

	
		if _, ok := r.GetNameOk(); ok  {
			if r.GetName() == page_name{
				d.SetId(r.GetId())
				// pages["id"] = r.GetId()
				// pages["name"] = r.GetName()
				// resources = append(resources, pages)
				break
			} 
			
			


			
		}
		

	}

	if f, fOk := d.GetOkExists("filter"); fOk {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourcePages().Schema["pages"].Elem.(*schema.Resource).Schema)
	}

	// if err := d.Set("pages", resources); err != nil {
	// 	return err
	// }

	return nil
}
