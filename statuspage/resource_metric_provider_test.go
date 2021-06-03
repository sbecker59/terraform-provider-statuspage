package statuspage

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccStatuspageMetricProvider_Basic(t *testing.T) {

	time.Sleep(10 * time.Second)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStatuspageMetricProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetricProviderConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("statuspage_metric_provider.default", "id"),
					resource.TestCheckResourceAttr("statuspage_metric_provider.default", "type", "Datadog"),
				),
			},
			{
				Config: testAccCheckMetricProviderConfigUpdated(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("statuspage_metric_provider.default", "id"),
					resource.TestCheckResourceAttr("statuspage_metric_provider.default", "type", "Datadog"),
				),
			},
		},
	})
}

func testAccCheckMetricProviderConfig() string {
	return fmt.Sprintf(`
	variable "pageid" {
		default = "%s"
	}
	variable "api_key" {
		default = "%s"
	}
	variable "application_key" {
		default = "%s"
	}
	resource "statuspage_metric_provider" "default" {
		page_id = var.pageid
		type = "Datadog"
		api_key = var.api_key
		application_key = var.application_key
		metric_base_uri = "https://app.datadoghq.eu/api/v1"

	}
	`, pageID, os.Getenv("DD_API_KEY"), os.Getenv("DD_APP_KEY"))
}

func testAccCheckMetricProviderConfigUpdated() string {
	return fmt.Sprintf(`
	variable "pageid" {
		default = "%s"
	}
	variable "api_key" {
		default = "%s"
	}
	variable "application_key" {
		default = "%s"
	}
	resource "statuspage_metric_provider" "default" {
		page_id = var.pageid
		type = "Datadog"
		api_key = var.api_key
		application_key = var.application_key
		metric_base_uri = "https://app.datadoghq.eu/api/v1"
	}
	`, pageID, os.Getenv("DD_API_KEY"), os.Getenv("DD_APP_KEY"))
}

func testAccCheckStatuspageMetricProviderDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ProviderConfiguration)
	statuspageClientV1 := conn.StatuspageClientV1
	authV1 := conn.AuthV1

	conn.Ratelimiter.Wait(authV1)

	for _, r := range s.RootModule().Resources {
		_, httpresp, err := statuspageClientV1.MetricProvidersApi.GetPagesPageIdMetricsProvidersMetricsProviderId(authV1, pageID, r.Primary.ID).Execute()
		if err.Error() != "" {
			if httpresp != nil && httpresp.StatusCode == 404 {
				continue
			}
			return translateClientError(err, "error retrieving Metric Provider")
		}
		return fmt.Errorf("Metric Provider still exists")
	}
	return nil
}
