package statuspage

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sp "github.com/sbecker59/statuspage-api-client-go/api/v1/statuspage"
)

var (
	testAccProviders map[string]*schema.Provider
	testAccProvider  *schema.Provider
	pageID           string
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"statuspage": testAccProvider,
	}
	pageID = os.Getenv("STATUSPAGE_PAGE_ID")
}

func isDebug() bool {
	return os.Getenv("DEBUG") == "true"
}

func isAPIKeySet() bool {
	if os.Getenv("SP_API_KEY") != "" {
		return true
	}
	if os.Getenv("STATUSPAGE_API_KEY") != "" {
		return true
	}
	return false
}

func isPageIdSet() bool {
	if os.Getenv("STATUSPAGE_PAGE_ID") != "" {
		return true
	}
	return false
}

func isDatadogApiKeySet() bool {
	if os.Getenv("DD_API_KEY") != "" {
		return true
	}
	return false
}

func isDatadogAppKeySet() bool {
	if os.Getenv("DD_APP_KEY") != "" {
		return true
	}
	return false
}

// testAccPreCheck validates the necessary test API keys exist
// in the testing environment
func testAccPreCheck(t *testing.T) {
	if !isAPIKeySet() {
		t.Fatal("STATUSPAGE_API_KEY or SP_API_KEY must be set for acceptance tests")
	}
	if !isPageIdSet() {
		t.Fatal("STATUSPAGE_PAGE_ID must be set for acceptance tests")
	}
	if !isDatadogApiKeySet() {
		t.Fatal("DD_API_KEY must be set for acceptance tests")
	}
	if !isDatadogAppKeySet() {
		t.Fatal("DD_APP_KEY must be set for acceptance tests")
	}
}

func buildStatuspageClientV1(httpClient *http.Client) *sp.APIClient {
	configV1 := sp.NewConfiguration()
	configV1.Debug = isDebug()
	configV1.HTTPClient = httpClient
	configV1.UserAgent = getUserAgent(configV1.UserAgent)
	return sp.NewAPIClient(configV1)
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestUnittranslateClientError_MsgEmptyAndGenericError(t *testing.T) {
	if err := translateClientError(errors.New(""), ""); err != nil {
		if !strings.Contains(err.Error(), "an error occurred") {
			t.Error("TestUnittranslateClientError_MsgEmptyAndGenericError")
		}
	}
}

func TestUnittranslateClientError_MsgAndGenericError(t *testing.T) {
	if err := translateClientError(errors.New(""), ""); err != nil {
		if !strings.Contains(err.Error(), "an error occurred") {
			t.Error("TestUnittranslateClientError_MsgAndGenericError")
		}
	}
}

type GenericOpenAPIError struct {
	body  []byte
	error string
	model interface{}
}

func (e GenericOpenAPIError) Error() string {
	return e.error
}

func TestUnittranslateClientError_MsgAndGenericOpenAPIError(t *testing.T) {

	genericOpenAPIError := GenericOpenAPIError{
		body:  make([]byte, 128),
		error: errors.New("GenericOpenAPIError").Error(),
		model: nil,
	}

	if err := translateClientError(genericOpenAPIError, ""); err != nil {
		if !strings.Contains(err.Error(), "an error occurred: GenericOpenAPIError") {
			t.Error("TestUnittranslateClientError_MsgAndGenericOpenAPIError")
		}
	}
}

func TestUnittranslateClientError_MsgAndUrlError(t *testing.T) {

	urlError := &url.Error{
		Op:  "Op",
		URL: "Url",
		Err: errors.New("Err"),
	}

	if err := translateClientError(urlError, ""); err != nil {
		if !strings.Contains(err.Error(), "an error occurred: (url.Error): Op") && !strings.Contains(err.Error(), "Url") && !strings.Contains(err.Error(), "Err") {
			t.Error("TestUnittranslateClientError_MsgAndUrlError")
		}
	}
}
