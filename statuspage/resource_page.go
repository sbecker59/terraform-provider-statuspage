package statuspage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	sp "github.com/sbecker59/statuspage-api-client-go/api/v1/statuspage"
)

// The Statuspage API has no POST or DELETE for /pages/{page_id} — only GET and
// PATCH. That shapes every choice below: Create can only error (there is
// nothing to create), Delete can only remove the resource from state (there
// is nothing to delete), and page_id is a required, ForceNew input rather
// than something this resource ever generates. See
// docs/TASK_statuspage_page_resource.md in svc_statuspage/ for the full
// design writeup this resource implements.

func resourcePage() *schema.Resource {
	return &schema.Resource{
		Create: resourcePageCreate,
		Read:   resourcePageRead,
		Update: resourcePageUpdate,
		// DeleteContext (not Delete) solely so it can return a Warning-severity
		// diagnostic instead of an error — the classic Delete signature can only
		// return `error`. A Warning still removes the resource from state (the
		// SDK does that itself once Apply sees no error-severity diagnostic),
		// it just also tells the practitioner the real page was left alone.
		DeleteContext: resourcePageDeleteContext,
		Importer: &schema.ResourceImporter{
			State: resourcePageImport,
		},

		Schema: map[string]*schema.Schema{
			"page_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Statuspage page code (e.g. `abc123def456`, from the manage.statuspage.io URL). Pages cannot be created or deleted through the API, so this must name an existing page; bring it under management with `terraform import` or an `import {}` block.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the page as displayed to visitors.",
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Custom domain (CNAME host) serving the page, e.g. `status.example.com`. " +
					"Changing this makes Statuspage start certificate provisioning and domain " +
					"validation for the new host — a flow the API does not expose, so it is " +
					"invisible to Terraform: `apply` can report success while issuance still " +
					"fails server-side. The CNAME must already point at the Statuspage edge " +
					"before this is set.",
			},
			"subdomain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The `*.statuspage.io` subdomain. Changing it changes the page's canonical URL — treat as risky.",
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "The website linked from the page's logo/image — not the status page's " +
					"own public address (the API has no field for that). It may coincidentally " +
					"equal the public URL without being derived from it.",
			},
			"time_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Timezone for the page, in Rails ActiveSupport form (e.g. `Berlin`, `UTC`) — not IANA (`Europe/Berlin`, `Etc/UTC`).",
			},
			"branding": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"basic", "premium"}, false),
				Description:  "Page template: `basic` or `premium`. The API rejects `premium` on plans that don't offer it.",
			},

			"hidden_from_search": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether the page hides itself from search engines.",
			},
			"viewers_must_be_team_members": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Restrict page visibility to team members only.",
			},
			"allow_page_subscribers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Allow visitors to subscribe to all notifications on the page.",
			},
			"allow_incident_subscribers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Allow visitors to subscribe to notifications for a single incident.",
			},
			"allow_email_subscribers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Allow visitors to receive notifications via email.",
			},
			"allow_sms_subscribers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Allow visitors to receive notifications via SMS.",
			},
			"allow_rss_atom_feeds": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Allow visitors to access incident feeds via RSS/Atom (not functional on audience-specific pages).",
			},
			"allow_webhook_subscribers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Allow visitors to receive notifications via webhooks.",
			},

			"notifications_from_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "From-address for page notification emails.",
			},
			"notifications_email_footer": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Footer appended to notification emails. Accepts Markdown.",
			},

			"css_body_background_color": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CSS color: page body background.",
			},
			"css_font_color": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CSS color: main font.",
			},
			"css_light_font_color": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CSS color: secondary/light font.",
			},
			"css_greens": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CSS color: operational status indicator.",
			},
			"css_yellows": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CSS color: degraded-performance status indicator.",
			},
			"css_oranges": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CSS color: partial-outage status indicator.",
			},
			"css_reds": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CSS color: major-outage status indicator.",
			},
			"css_blues": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CSS color: maintenance status indicator.",
			},
			"css_border_color": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CSS color: borders.",
			},
			"css_graph_color": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CSS color: uptime/metrics graphs.",
			},
			"css_link_color": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CSS color: links.",
			},
			"css_no_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CSS color: graph areas with no data.",
			},

			// Computed-only: readable via GET but absent from PATCH's schema —
			// 17 of Page's 45 fields, per the OpenAPI spec itself (not a client
			// limitation; regenerating the client would not add these). Setting
			// them in config would silently do nothing, so they are not Optional.
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp the page was created (RFC3339).",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp the page was last updated (RFC3339).",
			},
			"page_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Page description. Read-only through this API.",
			},
			"headline": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Page headline. Read-only through this API.",
			},
			"support_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Support URL shown on the page. Read-only through this API — changed in the UI, it moves in state with no way to drive it from config.",
			},
			"twitter_username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Twitter/X username shown on the page. Read-only through this API.",
			},
			"ip_restrictions": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP restriction list. Read-only through this API.",
			},
			"city": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "City shown on the page. Read-only through this API.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State/region shown on the page. Read-only through this API.",
			},
			"country": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Country shown on the page. Read-only through this API.",
			},
			"activity_score": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Statuspage's internal activity score for the page. Read-only.",
			},
		},
	}
}

// resourcePageCreate always fails: there is no POST /pages/{page_id}, so
// nothing can be created. This is deliberate, not a stopgap — it is what
// makes `import {}` / `terraform import` the only way this resource attaches
// to a page, which in turn means a typo in page_id is a hard failure instead
// of a silent PATCH against somebody else's page.
func resourcePageCreate(d *schema.ResourceData, m interface{}) error {
	pageID := d.Get("page_id").(string)
	return fmt.Errorf(
		"statuspage_page cannot create a page: the Statuspage API has no POST /pages endpoint. "+
			"Import the existing page instead: `terraform import statuspage_page.this %s`, or use an `import {}` block.",
		pageID,
	)
}

func resourcePageRead(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	pageID := d.Id()
	log.Printf("[INFO] Reading Status Page page '%s'", pageID)

	page, httpResp, err := statuspageClientV1.PagesAPI.GetPagesPageId(authV1, pageID).Execute()
	if err != nil {
		// A genuine 404 means the page is gone (or the ID was never valid):
		// drop it from state so the next plan proposes re-import rather than an
		// apply that can never succeed (pages cannot be re-created via this API).
		if httpResp != nil && httpResp.StatusCode == 404 {
			log.Printf("[INFO] Statuspage page %q not found, removing from state", pageID)
			d.SetId("")
			return nil
		}
		// 401/403 almost always mean a revoked or under-scoped API key, not a
		// deleted page. Leave state untouched — dropping it here would make the
		// next apply try (and fail) to re-create a page that still exists.
		if httpResp != nil && (httpResp.StatusCode == 401 || httpResp.StatusCode == 403) {
			return TranslateClientErrorDiag(err, fmt.Sprintf(
				"not authorized to read Status Page page %q (revoked or under-scoped API key); state was left unchanged", pageID))
		}
		return TranslateClientErrorDiag(err, "failed to get page using Status Page API")
	}

	setPageState(d, page)
	return nil
}

func resourcePageUpdate(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	pageID := d.Id()
	body := buildPagePatchBody(d, d.GetRawConfig())

	log.Printf("[INFO] Updating Status Page page '%s'", pageID)
	page, _, err := statuspageClientV1.PagesAPI.PatchPagesPageId(authV1, pageID).
		PatchPages(sp.PatchPages{Page: &body}).Execute()
	if err != nil {
		return TranslateClientErrorDiag(err, "failed to update page using Status Page API")
	}

	// The response, not the config, becomes state — that is what surfaces
	// drift on fields nobody manages instead of masking it by echoing config
	// straight back.
	setPageState(d, page)
	return nil
}

// resourcePageDeleteContext makes no API call: there is no DELETE
// /pages/{page_id}. Returning only a Warning-severity diagnostic (never
// Error) still removes the resource from state — the SDK does that on its
// own once Apply() sees no error-severity diagnostic — while telling the
// practitioner the real page was left untouched.
func resourcePageDeleteContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{
		{
			Severity: diag.Warning,
			Summary:  "Statuspage page was only removed from Terraform state",
			Detail: "Statuspage pages cannot be deleted via the API; the resource was removed " +
				"from Terraform state only. The page still exists.",
		},
	}
}

// resourcePageImport accepts a bare page_id as the import ID (there is no
// composite "page_id/id" here — the page IS the resource), sets page_id, and
// defers to Read for everything else.
func resourcePageImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	pageID := d.Id()
	log.Printf("[INFO] Importing Status Page page %q", pageID)

	if err := d.Set("page_id", pageID); err != nil {
		return nil, err
	}
	if err := resourcePageRead(d, m); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

// buildPagePatchBody builds the PATCH body from config, not from the plan.
// With Optional+Computed attributes, Terraform core pre-fills anything absent
// from config into the plan using the prior state value (objchange.ProposedNew)
// — so the plan is never null for a field the practitioner didn't set, and
// reading it here would resend every field that has ever landed in state.
// raw — the caller's d.GetRawConfig() — is the one place that still says
// "unset" honestly; it is a separate parameter (not read from d directly)
// so unit tests can supply a realistic raw config without going through
// Terraform core, which schema.TestResourceDataRaw cannot do — its diff
// never populates RawConfig, so d.GetRawConfig() reads as null in tests.
func buildPagePatchBody(d *schema.ResourceData, raw cty.Value) sp.PatchPagesPage {
	var page sp.PatchPagesPage

	if raw.IsNull() {
		return page
	}

	patchString(&page, raw, d, "name", (*sp.PatchPagesPage).SetName)
	patchString(&page, raw, d, "domain", (*sp.PatchPagesPage).SetDomain)
	patchString(&page, raw, d, "subdomain", (*sp.PatchPagesPage).SetSubdomain)
	patchString(&page, raw, d, "url", (*sp.PatchPagesPage).SetUrl)
	patchString(&page, raw, d, "time_zone", (*sp.PatchPagesPage).SetTimeZone)
	patchString(&page, raw, d, "branding", (*sp.PatchPagesPage).SetBranding)
	patchString(&page, raw, d, "notifications_from_email", (*sp.PatchPagesPage).SetNotificationsFromEmail)
	patchString(&page, raw, d, "notifications_email_footer", (*sp.PatchPagesPage).SetNotificationsEmailFooter)
	patchString(&page, raw, d, "css_body_background_color", (*sp.PatchPagesPage).SetCssBodyBackgroundColor)
	patchString(&page, raw, d, "css_font_color", (*sp.PatchPagesPage).SetCssFontColor)
	patchString(&page, raw, d, "css_light_font_color", (*sp.PatchPagesPage).SetCssLightFontColor)
	patchString(&page, raw, d, "css_greens", (*sp.PatchPagesPage).SetCssGreens)
	patchString(&page, raw, d, "css_yellows", (*sp.PatchPagesPage).SetCssYellows)
	patchString(&page, raw, d, "css_oranges", (*sp.PatchPagesPage).SetCssOranges)
	patchString(&page, raw, d, "css_reds", (*sp.PatchPagesPage).SetCssReds)
	patchString(&page, raw, d, "css_blues", (*sp.PatchPagesPage).SetCssBlues)
	patchString(&page, raw, d, "css_border_color", (*sp.PatchPagesPage).SetCssBorderColor)
	patchString(&page, raw, d, "css_graph_color", (*sp.PatchPagesPage).SetCssGraphColor)
	patchString(&page, raw, d, "css_link_color", (*sp.PatchPagesPage).SetCssLinkColor)
	patchString(&page, raw, d, "css_no_data", (*sp.PatchPagesPage).SetCssNoData)

	patchBool(&page, raw, d, "hidden_from_search", (*sp.PatchPagesPage).SetHiddenFromSearch)
	patchBool(&page, raw, d, "viewers_must_be_team_members", (*sp.PatchPagesPage).SetViewersMustBeTeamMembers)
	patchBool(&page, raw, d, "allow_page_subscribers", (*sp.PatchPagesPage).SetAllowPageSubscribers)
	patchBool(&page, raw, d, "allow_incident_subscribers", (*sp.PatchPagesPage).SetAllowIncidentSubscribers)
	patchBool(&page, raw, d, "allow_email_subscribers", (*sp.PatchPagesPage).SetAllowEmailSubscribers)
	patchBool(&page, raw, d, "allow_sms_subscribers", (*sp.PatchPagesPage).SetAllowSmsSubscribers)
	patchBool(&page, raw, d, "allow_rss_atom_feeds", (*sp.PatchPagesPage).SetAllowRssAtomFeeds)
	patchBool(&page, raw, d, "allow_webhook_subscribers", (*sp.PatchPagesPage).SetAllowWebhookSubscribers)

	return page
}

// patchString sets a PatchPagesPage string field only when the practitioner
// actually wrote attr in config (raw.GetAttr(attr) is non-null). See
// buildPagePatchBody for why that check has to run against raw config.
func patchString(page *sp.PatchPagesPage, raw cty.Value, d *schema.ResourceData, attr string, set func(*sp.PatchPagesPage, string)) {
	if v := raw.GetAttr(attr); !v.IsNull() {
		set(page, d.Get(attr).(string))
	}
}

// patchBool is patchString for booleans. It is a separate function rather
// than d.GetOkExists (the obvious shortcut) because GetOkExists is deprecated
// and cannot distinguish an explicit `false` from an attribute left out of
// config — exactly the distinction that matters for this resource's eight
// boolean toggles.
func patchBool(page *sp.PatchPagesPage, raw cty.Value, d *schema.ResourceData, attr string, set func(*sp.PatchPagesPage, bool)) {
	if v := raw.GetAttr(attr); !v.IsNull() {
		set(page, d.Get(attr).(bool))
	}
}

// setPageState maps a Page API response into Terraform state. Read and
// Update both funnel through this so state always reflects what the API
// actually returned.
func setPageState(d *schema.ResourceData, page *sp.Page) {
	d.SetId(page.GetId())
	d.Set("page_id", page.GetId())

	d.Set("name", page.GetName())
	d.Set("domain", page.GetDomain())
	d.Set("subdomain", page.GetSubdomain())
	d.Set("url", page.GetUrl())
	d.Set("time_zone", page.GetTimeZone())
	d.Set("branding", page.GetBranding())
	d.Set("hidden_from_search", page.GetHiddenFromSearch())
	d.Set("viewers_must_be_team_members", page.GetViewersMustBeTeamMembers())
	d.Set("allow_page_subscribers", page.GetAllowPageSubscribers())
	d.Set("allow_incident_subscribers", page.GetAllowIncidentSubscribers())
	d.Set("allow_email_subscribers", page.GetAllowEmailSubscribers())
	d.Set("allow_sms_subscribers", page.GetAllowSmsSubscribers())
	d.Set("allow_rss_atom_feeds", page.GetAllowRssAtomFeeds())
	d.Set("allow_webhook_subscribers", page.GetAllowWebhookSubscribers())
	d.Set("notifications_from_email", page.GetNotificationsFromEmail())
	d.Set("notifications_email_footer", page.GetNotificationsEmailFooter())
	d.Set("css_body_background_color", page.GetCssBodyBackgroundColor())
	d.Set("css_font_color", page.GetCssFontColor())
	d.Set("css_light_font_color", page.GetCssLightFontColor())
	d.Set("css_greens", page.GetCssGreens())
	d.Set("css_yellows", page.GetCssYellows())
	d.Set("css_oranges", page.GetCssOranges())
	d.Set("css_reds", page.GetCssReds())
	d.Set("css_blues", page.GetCssBlues())
	d.Set("css_border_color", page.GetCssBorderColor())
	d.Set("css_graph_color", page.GetCssGraphColor())
	d.Set("css_link_color", page.GetCssLinkColor())
	d.Set("css_no_data", page.GetCssNoData())

	// Read-only fields: see the "Computed-only" schema comment above for why
	// these are never sent.
	d.Set("created_at", page.GetCreatedAt().Format(time.RFC3339))
	d.Set("updated_at", page.GetUpdatedAt().Format(time.RFC3339))
	d.Set("page_description", page.GetPageDescription())
	d.Set("headline", page.GetHeadline())
	d.Set("support_url", page.GetSupportUrl())
	d.Set("twitter_username", page.GetTwitterUsername())
	d.Set("ip_restrictions", page.GetIpRestrictions())
	d.Set("city", page.GetCity())
	d.Set("state", page.GetState())
	d.Set("country", page.GetCountry())
	d.Set("activity_score", float64(page.GetActivityScore()))
}
