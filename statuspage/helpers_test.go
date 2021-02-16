package statuspage

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestUnitGenerateDataSourceHashID_returnHashID(t *testing.T) {
	hashID := GenerateDataSourceHashID("dataSource", &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key1": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key2": {
				Type: schema.TypeList,
			},
		},
	}, &schema.ResourceData{})
	if hashID != "dataSource0" {
		t.Error("TestUnitGenerateDataSourceHashID_resourceShemaNullAndresourceDataNull")
	}
}

func TestUnitGenerateDataSourceHashID_returnNull(t *testing.T) {
	hashID := GenerateDataSourceHashID("dataSource", nil, nil)
	if hashID != "" {
		t.Error("TestUnitGenerateDataSourceHashID_returnNull")
	}
}
