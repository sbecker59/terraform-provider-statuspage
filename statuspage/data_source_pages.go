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
		},
	}
}

func dataSourcePagesRead(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1
	page_name := d.Get("page_name").(string)

	res, _, err := statuspageClientV1.PagesApi.GetPages(authV1).Execute()

	if err != nil {
		return TranslateClientErrorDiag(err, "error querying pages list")
	}

	for _, r := range res {
		if _, ok := r.GetNameOk(); ok {
			if r.GetName() == page_name {
				d.SetId(r.GetId())
				break
			}
		}
	}
	return nil
}
