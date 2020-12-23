package statuspage

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	sp "github.com/sbecker59/statuspage-api-client-go/api/v1/statuspage"
)

func resourceComponentRead(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	name := d.Get("name").(string)
	log.Printf("[INFO] Reading Status Page component '%s'", name)

	component, _, err := statuspageClientV1.ComponentsApi.GetPagesPageIdComponentsComponentId(authV1, d.Get("page_id").(string), d.Id()).Execute()
	if err.Error() != "" {
		return translateClientError(err, "failed to get component using Status Page API")
	}

	if &component == nil {
		log.Printf("[INFO] Statuspage could not find component with ID: %s\n", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("description", component.GetDescription())
	d.Set("name", component.GetName())
	d.Set("only_show_if_degraded", component.GetOnlyShowIfDegraded())
	d.Set("showcase", component.GetShowcase())
	d.Set("status", component.GetStatus())

	return nil
}

func resourceComponentCreate(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	name := d.Get("name").(string)
	status := d.Get("status").(string)
	description := d.Get("description").(string)
	showcase := d.Get("showcase").(bool)
	onlyShowIfDegraded := d.Get("only_show_if_degraded").(bool)

	var component sp.PostPagesPageIdComponentsComponent

	component.SetName(name)
	component.SetDescription(description)
	component.SetStatus(status)
	component.SetOnlyShowIfDegraded(onlyShowIfDegraded)
	component.SetShowcase(showcase)
	if r, ok := d.GetOk("start_date"); ok {
		component.SetStartDate(r.(string))
	}

	o := *sp.NewPostPagesPageIdComponents()
	o.SetComponent(component)

	log.Printf("[INFO] Creating Status Page componant '%s'", name)
	result, _, err := statuspageClientV1.ComponentsApi.PostPagesPageIdComponents(authV1, d.Get("page_id").(string)).PostPagesPageIdComponents(o).Execute()

	if err.Error() != "" {
		return translateClientError(err, "failed to create component using Status Page API")
	}

	d.SetId(result.GetId())

	return resourceComponentRead(d, m)

}

func resourceComponentUpdate(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	name := d.Get("name").(string)
	status := d.Get("status").(string)
	description := d.Get("description").(string)
	showcase := d.Get("showcase").(bool)
	onlyShowIfDegraded := d.Get("only_show_if_degraded").(bool)

	var component sp.PostPagesPageIdComponentsComponent

	component.SetName(name)
	component.SetDescription(description)
	component.SetStatus(status)
	component.SetOnlyShowIfDegraded(onlyShowIfDegraded)
	component.SetShowcase(showcase)
	if r, ok := d.GetOk("start_date"); ok {
		component.SetStartDate(r.(string))
	}

	o := *sp.NewPatchPagesPageIdComponents()
	o.SetComponent(component)

	log.Printf("[INFO] Update Status Page componant '%s'", name)
	result, _, err := statuspageClientV1.ComponentsApi.PatchPagesPageIdComponentsComponentId(authV1, d.Get("page_id").(string), d.Id()).PatchPagesPageIdComponents(o).Execute()

	if err.Error() != "" {
		return translateClientError(err, "failed to update component using Status Page API")
	}

	d.SetId(result.GetId())

	return resourceComponentRead(d, m)
}

func resourceComponentDelete(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	_, err := statuspageClientV1.ComponentsApi.DeletePagesPageIdComponentsComponentId(authV1, d.Get("page_id").(string), d.Id()).Execute()

	if err.Error() != "" {
		return translateClientError(err, "failed to delete component using Status Page API")
	}

	return nil
}

func resourceComponent() *schema.Resource {
	return &schema.Resource{
		Create: resourceComponentCreate,
		Read:   resourceComponentRead,
		Update: resourceComponentUpdate,
		Delete: resourceComponentDelete,
		Schema: map[string]*schema.Schema{
			"page_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the ID of the page this component belongs to",
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Display Name for the component",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "More detailed description for the component",
				Optional:    true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"operational", "under_maintenance", "degraded_performance", "partial_outage", "major_outage", ""}, false),
				Default:      "operational",
			},
			"showcase": {
				Type:        schema.TypeBool,
				Description: "Should this component be shown component only if in degraded state",
				Optional:    true,
			},
			"only_show_if_degraded": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"start_date": {
				Type:        schema.TypeString,
				Description: "Should this component be showcased",
				Optional:    true,
			},
		},
	}
}
