package statuspage

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	sp "github.com/sbecker59/statuspage-api-client-go/api/v1/statuspage"
)

func resourceIncidentRead(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	name := d.Get("name").(string)
	log.Printf("[INFO] Reading Status Page incident '%s'", name)

	incident, _, err := statuspageClientV1.IncidentsApi.GetPagesPageIdIncidentsIncidentId(authV1, d.Get("page_id").(string), d.Id()).Execute()
	if err != nil {
		return TranslateClientErrorDiag(err, "failed to get incident using Status Page API")
	}

	if &incident == nil {
		log.Printf("[INFO] Statuspage could not find incident with ID: %s\n", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("name", incident.GetName())
	d.Set("status", incident.GetStatus())
	d.Set("impact_override", incident.GetImpactOverride())
	d.Set("scheduled_remind_prior", incident.GetScheduledRemindPrior())
	d.Set("scheduled_auto_in_progress", incident.GetScheduledAutoInProgress())
	d.Set("scheduled_auto_completed", incident.GetScheduledAutoCompleted())

	components := make([]interface{}, len(incident.GetComponents()))
	for i, statuspage_component := range incident.GetComponents() {
		component := make(map[string]interface{})
		component["id"] = statuspage_component.GetId()
		component["name"] = statuspage_component.GetName()
		component["status"] = statuspage_component.GetStatus()
		components[i] = component
	}
	d.Set("component", components)

	return nil
}

func resourceIncidentCreate(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	name := d.Get("name").(string)
	status := d.Get("status").(string)
	impact_override := d.Get("impact_override").(string)
	body := d.Get("body").(string)

	scheduled_remind_prior := d.Get("scheduled_remind_prior").(bool)
	scheduled_auto_in_progress := d.Get("scheduled_auto_in_progress").(bool)
	scheduled_auto_completed := d.Get("scheduled_auto_completed").(bool)

	terraformComponents := d.Get("component").(*schema.Set).List()

	var component_ids []string
	components := make(map[string]interface{})

	for _, terraformComponent := range terraformComponents {

		if id, ok := terraformComponent.(map[string]interface{})["id"].(string); ok && len(id) != 0 {
			component_ids = append(component_ids, id)
			if status, ok := terraformComponent.(map[string]interface{})["status"].(string); ok && len(status) != 0 {
				components[id] = status
			}
		}

	}

	var component sp.PostPagesPageIdIncidentsIncident

	component.SetName(name)
	component.SetStatus(status)
	component.SetImpactOverride(impact_override)
	component.SetBody(body)

	component.SetScheduledRemindPrior(scheduled_remind_prior)
	component.SetScheduledAutoInProgress(scheduled_auto_in_progress)
	component.SetScheduledAutoCompleted(scheduled_auto_completed)

	component.SetComponentIds(component_ids)
	component.SetComponents(components)

	o := *sp.NewPostPagesPageIdIncidents()
	o.SetIncident(component)

	log.Printf("[INFO] Creating Status Page incident '%s'", name)
	result, _, err := statuspageClientV1.IncidentsApi.PostPagesPageIdIncidents(authV1, d.Get("page_id").(string)).PostPagesPageIdIncidents(o).Execute()

	if err != nil {
		return TranslateClientErrorDiag(err, "failed to create incident using Status Page API")
	}

	d.SetId(result.GetId())

	return resourceIncidentRead(d, m)

}

func resourceIncidentUpdate(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	name := d.Get("name").(string)
	status := d.Get("status").(string)
	impact_override := d.Get("impact_override").(string)
	body := d.Get("body").(string)

	scheduled_remind_prior := d.Get("scheduled_remind_prior").(bool)
	scheduled_auto_in_progress := d.Get("scheduled_auto_in_progress").(bool)
	scheduled_auto_completed := d.Get("scheduled_auto_completed").(bool)

	terraformComponents := d.Get("component").(*schema.Set).List()

	var component_ids []string
	components := make(map[string]interface{})

	for _, terraformComponent := range terraformComponents {

		if id, ok := terraformComponent.(map[string]interface{})["id"].(string); ok && len(id) != 0 {
			component_ids = append(component_ids, id)
			if status, ok := terraformComponent.(map[string]interface{})["status"].(string); ok && len(status) != 0 {
				components[id] = status
			}
		}

	}

	var component sp.PatchPagesPageIdIncidentsIncident

	component.SetName(name)
	component.SetStatus(status)
	component.SetImpactOverride(impact_override)
	component.SetBody(body)

	component.SetScheduledRemindPrior(scheduled_remind_prior)
	component.SetScheduledAutoInProgress(scheduled_auto_in_progress)
	component.SetScheduledAutoCompleted(scheduled_auto_completed)

	component.SetComponentIds(component_ids)
	component.SetComponents(components)

	o := *sp.NewPatchPagesPageIdIncidents()
	o.SetIncident(component)

	log.Printf("[INFO] Update Status Page incident '%s'", name)
	result, _, err := statuspageClientV1.IncidentsApi.PatchPagesPageIdIncidentsIncidentId(authV1, d.Get("page_id").(string), d.Id()).PatchPagesPageIdIncidents(o).Execute()

	if err != nil {
		return TranslateClientErrorDiag(err, "failed to update incident using Status Page API")
	}

	d.SetId(result.GetId())

	return resourceIncidentRead(d, m)
}

func resourceIncidentDelete(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	_, _, err := statuspageClientV1.IncidentsApi.DeletePagesPageIdIncidentsIncidentId(authV1, d.Get("page_id").(string), d.Id()).Execute()

	if err != nil {
		return TranslateClientErrorDiag(err, "failed to delete incident using Status Page API")
	}

	return nil
}

func resourceIncidentImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	if len(strings.Split(d.Id(), "/")) != 2 {
		return []*schema.ResourceData{}, fmt.Errorf("[ERROR] Invalid resource format: %s. Please use 'page-id/incident-id'", d.Id())
	}

	pageID := strings.Split(d.Id(), "/")[0]
	incidentID := strings.Split(d.Id(), "/")[1]

	log.Printf("[INFO] Importing Incident %s from Page %s", incidentID, pageID)

	d.Set("page_id", pageID)
	d.SetId(incidentID)

	err := resourceIncidentRead(d, m)
	return []*schema.ResourceData{d}, err

}

func resourceIncident() *schema.Resource {
	return &schema.Resource{
		Create: resourceIncidentCreate,
		Read:   resourceIncidentRead,
		Update: resourceIncidentUpdate,
		Delete: resourceIncidentDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIncidentImport,
		},
		Schema: map[string]*schema.Schema{
			"page_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the ID of the page this component belongs to",
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Incident Name",
				Required:    true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The incident status. For realtime incidents, valid values are investigating, identified, monitoring, and resolved. For scheduled incidents, valid values are scheduled, in_progress, verifying, and completed.",
				ValidateFunc: validation.StringInSlice([]string{"investigating", "identified", "monitoring", "resolved"}, false),
				Default:      "investigating",
			},
			"impact_override": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "value to override calculated impact value",
				ValidateFunc: validation.StringInSlice([]string{"maintenance", "none", "critical", "major", "minor"}, false),
				Default:      "none",
			},
			"scheduled_remind_prior": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"scheduled_auto_in_progress": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"scheduled_auto_completed": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"component": {
				Type:        schema.TypeSet,
				Description: "List of component_ids affected by this incident",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Identifier for component",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"name": {
							Type:         schema.TypeString,
							Description:  "Display name for component",
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"status": {
							Type:         schema.TypeString,
							Description:  "Status of component",
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"operational", "under_maintenance", "degraded_performance", "partial_outage", "major_outage", ""}, false),
							Default:      "operational",
						},
					},
				},
			},
			"body": {
				Type:        schema.TypeString,
				Description: "The initial message, created as the first incident update",
				Optional:    true,
			},
		},
	}
}
