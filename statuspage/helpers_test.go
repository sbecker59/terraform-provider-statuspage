package statuspage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestGenerateDataSourceHashID(t *testing.T) {

	schemaOneKeyOptional := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key1": {Type: schema.TypeString, Optional: true},
		},
	}
	schemaOneKeyRequired := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key1": {Type: schema.TypeString, Required: true},
		},
	}
	schemaOneKey := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key1": {Type: schema.TypeString},
		},
	}
	schemaMultipleKey := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key1": {Type: schema.TypeString, Required: true},
			"key2": {Type: schema.TypeList},
			"key3": {Type: schema.TypeSet},
			"key4": {Type: schema.TypeMap},
		},
	}

	d := &schema.ResourceData{}
	d = schema.TestResourceDataRaw(t, schemaOneKeyRequired.Schema, map[string]interface{}{
		"key1": "value1",
	})

	if element, ok := d.GetOkExists("key1"); ok {
		fmt.Printf("%v-", element)
	}

	type args struct {
		idPrefix       string
		resourceSchema *schema.Resource
		resourceData   *schema.ResourceData
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		struct {
			name string
			args args
			want string
		}{name: "resourceSchemaNilAndresourceDataNil", args: args{idPrefix: "datasource", resourceSchema: nil, resourceData: nil}, want: ""},
		{name: "oneKeyOptional", args: args{idPrefix: "datasource", resourceSchema: schemaOneKeyOptional, resourceData: &schema.ResourceData{}}, want: "datasource0"},
		{name: "oneKeyRequired", args: args{idPrefix: "datasource", resourceSchema: schemaOneKeyRequired, resourceData: &schema.ResourceData{}}, want: "datasource0"},
		{name: "oneKey", args: args{idPrefix: "datasource", resourceSchema: schemaOneKey, resourceData: &schema.ResourceData{}}, want: "datasource0"},
		{name: "multipleKey", args: args{idPrefix: "datasource", resourceSchema: schemaMultipleKey, resourceData: &schema.ResourceData{}}, want: "datasource0"},
		{name: "oneKeyAndData", args: args{idPrefix: "datasource", resourceSchema: schemaOneKeyRequired, resourceData: d}, want: "datasource482442878"},
		{name: "multipleKeyAndData", args: args{idPrefix: "datasource", resourceSchema: schemaMultipleKey, resourceData: d}, want: "datasource482442878"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateDataSourceHashID(tt.args.idPrefix, tt.args.resourceSchema, tt.args.resourceData); got != tt.want {
				t.Errorf("GenerateDataSourceHashID() = %v, want %v", got, tt.want)
			}
		})
	}
}
