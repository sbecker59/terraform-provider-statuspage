package statuspage

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestUnitGenerateDataSourceHashID_returnHashID(t *testing.T) {
	hashID := GenerateDataSourceHashID("dataSource", &schema.Resource{}, &schema.ResourceData{})
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
