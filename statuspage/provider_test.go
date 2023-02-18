package statuspage

import (
	"errors"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	testAccProviders       map[string]*schema.Provider
	testAccProvider        *schema.Provider
	pageID                 string
	pageName               string
	audienceSpecificPageID string
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"statuspage": testAccProvider,
	}
	pageID = os.Getenv("STATUSPAGE_PAGE_ID")
	pageName = os.Getenv("STATUSPAGE_PAGE_NAME")
	audienceSpecificPageID = os.Getenv("STATUSPAGE_AUDIENCE_SPECIFIC_PAGE_ID")
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
	return os.Getenv("STATUSPAGE_PAGE_ID") != ""
}

func isPageNameSet() bool {
	return os.Getenv("STATUSPAGE_PAGE_NAME") != ""
}

func isAudienceSpecificPageIdSet() bool {
	return os.Getenv("STATUSPAGE_AUDIENCE_SPECIFIC_PAGE_ID") != ""
}

func isDatadogApiKeySet() bool {
	return os.Getenv("DD_API_KEY") != ""
}

func isDatadogAppKeySet() bool {
	return os.Getenv("DD_APP_KEY") != ""
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
	if !isPageNameSet() {
		t.Fatal("STATUSPAGE_PAGE_NAME must be set for acceptance tests")
	}
	if !isAudienceSpecificPageIdSet() {
		t.Fatal("STATUSPAGE_AUDIENCE_SPECIFIC_PAGE_ID must be set for acceptance tests")
	}
	if !isDatadogApiKeySet() {
		t.Fatal("DD_API_KEY must be set for acceptance tests")
	}
	if !isDatadogAppKeySet() {
		t.Fatal("DD_APP_KEY must be set for acceptance tests")
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestUnitTranslateClientErrorDiag_MsgEmptyAndGenericError(t *testing.T) {
	if err := TranslateClientErrorDiag(errors.New(""), ""); err != nil {
		if !strings.Contains(err.Error(), "an error occurred") {
			t.Error("TestUnitTranslateClientErrorDiag_MsgEmptyAndGenericError")
		}
	}
}

func TestUnitTranslateClientErrorDiag_MsgAndGenericError(t *testing.T) {
	if err := TranslateClientErrorDiag(errors.New(""), ""); err != nil {
		if !strings.Contains(err.Error(), "an error occurred") {
			t.Error("TestUnitTranslateClientErrorDiag_MsgAndGenericError")
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

func TestUnitTranslateClientErrorDiag_MsgAndGenericOpenAPIError(t *testing.T) {

	genericOpenAPIError := GenericOpenAPIError{
		body:  make([]byte, 128),
		error: errors.New("GenericOpenAPIError").Error(),
		model: nil,
	}

	if err := TranslateClientErrorDiag(genericOpenAPIError, ""); err != nil {
		if !strings.Contains(err.Error(), "an error occurred: GenericOpenAPIError") {
			t.Error("TestUnitTranslateClientErrorDiag_MsgAndGenericOpenAPIError")
		}
	}
}

func TestUnitTranslateClientErrorDiag_MsgAndUrlError(t *testing.T) {

	urlError := &url.Error{
		Op:  "Op",
		URL: "Url",
		Err: errors.New("Err"),
	}

	if err := TranslateClientErrorDiag(urlError, ""); err != nil {
		if !strings.Contains(err.Error(), "an error occurred: (url.Error): Op") && !strings.Contains(err.Error(), "Url") && !strings.Contains(err.Error(), "Err") {
			t.Error("TestUnitTranslateClientErrorDiag_MsgAndUrlError")
		}
	}
}
