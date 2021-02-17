package statuspage

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	sp "github.com/sbecker59/statuspage-api-client-go/api/v1/statuspage"
)

func resourceMetricProviderRead(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	err := providerConf.Ratelimiter.Wait(authV1) // This is a blocking call. Honors the rate limit
	if err != nil {
		return translateClientError(err, "error Ratelimiter")
	}

	metricProvider, _, err := statuspageClientV1.MetricProvidersApi.GetPagesPageIdMetricsProvidersMetricsProviderId(authV1, d.Get("page_id").(string), d.Id()).Execute()

	if err.Error() != "" {
		return translateClientError(err, "failed to get metric provider using Status Page API")
	}

	if &metricProvider == nil {
		d.SetId("")
		return nil
	}

	d.Set("type", metricProvider.Type)

	return nil

}

func resourceMetricProviderCreate(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	err := providerConf.Ratelimiter.Wait(authV1) // This is a blocking call. Honors the rate limit
	if err != nil {
		return translateClientError(err, "error Ratelimiter")
	}

	email := d.Get("email").(string)
	password := d.Get("password").(string)
	apiKey := d.Get("api_key").(string)
	apiToken := d.Get("api_token").(string)
	applicationKey := d.Get("application_key").(string)
	metricBaseURI := d.Get("metric_base_uri").(string)
	metricType := d.Get("type").(string)

	var metricProvider sp.PostPagesPageIdMetricsProvidersMetricsProvider

	metricProvider.SetApiKey(apiKey)
	metricProvider.SetApiToken(apiToken)
	metricProvider.SetApplicationKey(applicationKey)
	metricProvider.SetEmail(email)
	metricProvider.SetMetricBaseUri(metricBaseURI)
	metricProvider.SetPassword(password)
	metricProvider.SetType(metricType)

	o := *sp.NewPostPagesPageIdMetricsProviders()
	o.SetMetricsProvider(metricProvider)

	resp, _, err := statuspageClientV1.MetricProvidersApi.PostPagesPageIdMetricsProviders(authV1, d.Get("page_id").(string)).PostPagesPageIdMetricsProviders(o).Execute()

	if err.Error() != "" {
		return translateClientError(err, "failed to create metric provider using Status Page API")
	}

	d.SetId(resp.GetId())

	return resourceMetricProviderRead(d, m)
}

func resourceMetricProviderUpdate(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	err := providerConf.Ratelimiter.Wait(authV1) // This is a blocking call. Honors the rate limit
	if err != nil {
		return translateClientError(err, "error Ratelimiter")
	}

	metricBaseURI := d.Get("metric_base_uri").(string)
	metricType := d.Get("type").(string)

	var metricProvider sp.PatchPagesPageIdMetricsProvidersMetricsProvider

	metricProvider.SetMetricBaseUri(metricBaseURI)
	metricProvider.SetType(metricType)

	o := *sp.NewPatchPagesPageIdMetricsProviders()
	o.SetMetricsProvider(metricProvider)

	resp, _, err := statuspageClientV1.MetricProvidersApi.PatchPagesPageIdMetricsProvidersMetricsProviderId(authV1, d.Get("page_id").(string), d.Id()).PatchPagesPageIdMetricsProviders(o).Execute()

	if err.Error() != "" {
		return translateClientError(err, "failed to create metric provider using Status Page API")
	}

	d.SetId(resp.GetId())

	return resourceMetricProviderRead(d, m)
}

func resourceMetricProviderDelete(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	err := providerConf.Ratelimiter.Wait(authV1) // This is a blocking call. Honors the rate limit
	if err != nil {
		return translateClientError(err, "error Ratelimiter")
	}

	_, _, err = statuspageClientV1.MetricProvidersApi.DeletePagesPageIdMetricsProvidersMetricsProviderId(authV1, d.Get("page_id").(string), d.Id()).Execute()

	if err.Error() != "" {
		return translateClientError(err, "failed to delete component using Status Page API")
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
				Description: "the ID of the page this component group belongs to",
				Required:    true,
				ForceNew:    true,
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
				Type:        schema.TypeString,
				Description: "One of 'Pingdom', 'NewRelic', 'Librato', 'Datadog', or 'Self'",
				Required:    true,
				ValidateFunc: validation.StringInSlice(
					[]string{"Pingdom", "NewRelic", "Librato", "Datadog", "Self"},
					false,
				),
			},
			"metric_base_uri": {
				Type:        schema.TypeString,
				Description: "Required by the NewRelic-type metrics provider",
				Optional:    true,
			},
		},
	}
}
