package statuspage

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	sp "github.com/sbecker59/statuspage-api-client-go/api/v1/statuspage"
)

func resourceMetricProviderRead(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	name := d.Get("name").(string)
	log.Printf("[INFO] Reading OpsGenie metric providers '%s'", name)

	metricProviders, _, err := statuspageClientV1.MetricProvidersApi.GetPagesPageIdMetricsProvidersMetricsProviderId(authV1, d.Get("page_id").(string), d.Id())
	if err != nil {
		return translateClientError(err, "failed to get metric providers using Status Page API")
	}

	d.Set("type", metricProviders.Type)
	d.Set("metric_base_uri", metricProviders.MetricBaseUri)

	return nil
}

func resourceMetricProviderCreate(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	email := d.Get("email").(string)
	passowrd := d.Get("passowrd").(string)
	apiKey := d.Get("api_key").(string)
	apiToken := d.Get("api_token").(string)
	applicationKey := d.Get("application_key").(string)
	t := d.Get("type").(string)
	metricBaseUri := d.Get("metric_base_uri").(string)

	p := sp.PostPagesPageIdMetricsProviders{
		MetricsProvider: &sp.PostPagesPageIdMetricsProvidersMetricsProvider{
			Email:          email,
			Password:       passowrd,
			ApiKey:         apiKey,
			ApiToken:       apiToken,
			ApplicationKey: applicationKey,
			Type:           t,
			MetricBaseUri:  metricBaseUri,
		},
	}

	log.Printf("[INFO] Creating Status Page metric providers '%s'", t)
	result, _, err := statuspageClientV1.MetricProvidersApi.PostPagesPageIdMetricsProviders(authV1, d.Get("page_id").(string), p)

	if err != nil {
		return translateClientError(err, "failed to create metric providers using Status Page API")
	}

	d.SetId(result.Id)

	return resourceMetricProviderRead(d, m)

}

func resourceMetricProviderUpdate(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	t := d.Get("type").(string)
	metricBaseUri := d.Get("metric_base_uri").(string)

	p := sp.PatchPagesPageIdMetricsProviders{
		MetricsProvider: &sp.PatchPagesPageIdMetricsProvidersMetricsProvider{
			Type:          t,
			MetricBaseUri: metricBaseUri,
		},
	}

	log.Printf("[INFO] Update Status Page metric providers '%s'", t)
	_, _, err := statuspageClientV1.MetricProvidersApi.PatchPagesPageIdMetricsProvidersMetricsProviderId(authV1, d.Get("page_id").(string), d.Id(), p)

	if err != nil {
		return translateClientError(err, "failed to update metric providers using Status Page API")
	}

	return nil
}

func resourceMetricProviderDelete(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	_, _, err := statuspageClientV1.MetricProvidersApi.DeletePagesPageIdMetricsProvidersMetricsProviderId(authV1, d.Get("page_id").(string), d.Id())

	if err != nil {
		return translateClientError(err, "failed to delete metric providers using Status Page API")
	}

	return nil
}

func resourceMetricProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetricProviderCreate,
		Read:   resourceMetricProviderRead,
		Update: resourceMetricProviderUpdate,
		Delete: resourceMetricProviderDelete,
		Schema: map[string]*schema.Schema{
			"page_id": {
				Type:        schema.TypeString,
				Description: "The ID of the page this metric provider belongs to",
				Required:    true,
			},
			"email": {
				Type:        schema.TypeString,
				Description: "Required by the Librato and Pingdom type metrics providers",
				Optional:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "Required by the Pingdom-type metrics provider",
				Optional:    true,
				Sensitive:   true,
			},
			"api_key": {
				Type:        schema.TypeString,
				Description: "Required by the Datadog and NewRelic type metrics providers",
				Optional:    true,
				Sensitive:   true,
			},
			"api_token": {
				Type:        schema.TypeString,
				Description: "Required by the Librato, Pingdom and Datadog type metrics providers",
				Optional:    true,
				Sensitive:   true,
			},
			"application_key": {
				Type:        schema.TypeString,
				Description: "Required by the Pingdom-type metrics provider",
				Optional:    true,
				Sensitive:   true,
			},
			"type": {
				Type:         schema.TypeString,
				Description:  "One of 'Pingdom', 'NewRelic', 'Librato', 'Datadog', or 'Self'",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Pingdom", "NewRelic", "Librato", "Datadog", "Self"}, false),
			},
			"metric_base_uri": {
				Type:        schema.TypeString,
				Description: "Required by the NewRelic-type metrics provider",
				Optional:    true,
			},
		},
	}
}
