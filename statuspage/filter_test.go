package statuspage

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Not supplying filters should not restrict results
func TestUnitApplyFilters_passThrough(t *testing.T) {
	items := []map[string]interface{}{
		{},
		{},
		{},
	}
	testSchema := map[string]*schema.Schema{
		"letter": {
			Type: schema.TypeString,
		},
	}

	res := ApplyFilters(nil, items, testSchema)
	if len(res) != 3 {
		t.Errorf("Expected 3 results, got %d", len(res))
	}
}

// Filtering against a nonexistent property should throw no errors and return no results
func TestUnitApplyFilters_nonExistentProperty(t *testing.T) {
	items := []map[string]interface{}{
		{"letter": "a"},
	}
	testSchema := map[string]*schema.Schema{
		"letter": {
			Type: schema.TypeString,
		},
	}

	filters := &schema.Set{F: func(interface{}) int { return 1 }}
	filters.Add(map[string]interface{}{
		"name":   "number",
		"values": []interface{}{"1"},
	})

	res := ApplyFilters(filters, items, testSchema)
	if len(res) > 0 {
		t.Errorf("Expected 0 results, got %d", len(res))
	}
}

// Filtering against an empty resource set should not throw errors
func TestUnitApplyFilters_noResources(t *testing.T) {
	items := []map[string]interface{}{}

	testSchema := map[string]*schema.Schema{
		"number": {
			Type: schema.TypeString,
		},
	}

	filters := &schema.Set{F: func(interface{}) int { return 1 }}
	filters.Add(map[string]interface{}{
		"name":   "number",
		"values": []interface{}{"1"},
	})

	res := ApplyFilters(filters, items, testSchema)
	if len(res) != 0 {
		t.Errorf("Expected 0 results, got %d", len(res))
	}
}

func TestUnitApplyFilters_basic(t *testing.T) {
	items := []map[string]interface{}{
		{"letter": "a"},
		{"letter": "b"},
		{"letter": "c"},
	}

	testSchema := map[string]*schema.Schema{
		"letter": {
			Type: schema.TypeString,
		},
	}

	filters := &schema.Set{F: func(interface{}) int { return 1 }}
	filters.Add(map[string]interface{}{
		"name":   "letter",
		"values": []interface{}{"b"},
	})

	res := ApplyFilters(filters, items, testSchema)
	if len(res) != 1 {
		t.Errorf("Expected 1 result, got %d", len(res))
	}
}

func TestUnitApplyFilters_duplicates(t *testing.T) {
	items := []map[string]interface{}{
		{"letter": "a"},
		{"letter": "a"},
		{"letter": "c"},
	}
	testSchema := map[string]*schema.Schema{
		"letter": {
			Type: schema.TypeString,
		},
	}

	filters := &schema.Set{F: func(v interface{}) int {
		return schema.HashString(v.(map[string]interface{})["name"])
	}}
	filters.Add(map[string]interface{}{
		"name":   "letter",
		"values": []interface{}{"a"},
	})

	res := ApplyFilters(filters, items, testSchema)
	if len(res) != 2 {
		t.Errorf("Expected 2 results, got %d", len(res))
	}
}

func TestUnitApplyFilters_OR(t *testing.T) {
	items := []map[string]interface{}{
		{"letter": "a"},
		{"letter": "b"},
		{"letter": "c"},
	}

	testSchema := map[string]*schema.Schema{
		"letter": {
			Type: schema.TypeString,
		},
	}

	filters := &schema.Set{
		F: func(v interface{}) int {
			elems := v.(map[string]interface{})["values"].([]interface{})
			res := make([]string, len(elems))
			for i, v := range elems {
				res[i] = v.(string)
			}
			return schema.HashString(strings.Join(res, ""))
		},
	}
	filters.Add(map[string]interface{}{
		"name":   "letter",
		"values": []interface{}{"a", "b"},
	})

	res := ApplyFilters(filters, items, testSchema)
	if len(res) != 2 {
		t.Errorf("Expected 2 results, got %d", len(res))
	}
}

func TestUnitApplyFilters_cascadeAND(t *testing.T) {
	items := []map[string]interface{}{
		{"letter": "a"},
		{"letter": "b"},
		{"letter": "c"},
	}
	testSchema := map[string]*schema.Schema{
		"letter": {
			Type: schema.TypeString,
		},
	}

	filters := &schema.Set{
		F: func(v interface{}) int {
			elems := v.(map[string]interface{})["values"].([]interface{})
			res := make([]string, len(elems))
			for i, v := range elems {
				res[i] = v.(string)
			}
			return schema.HashString(strings.Join(res, ""))
		},
	}
	filters.Add(map[string]interface{}{
		"name":   "letter",
		"values": []interface{}{"a", "b"},
	})
	filters.Add(map[string]interface{}{
		"name":   "letter",
		"values": []interface{}{"c"},
	})

	res := ApplyFilters(filters, items, testSchema)
	if len(res) != 0 {
		t.Errorf("Expected 0 results, got %d", len(res))
	}
}

func TestUnitApplyFilters_regex(t *testing.T) {
	items := []map[string]interface{}{
		{"string": "xblx:PHX-AD-1"},
		{"string": "xblx:PHX-AD-2"},
		{"string": "xblx:PHX-AD-3"},
	}

	testSchema := map[string]*schema.Schema{
		"string": {
			Type: schema.TypeString,
		},
	}

	filters := &schema.Set{F: func(v interface{}) int {
		return schema.HashString(v.(map[string]interface{})["name"])
	}}
	filters.Add(map[string]interface{}{
		"name":   "string",
		"values": []interface{}{"\\w*:PHX-AD-2"},
		"regex":  true,
	})

	res := ApplyFilters(filters, items, testSchema)
	if len(res) != 1 {
		t.Errorf("Expected 1 result, got %d", len(res))
	}
}

// Filters should test against an array of strings
func TestUnitApplyFilters_arrayOfStrings(t *testing.T) {
	items := []map[string]interface{}{
		{"letters": []string{"a"}},
		{"letters": []string{"b", "c"}},
		{"letters": []string{"c", "d", "e"}},
		{"letters": []string{"e", "f"}},
	}

	testSchema := map[string]*schema.Schema{
		"letters": {
			Type: schema.TypeList,
			Elem: schema.TypeString,
		},
	}

	filters := &schema.Set{F: func(interface{}) int { return 1 }}
	filters.Add(map[string]interface{}{
		"name":   "letters",
		"values": []interface{}{"a", "c"},
	})

	res := ApplyFilters(filters, items, testSchema)
	if len(res) != 3 {
		t.Errorf("Expected 3 result, got %d", len(res))
	}

	filters = &schema.Set{F: func(interface{}) int { return 1 }}
	filters.Add(map[string]interface{}{
		"name":   "letters",
		"values": []interface{}{"a", "f"},
	})

	res = ApplyFilters(filters, items, testSchema)
	if len(res) != 2 {
		t.Errorf("Expected 2 result, got %d", len(res))
	}
}

type CustomStringTypeA string
type CustomStringTypeB CustomStringTypeA

// Test fields that aren't supported: list of non-strings or structured objects
func TestUnitApplyFilters_unsupportedTypes(t *testing.T) {
	items := []map[string]interface{}{
		{
			"nums": []int{1, 2, 3},
		},
		{
			"nums": []int{3, 4, 5},
		},
		{
			"nums": []int{5, 6, 7},
		},
	}

	testSchema := map[string]*schema.Schema{
		"nums": {
			Type: schema.TypeList,
			Elem: schema.TypeInt,
		},
	}

	filters := &schema.Set{
		F: func(v interface{}) int {
			return schema.HashString(v.(map[string]interface{})["name"])
		},
	}

	intArrayFilter := map[string]interface{}{
		"name":   "nums",
		"values": []interface{}{"1", "3", "5"},
	}
	filters.Add(intArrayFilter)

	res := ApplyFilters(filters, items, testSchema)
	if len(res) != 0 {
		t.Errorf("Expected 0 result, got %d", len(res))
	}
}

func TestUnitApplyFilters_booleanTypes(t *testing.T) {
	items := []map[string]interface{}{
		{
			"enabled": true,
		},
		{
			"enabled": "true",
		},
		{
			"enabled": "1",
		},
		{
			"enabled": false,
		},
		{
			"enabled": "false",
		},
		{
			"enabled": "0",
		},
	}

	testSchema := map[string]*schema.Schema{
		"enabled": {
			Type: schema.TypeBool,
		},
	}

	filters := &schema.Set{
		F: func(v interface{}) int {
			return schema.HashString(v.(map[string]interface{})["name"])
		},
	}

	truthyBooleanFilter := map[string]interface{}{
		"name":   "enabled",
		"values": []interface{}{"true", "1"}, // while we can pass an actual boolean true here in the test, terraform
		// doesnt, so keep coercion logic simple in filters.go
	}
	filters.Add(truthyBooleanFilter)

	res := ApplyFilters(filters, items, testSchema)

	for _, i := range res {
		switch enabled := i["enabled"].(type) {
		case bool:
			if !enabled {
				t.Errorf("Expected a truthy value, got %t", enabled)
			}
		case string:
			enabledBool, _ := strconv.ParseBool(enabled)
			if !enabledBool {
				t.Errorf("Expected a truthy value, got %s", enabled)
			}
		}
	}

	if len(res) != 3 {
		t.Errorf("Expected 3 results, got %d", len(res))
	}
	filters.Remove(truthyBooleanFilter)

	falsyBooleanFilter := map[string]interface{}{
		"name":   "enabled",
		"values": []interface{}{"false", "0"},
	}
	filters.Add(falsyBooleanFilter)

	res = ApplyFilters(filters, items, testSchema)

	for _, i := range res {
		switch enabled := i["enabled"].(type) {
		case bool:
			if enabled {
				t.Errorf("Expected a falsy value, got %t", enabled)
			}
		case string:
			enabledBool, _ := strconv.ParseBool(enabled)
			if enabledBool {
				t.Errorf("Expected a falsy value, got %s", enabled)
			}
		}
	}

	if len(res) != 3 {
		t.Errorf("Expected 3 results, got %d", len(res))
	}
	filters.Remove(falsyBooleanFilter)
}

func TestUnitApplyFilters_numberTypes(t *testing.T) {
	items := []map[string]interface{}{
		{
			"integer": 1,
			"float":   1.1,
		},
		{
			"integer": 2,
			"float":   2.2,
		},
		{
			"integer": 3,
			"float":   3.3,
		},
	}

	testSchema := map[string]*schema.Schema{
		"integer": {
			Type: schema.TypeInt,
		},
		"float": {
			Type: schema.TypeFloat,
		},
	}

	filters := &schema.Set{
		F: func(v interface{}) int {
			return schema.HashString(v.(map[string]interface{})["name"])
		},
	}

	// int filter with single target value
	intFilter := map[string]interface{}{
		"name":   "integer",
		"values": []interface{}{"2"},
	}
	filters.Add(intFilter)

	res := ApplyFilters(filters, items, testSchema)
	if len(res) != 1 {
		t.Errorf("Expected 1 result, got %d", len(res))
	}
	filters.Remove(intFilter)

	// test filter with multiple target value
	intsFilter := map[string]interface{}{
		"name":   "integer",
		"values": []interface{}{"1", "3"},
	}
	filters.Add(intsFilter)

	res = ApplyFilters(filters, items, testSchema)
	if len(res) != 2 {
		t.Errorf("Expected 2 results, got %d", len(res))
	}
	filters.Remove(intsFilter)

	// test float filter
	floatFilter := map[string]interface{}{
		"name":   "float",
		"values": []interface{}{"1.1", "3.3"},
	}
	filters.Add(floatFilter)

	res = ApplyFilters(filters, items, testSchema)
	if len(res) != 2 {
		t.Errorf("Expected 2 results, got %d", len(res))
	}
	filters.Remove(floatFilter)
}

func TestUnitApplyFilters_multiProperty(t *testing.T) {
	items := []map[string]interface{}{
		{
			"letter": "a",
			"number": "1",
			"symbol": "!",
		},
		{
			"letter": "b",
			"number": "2",
			"symbol": "@",
		},
		{
			"letter": "c",
			"number": "3",
			"symbol": "#",
		},
		{
			"letter": "d",
			"number": "4",
			"symbol": "$",
		},
	}

	testSchema := map[string]*schema.Schema{
		"letter": {
			Type: schema.TypeString,
		},
		"number": {
			Type: schema.TypeInt,
		},
		"symbol": {
			Type: schema.TypeString,
		},
	}

	filters := &schema.Set{
		F: func(v interface{}) int {
			return schema.HashString(v.(map[string]interface{})["name"])
		},
	}
	filters.Add(map[string]interface{}{
		"name":   "letter",
		"values": []interface{}{"a", "b", "c"},
	})
	filters.Add(map[string]interface{}{
		"name":   "number",
		"values": []interface{}{"2", "3", "4"},
	})
	filters.Add(map[string]interface{}{
		"name":   "symbol",
		"values": []interface{}{"#", "$"},
	})

	res := ApplyFilters(filters, items, testSchema)
	if len(res) != 1 {
		t.Errorf("Expected 1 result, got %d", len(res))
	}
}

// Test to validate that the Apply filters do not impact the original item order
func TestUnitApplyFilters_ElementOrder(t *testing.T) {
	items := []map[string]interface{}{
		{"letter": "a"},
		{"letter": "b"},
		{"letter": "c"},
		{"letter": "d"},
	}

	testSchema := map[string]*schema.Schema{
		"letter": {
			Type: schema.TypeString,
		},
	}

	filters := &schema.Set{F: func(interface{}) int { return 1 }}
	filters.Add(map[string]interface{}{
		"name":   "letter",
		"values": []interface{}{"a", "d"},
	})

	res := ApplyFilters(filters, items, testSchema)
	if len(res) != 2 {
		t.Errorf("Expected 2 result, got %d", len(res))
	}
	if res[0]["letter"] != "a" || res[1]["letter"] != "d" {
		t.Errorf("Expected sort order not retained, got %s %s", res[0]["letter"], res[1]["letter"])
	}

}

func TestUnitGetValue_EmptyMap(t *testing.T) {
	item := map[string]interface{}{}

	_, singleLevelGetOk := getValueFromPath(item, []string{"path"})
	_, multiLevelGetOk := getValueFromPath(item, []string{"path", "to", "target"})

	if singleLevelGetOk || multiLevelGetOk {
		t.Error("Expected non OK result")
	}
}

func TestUnitGetValue_MultiLevelMap(t *testing.T) {
	item := map[string]interface{}{
		"level1": map[string]interface{}{
			"level2": map[string]interface{}{
				"level3": "value",
			},
		},
	}

	singleLevelGet, singleLevelGetOk := getValueFromPath(item, []string{"level1"})
	multiLevelGet, multiLevelGetOk := getValueFromPath(item, []string{"level1", "level2", "level3"})

	if !singleLevelGetOk || !multiLevelGetOk {
		t.Errorf("Expected OK result for topLevel %v multi level %v", singleLevelGetOk, multiLevelGetOk)
	}

	if multiLevelGet != "value" {
		t.Errorf("Expected = value, Got = %s", multiLevelGet)
	}

	if len(singleLevelGet.(map[string]interface{})) != 1 {
		t.Error("Expected size of map is 1")
	}
}

func TestUnitNestedMap(t *testing.T) {
	item := map[string]interface{}{
		"level1": map[string]interface{}{
			"level2": map[string]interface{}{
				"level3":   "value",
				"level3_1": []string{"A", "B", "C"},
				"level3_2": []int{2, 3, 4},
				"level3_3": []float64{2.1, 3.1, 4.1},
				"level3_4": []interface{}{2, 3, 4},
			},
		},
	}
	services := genericMapToJsonMap(item)

	if len(services) != 1 {
		t.Errorf("unexpected number of values returned in map")
	}
}

// Helper to marshal JSON objects from service into strings that can be stored in state.
// This limitation exists because Terraform doesn't support maps of nested objects and so we use JSON strings representation
// as a workaround.
func genericMapToJsonMap(genericMap map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}

	for key, value := range genericMap {
		switch v := value.(type) {
		case string:
			result[key] = v
		default:
			bytes, err := json.Marshal(v)
			if err != nil {
				continue
			}
			result[key] = string(bytes)
		}
	}

	return result
}

func Test_convertToObjectMap(t *testing.T) {

	m := make(map[string]string)
	m["key1"] = "value1"

	m1 := make(map[string]string)
	m1["key1"] = ""

	type args struct {
		stringTostring map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{name: "stringNoEmpty", args: args{m}, want: map[string]interface{}{"key1": "value1"}},
		{name: "stringEmpty", args: args{m1}, want: map[string]interface{}{"key1": ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertToObjectMap(tt.args.stringTostring); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToObjectMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkAndConvertMap(t *testing.T) {

	m := make(map[string]string)
	m["key1"] = "value1"

	m1 := make(map[string]interface{})
	m1["key1"] = "value1"

	m2 := make(map[string]int)
	m2["key1"] = 1

	type args struct {
		element interface{}
	}
	tests := []struct {
		name  string
		args  args
		want  map[string]interface{}
		want1 bool
	}{
		{name: "string", args: args{m}, want: map[string]interface{}{"key1": "value1"}, want1: true},
		{name: "interface", args: args{m1}, want: map[string]interface{}{"key1": "value1"}, want1: true},
		{name: "int", args: args{m2}, want: nil, want1: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := checkAndConvertMap(tt.args.element)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkAndConvertMap() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("checkAndConvertMap() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
