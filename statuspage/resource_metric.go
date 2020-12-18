package statuspage

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	sp "github.com/sbecker59/statuspage-api-client-go/api/v1/statuspage"
)

func resourceMetricCreate(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	name := d.Get("name").(string)
	metricIdentifier := d.Get("metric_identifier").(string)
	transform := d.Get("transform").(string)
	suffix := d.Get("suffix").(string)
	yAxisMin := d.Get("y_axis_min").(int32)
	yAxisMax := d.Get("y_axis_max").(int32)
	yAxisHidden := d.Get("y_axis_hidden").(bool)
	display := d.Get("display").(bool)
	decimalPlaces := d.Get("decimal_places").(int32)
	tooltipDescription := d.Get("tooltip_description").(string)

	p := sp.PostPagesPageIdMetricsProvidersMetricsProviderIdMetrics{
		Metric: &sp.PostPagesPageIdMetricsProvidersMetricsProviderIdMetricsMetric{
			Name:               name,
			MetricIdentifier:   metricIdentifier,
			Transform:          transform,
			Suffix:             suffix,
			YAxisMin:           yAxisMin,
			YAxisMax:           yAxisMax,
			YAxisHidden:        yAxisHidden,
			Display:            display,
			DecimalPlaces:      decimalPlaces,
			TooltipDescription: tooltipDescription,
		},
	}

	log.Printf("[INFO] Creating Status Page metric providers '%s'", t)
	result, _, err := statuspageClientV1.MetricsApi.PostPagesPageIdMetricsProvidersMetricsProviderIdMetrics(authV1, d.Get("page_id").(string), d.Get("metric_provider_id").(string), p)

	if err != nil {
		return translateClientError(err, "failed to create metric providers using Status Page API")
	}

	d.SetId(result.Id)

	return resourceMetricRead(d, m)
}

func resourceMetricRead(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	metric, _, err := statuspageClientV1.MetricsApi.GetPagesPageIdMetricsMetricId(authV1, d.Get("page_id").(string), d.Id())
	if err != nil {
		return translateClientError(err, "failed to get metric using Status Page API")
	}

	d.Set("type", metric.Name)
	d.Set("metric_base_uri", metric.MetricIdentifier)
	d.Set("suffix", metric.Suffix)
	d.Set("y_axis_min", metric.YAxisMin)
	d.Set("y_axis_max", metric.YAxisMax)
	d.Set("y_axis_hidden", metric.YAxisHidden)
	d.Set("display", metric.Display)
	d.Set("decimal_places", metric.DecimalPlaces)
	d.Set("tooltip_description", metric.TooltipDescription)

	return nil
}

func resourceMetricUpdate(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	name := d.Get("name").(string)
	metricIdentifier := d.Get("metric_identifier").(string)
	transform := d.Get("transform").(string)
	suffix := d.Get("suffix").(string)
	yAxisMin := d.Get("y_axis_min").(int32)
	yAxisMax := d.Get("y_axis_max").(int32)
	yAxisHidden := d.Get("y_axis_hidden").(bool)
	display := d.Get("display").(bool)
	decimalPlaces := d.Get("decimal_places").(int32)
	tooltipDescription := d.Get("tooltip_description").(string)

	p := sp.PutPagesPageIdMetrics{
		Metric: &sp.PatchPagesPageIdMetricsMetric{
			Name:               name,
			MetricIdentifier:   metricIdentifier,
			Transform:          transform,
			Suffix:             suffix,
			YAxisMin:           yAxisMin,
			YAxisMax:           yAxisMax,
			YAxisHidden:        yAxisHidden,
			Display:            display,
			DecimalPlaces:      decimalPlaces,
			TooltipDescription: tooltipDescription,
		},
	}

	log.Printf("[INFO] Update Status Page metric providers '%s'", t)
	_, _, err := statuspageClientV1.MetricsApi.PatchPagesPageIdMetricsMetricId(authV1, d.Get("page_id").(string), d.Id(), p)

	if err != nil {
		return translateClientError(err, "failed to update metric providers using Status Page API")
	}

	return nil
}

func resourceMetricDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceMetric() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetricCreate,
		Read:   resourceMetricRead,
		Update: resourceMetricUpdate,
		Delete: resourceMetricDelete,
		Schema: map[string]*schema.Schema{
			"page_id": {
				Type:        schema.TypeString,
				Description: "The ID of the page this metric belongs to",
				Required:    true,
			},
			"metric_provider_id": {
				Type:        schema.TypeString,
				Description: "ID of the metric provider",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Display name for the metric",
				Optional:    true,
			},
			"metric_identifier": {
				Type:        schema.TypeString,
				Description: "The identifier used to look up the metric data from the provider",
				Optional:    true,
			},
			"transform": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The transform to apply to metric before pulling into Statuspage. One of: 'average', 'count', 'max', 'min', 'sum', 'response_time' or 'uptime'",
				ValidateFunc: validation.StringInSlice(
					[]string{"average", "count", "max", "min", "sum", "response_time", "uptime"},
					false,
				),
			},
			"suffix": {
				Type:        schema.TypeString,
				Description: "Suffix to describe the units on the graph",
				Optional:    true,
			},
			"y_axis_min": {
				Type:        schema.TypeFloat,
				Description: "The lower bound of the y axis",
				Optional:    true,
			},
			"y_axis_max": {
				Type:        schema.TypeFloat,
				Description: "The upper bound of the y axis",
				Optional:    true,
			},
			"y_axis_hidden": {
				Type:        schema.TypeBool,
				Description: "Should the values on the y axis be hidden on render",
				Optional:    true,
			},
			"display": {
				Type:        schema.TypeBool,
				Description: "Should the metric be displayed",
				Optional:    true,
			},
			"decimal_places": {
				Type:        schema.TypeInt,
				Description: "How many decimal places to render on the graph",
				Optional:    true,
			},
			"tooltip_description": {
				Type:        schema.TypeString,
				Description: "Tooltip for the metric",
				Optional:    true,
			},
		},
	}
}
