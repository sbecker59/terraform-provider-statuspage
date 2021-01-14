package statuspage

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform/helper/hashcode"
)

func GenerateDataSourceHashID(idPrefix string, resourceSchema *schema.Resource, resourceData *schema.ResourceData) string {
	// Important, if you don't have an ID, make one up for your datasource
	// or things will end in tears.

	if resourceSchema == nil || resourceData == nil {
		return ""
	}

	var buf bytes.Buffer
	// sort keys of the map
	keys := make([]string, 0, len(resourceSchema.Schema))
	for key := range resourceSchema.Schema {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		// parse schema field is user input
		value := resourceSchema.Schema[key]
		if !(value.Required || value.Optional) {
			continue
		}

		// Ignoring TypeList, TypeSet and TypeMap
		if value.Type == schema.TypeList || value.Type == schema.TypeSet || value.Type == schema.TypeMap {
			continue
		}

		if element, ok := resourceData.GetOkExists(key); ok {
			buf.WriteString(fmt.Sprintf("%v-", element))
		}
	}
	return fmt.Sprintf("%s%d", idPrefix, hashcode.String(buf.String()))
}
