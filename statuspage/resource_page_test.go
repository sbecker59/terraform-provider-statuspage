package statuspage

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	sp "github.com/sbecker59/statuspage-api-client-go/api/v1/statuspage"
)

// resourcePageTestData builds a *schema.ResourceData for resourcePage(),
// pre-populated with attrs. It is enough for d.Get/d.Set/d.SetId, which is
// everything Read/Update/Create/Delete need — but NOT for d.GetRawConfig():
// schema.TestResourceDataRaw drives the legacy Diff() path, which never
// populates InstanceDiff.RawConfig, so GetRawConfig() reads back as null
// here regardless of attrs. Tests that need a realistic raw config (i.e.
// buildPagePatchBody's tests) use pageRawConfig below instead.
func resourcePageTestData(t *testing.T, attrs map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := resourcePage()
	d := schema.TestResourceDataRaw(t, r.Schema, attrs)
	return d
}

// pageRawConfig builds the cty.Value that a real `terraform apply` would hand
// to d.GetRawConfig(): every schema attribute present, null unless attrs
// names it explicitly. Built straight from the resource's own implied type,
// so it can't drift from the schema.
func pageRawConfig(t *testing.T, attrs map[string]interface{}) cty.Value {
	t.Helper()
	implied := resourcePage().CoreConfigSchema().ImpliedType()

	vals := make(map[string]cty.Value, len(implied.AttributeTypes()))
	for name, ty := range implied.AttributeTypes() {
		v, ok := attrs[name]
		if !ok {
			vals[name] = cty.NullVal(ty)
			continue
		}
		switch x := v.(type) {
		case string:
			vals[name] = cty.StringVal(x)
		case bool:
			vals[name] = cty.BoolVal(x)
		default:
			t.Fatalf("pageRawConfig: unsupported value type %T for %q", v, name)
		}
	}
	return cty.ObjectVal(vals)
}

func TestUnitBuildPagePatchBody_OnlyConfigSetFieldsAreSent(t *testing.T) {
	attrs := map[string]interface{}{
		"page_id": "abc123def456",
		"name":    "prisma CLOUD status page",
		"domain":  "status.example.com",
		// every other attribute intentionally absent from attrs
	}
	d := resourcePageTestData(t, attrs)

	body := buildPagePatchBody(d, pageRawConfig(t, attrs))

	if got, want := body.GetName(), "prisma CLOUD status page"; got != want {
		t.Errorf("Name = %q, want %q", got, want)
	}
	if got, want := body.GetDomain(), "status.example.com"; got != want {
		t.Errorf("Domain = %q, want %q", got, want)
	}
	if body.Subdomain != nil {
		t.Errorf("Subdomain = %v, want nil (not set in config)", *body.Subdomain)
	}
	if body.TimeZone != nil {
		t.Errorf("TimeZone = %v, want nil (not set in config)", *body.TimeZone)
	}
	if body.HiddenFromSearch != nil {
		t.Errorf("HiddenFromSearch = %v, want nil (not set in config)", *body.HiddenFromSearch)
	}
}

func TestUnitBuildPagePatchBody_ExplicitFalseBoolIsSent(t *testing.T) {
	// The case d.GetOkExists cannot handle: an explicit `false` must appear in
	// the body exactly like an explicit `true` would, and must be
	// distinguishable from a bool the practitioner never mentioned.
	attrs := map[string]interface{}{
		"page_id":                   "abc123def456",
		"allow_webhook_subscribers": false,
		"allow_sms_subscribers":     true,
		// allow_rss_atom_feeds intentionally absent
	}
	d := resourcePageTestData(t, attrs)

	body := buildPagePatchBody(d, pageRawConfig(t, attrs))

	if body.AllowWebhookSubscribers == nil {
		t.Fatal("AllowWebhookSubscribers = nil, want a set false (it was written in config)")
	}
	if *body.AllowWebhookSubscribers != false {
		t.Errorf("AllowWebhookSubscribers = %v, want false", *body.AllowWebhookSubscribers)
	}
	if body.AllowSmsSubscribers == nil || *body.AllowSmsSubscribers != true {
		t.Errorf("AllowSmsSubscribers = %v, want true", body.AllowSmsSubscribers)
	}
	if body.AllowRssAtomFeeds != nil {
		t.Errorf("AllowRssAtomFeeds = %v, want nil (not set in config)", *body.AllowRssAtomFeeds)
	}
}

func TestUnitBuildPagePatchBody_EmptyConfigProducesEmptyBody(t *testing.T) {
	attrs := map[string]interface{}{
		"page_id": "abc123def456",
	}
	d := resourcePageTestData(t, attrs)

	body := buildPagePatchBody(d, pageRawConfig(t, attrs))

	empty := sp.PatchPagesPage{}
	if body != empty {
		t.Errorf("buildPagePatchBody() = %+v, want an all-nil PatchPagesPage", body)
	}
}

func TestUnitSetPageState_MapsResponseNotConfig(t *testing.T) {
	// The state must come from the API response even when it disagrees with
	// whatever was last in config/state — that's what makes drift visible.
	d := resourcePageTestData(t, map[string]interface{}{
		"page_id": "abc123def456",
		"name":    "stale local value",
	})

	var page sp.Page
	page.SetId("abc123def456")
	page.SetName("prisma CLOUD status page")
	page.SetDomain("status.prismacloud.com")
	page.SetTimeZone("Berlin")
	page.SetBranding("premium")
	page.SetHiddenFromSearch(true)
	page.SetAllowWebhookSubscribers(false)
	page.SetSupportUrl("https://support.example.com")
	page.SetActivityScore(185)
	page.SetSubdomain("prismacloud1")

	setPageState(d, &page)

	if got, want := d.Id(), "abc123def456"; got != want {
		t.Errorf("Id() = %q, want %q", got, want)
	}
	if got, want := d.Get("name").(string), "prisma CLOUD status page"; got != want {
		t.Errorf("name = %q, want %q (response, not stale config)", got, want)
	}
	if got, want := d.Get("domain").(string), "status.prismacloud.com"; got != want {
		t.Errorf("domain = %q, want %q", got, want)
	}
	if got, want := d.Get("time_zone").(string), "Berlin"; got != want {
		t.Errorf("time_zone = %q, want %q", got, want)
	}
	if got, want := d.Get("hidden_from_search").(bool), true; got != want {
		t.Errorf("hidden_from_search = %v, want %v", got, want)
	}
	if got, want := d.Get("allow_webhook_subscribers").(bool), false; got != want {
		t.Errorf("allow_webhook_subscribers = %v, want %v", got, want)
	}
	// Computed-only field: never sent, but must round-trip on Read.
	if got, want := d.Get("support_url").(string), "https://support.example.com"; got != want {
		t.Errorf("support_url = %q, want %q", got, want)
	}
	if got, want := d.Get("activity_score").(float64), 185.0; got != want {
		t.Errorf("activity_score = %v, want %v", got, want)
	}
}

func TestUnitResourcePageCreate_AlwaysErrorsWithImportHint(t *testing.T) {
	d := resourcePageTestData(t, map[string]interface{}{
		"page_id": "abc123def456",
	})

	err := resourcePageCreate(d, nil)
	if err == nil {
		t.Fatal("resourcePageCreate() = nil error, want an error (pages cannot be created)")
	}
	for _, want := range []string{"cannot create", "abc123def456", "import"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not mention %q", err.Error(), want)
		}
	}
}

func TestUnitResourcePageDeleteContext_NoRequestWarningOnly(t *testing.T) {
	d := resourcePageTestData(t, map[string]interface{}{
		"page_id": "abc123def456",
	})
	d.SetId("abc123def456")

	// No ProviderConfiguration/HTTP client is wired into m at all — if the
	// implementation ever tried to make a request, this test would panic on
	// the nil-pointer dereference rather than silently pass.
	diags := resourcePageDeleteContext(context.Background(), d, nil)

	if diags.HasError() {
		t.Fatalf("DeleteContext returned an error-severity diagnostic: %v", diags)
	}
	if len(diags) != 1 {
		t.Fatalf("DeleteContext returned %d diagnostics, want exactly 1 warning: %v", len(diags), diags)
	}
	if diags[0].Severity != diag.Warning {
		t.Errorf("diagnostic severity = %v, want Warning", diags[0].Severity)
	}
	if !strings.Contains(diags[0].Detail, "cannot be deleted") {
		t.Errorf("diagnostic detail %q does not explain the page still exists", diags[0].Detail)
	}
}

// --- Read: 404 vs 401/403, against a local httptest server -----------------

func newTestPagesClient(t *testing.T, handler http.HandlerFunc) (*sp.APIClient, context.Context) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	cfg := sp.NewConfiguration()
	cfg.Servers = sp.ServerConfigurations{{URL: srv.URL}}

	authV1 := context.WithValue(context.Background(), sp.ContextAPIKeys, map[string]sp.APIKey{
		"api_key": {Key: "test-key", Prefix: "oauth"},
	})
	return sp.NewAPIClient(cfg), authV1
}

func TestUnitResourcePageRead_404RemovesFromState(t *testing.T) {
	client, authV1 := newTestPagesClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"not found"}`))
	})

	d := resourcePageTestData(t, map[string]interface{}{"page_id": "abc123def456"})
	d.SetId("abc123def456")
	m := &ProviderConfiguration{StatuspageClientV1: client, AuthV1: authV1}

	if err := resourcePageRead(d, m); err != nil {
		t.Fatalf("resourcePageRead() error = %v, want nil (404 is handled, not surfaced)", err)
	}
	if d.Id() != "" {
		t.Errorf("Id() = %q after a 404, want empty (removed from state)", d.Id())
	}
}

func TestUnitResourcePageRead_401KeepsState(t *testing.T) {
	client, authV1 := newTestPagesClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"Could not authenticate"}`))
	})

	d := resourcePageTestData(t, map[string]interface{}{"page_id": "abc123def456"})
	d.SetId("abc123def456")
	m := &ProviderConfiguration{StatuspageClientV1: client, AuthV1: authV1}

	err := resourcePageRead(d, m)
	if err == nil {
		t.Fatal("resourcePageRead() error = nil, want an error on 401")
	}
	if d.Id() != "abc123def456" {
		t.Errorf("Id() = %q after a 401, want it unchanged (the page was not proven gone)", d.Id())
	}
}

func TestUnitResourcePageRead_MapsSuccessfulResponse(t *testing.T) {
	client, authV1 := newTestPagesClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":"abc123def456","name":"live name","time_zone":"Berlin"}`))
	})

	d := resourcePageTestData(t, map[string]interface{}{"page_id": "abc123def456"})
	d.SetId("abc123def456")
	m := &ProviderConfiguration{StatuspageClientV1: client, AuthV1: authV1}

	if err := resourcePageRead(d, m); err != nil {
		t.Fatalf("resourcePageRead() error = %v, want nil", err)
	}
	if got, want := d.Get("name").(string), "live name"; got != want {
		t.Errorf("name = %q, want %q", got, want)
	}
	if got, want := d.Get("time_zone").(string), "Berlin"; got != want {
		t.Errorf("time_zone = %q, want %q", got, want)
	}
}

func TestUnitResourcePageUpdate_SendsBodyAndMapsResponse(t *testing.T) {
	// This exercises resourcePageUpdate end to end: it hits PATCH and maps
	// the response back into state. It does NOT assert which fields end up
	// in the body — that's TestUnitBuildPagePatchBody_*'s job, and it needs a
	// realistic raw config (see pageRawConfig) that d.GetRawConfig() cannot
	// supply when d comes from schema.TestResourceDataRaw.
	var gotMethod, gotBody string
	client, authV1 := newTestPagesClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		buf := make([]byte, r.ContentLength)
		r.Body.Read(buf)
		gotBody = string(buf)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":"abc123def456","name":"prisma CLOUD status page"}`))
	})

	d := resourcePageTestData(t, map[string]interface{}{
		"page_id": "abc123def456",
		"name":    "prisma CLOUD status page",
	})
	d.SetId("abc123def456")
	m := &ProviderConfiguration{StatuspageClientV1: client, AuthV1: authV1}

	if err := resourcePageUpdate(d, m); err != nil {
		t.Fatalf("resourcePageUpdate() error = %v, want nil", err)
	}
	if gotMethod != http.MethodPatch {
		t.Errorf("method = %s, want PATCH", gotMethod)
	}
	if !strings.Contains(gotBody, `"page"`) {
		t.Errorf("PATCH body = %s, want it namespaced under \"page\"", gotBody)
	}
	if got, want := d.Get("name").(string), "prisma CLOUD status page"; got != want {
		t.Errorf("name after update = %q, want %q (mapped from response)", got, want)
	}
}

// --- Acceptance test -------------------------------------------------------
//
// WARNING: TestAccStatuspagePage_basic PATCHes a real Statuspage page. Point
// STATUSPAGE_TEST_PAGE_ID at a disposable test page — never at a production
// page. Statuspage pages cannot be created via the API, so a disposable page
// has to already exist (and is billable); that is exactly why this test is
// opt-in and, per docs/TASK_statuspage_page_resource.md, not part of the
// resource's acceptance criteria. It has not been run against a real page as
// part of this change — no disposable test page has been provisioned yet.
//
// Skips unless TF_ACC=1, STATUSPAGE_TEST_PAGE_ID, and an API key
// (STATUSPAGE_API_KEY or SP_API_KEY) are all set.

func isTestPageIDSet() bool {
	return os.Getenv("STATUSPAGE_TEST_PAGE_ID") != ""
}

func testAccPreCheckPage(t *testing.T) {
	if !isAPIKeySet() {
		t.Fatal("STATUSPAGE_API_KEY or SP_API_KEY must be set for TestAccStatuspagePage_basic")
	}
	if !isTestPageIDSet() {
		t.Fatal("STATUSPAGE_TEST_PAGE_ID must be set for TestAccStatuspagePage_basic — point it at a " +
			"disposable page, never at production; pages cannot be created via the API")
	}
}

func TestAccStatuspagePage_basic(t *testing.T) {
	testPageID := os.Getenv("STATUSPAGE_TEST_PAGE_ID")

	// Snapshot the live name before touching anything, and restore it via
	// t.Cleanup rather than only in a final test step — a t.Cleanup runs even
	// if an earlier step fails or the test is aborted, so the page can't be
	// left mutated by a broken run the way a final-step-only restore could.
	var originalName string
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckPage(t)

			// PreCheck runs before resource.Test configures the provider (that
			// happens once the first Config step plans), so testAccProvider.Meta()
			// is not populated yet here — configure a throwaway client directly
			// instead, purely to snapshot the pre-test state.
			apiKey := os.Getenv("STATUSPAGE_API_KEY")
			if apiKey == "" {
				apiKey = os.Getenv("SP_API_KEY")
			}
			providerData := schema.TestResourceDataRaw(t, testAccProvider.Schema, map[string]interface{}{
				"api_key": apiKey,
			})
			raw, err := providerConfigure(providerData)
			if err != nil {
				t.Fatalf("failed to configure provider for pre-test snapshot: %s", err)
			}
			pc := raw.(*ProviderConfiguration)
			page, _, err := pc.StatuspageClientV1.PagesAPI.GetPagesPageId(pc.AuthV1, testPageID).Execute()
			if err != nil {
				t.Fatalf("failed to snapshot original page state: %s", err)
			}
			originalName = page.GetName()

			t.Cleanup(func() {
				var restore sp.PatchPagesPage
				restore.SetName(originalName)
				if _, _, err := pc.StatuspageClientV1.PagesAPI.PatchPagesPageId(pc.AuthV1, testPageID).
					PatchPages(sp.PatchPages{Page: &restore}).Execute(); err != nil {
					t.Errorf("failed to restore original page name %q: %s", originalName, err)
				}
			})
		},
		Providers:    testAccProviders,
		CheckDestroy: func(s *terraform.State) error { return nil }, // Delete is a no-op by design.
		Steps: []resource.TestStep{
			{
				// Import first: there is no Create, so the only way in is import.
				// ImportStateVerify can't be used here: it diffs the imported state
				// against a prior *applied* state, but nothing has ever been applied
				// (Create always errors by design) — with nothing to compare against,
				// the SDK reports "resource with ID ... not found" rather than a real
				// mismatch. ImportStatePersist carries the imported state forward so
				// the next step's Update has a resource to operate on.
				ResourceName:       "statuspage_page.this",
				ImportState:        true,
				ImportStateId:      testPageID,
				ImportStatePersist: true,
				Config:             testAccStatuspagePageConfig(testPageID, originalName),
			},
			{
				Config: testAccStatuspagePageConfig(testPageID, "tf-testacc-page-renamed"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_page.this", "name", "tf-testacc-page-renamed"),
					testAccCheckStatuspagePageName(testPageID, "tf-testacc-page-renamed"),
				),
			},
		},
	})
}

func testAccStatuspagePageConfig(pageID, name string) string {
	return fmt.Sprintf(`
resource "statuspage_page" "this" {
  page_id = %q
  name    = %q
}
`, pageID, name)
}

func testAccCheckStatuspagePageName(pageID, want string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*ProviderConfiguration)
		page, _, err := conn.StatuspageClientV1.PagesAPI.GetPagesPageId(conn.AuthV1, pageID).Execute()
		if err != nil {
			return TranslateClientErrorDiag(err, "error retrieving page for verification")
		}
		if got := page.GetName(); got != want {
			return fmt.Errorf("page name = %q on the API, want %q — PATCH did not stick", got, want)
		}
		return nil
	}
}
