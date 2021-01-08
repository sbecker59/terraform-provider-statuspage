package statuspage

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	sp "github.com/sbecker59/statuspage-api-client-go/api/v1/statuspage"
)

func resourceSubscriberRead(d *schema.ResourceData, m interface{}) error {
	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	resp, _, err := statuspageClientV1.SubscribersApi.GetPagesPageIdSubscribersSubscriberId(authV1, d.Get("page_id").(string), d.Id()).Execute()
	if err.Error() != "" {
		return translateClientError(err, "failed to get component groups using Status Page API")
	}

	if &resp == nil {
		log.Printf("[INFO] Statuspage could not find subscriber with ID: %s\n", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("email", resp.GetEmail())
	d.Set("endpoint", resp.GetEndpoint())

	return nil

}

func resourceSubscriberCreate(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	var subscriber sp.PostPagesPageIdSubscribersSubscriber

	if r, ok := d.GetOk("email"); ok {
		subscriber.SetEmail(r.(string))
	}
	if r, ok := d.GetOk("endpoint"); ok {
		subscriber.SetEndpoint(r.(string))
	}

	o := *sp.NewPostPagesPageIdSubscribers()
	o.SetSubscriber(subscriber)

	result, _, err := statuspageClientV1.SubscribersApi.PostPagesPageIdSubscribers(authV1, d.Get("page_id").(string)).PostPagesPageIdSubscribers(o).Execute()

	if err.Error() != "" {
		return translateClientError(err, "failed to create subscriber using Status Page API")
	}

	d.SetId(result.GetId())

	return resourceSubscriberRead(d, m)

}

func resourceSubscriberUpdate(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	o := *sp.NewPatchPagesPageIdSubscribers()

	result, _, err := statuspageClientV1.SubscribersApi.PatchPagesPageIdSubscribersSubscriberId(authV1, d.Get("page_id").(string), d.Id()).PatchPagesPageIdSubscribers(o).Execute()

	if err.Error() != "" {
		return translateClientError(err, "failed to update subscriber using Status Page API")
	}

	d.SetId(result.GetId())

	return resourceSubscriberRead(d, m)

}

func resourceSubscriberDelete(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	_, _, err := statuspageClientV1.SubscribersApi.DeletePagesPageIdSubscribersSubscriberId(authV1, d.Get("page_id").(string), d.Id()).Execute()

	if err.Error() != "" {
		return translateClientError(err, "failed to delete subscriber using Status Page API")
	}

	return nil

}

func resourceSubscriber() *schema.Resource {
	return &schema.Resource{
		Create: resourceSubscriberCreate,
		Read:   resourceSubscriberRead,
		Update: resourceSubscriberUpdate,
		Delete: resourceSubscriberDelete,
		Schema: map[string]*schema.Schema{
			"page_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the ID of the page this component belongs to",
				ForceNew:    true,
			},
			"email": {
				Type:        schema.TypeString,
				Description: "the email address for creating Email and Webhook subscribers",
				Optional:    true,
			},
			"endpoint": {
				Type:        schema.TypeString,
				Description: "The endpoint URI for creating Webhook subscribers",
				Optional:    true,
			},
			// "phone_country": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Description: "The two-character country where the phone number is located to use for the new SMS subscriber",
			// },
			// "phone_number": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Description: "The phone number (as you would dial from the phone_country) to use for the new SMS subscriber",
			// },
			// "skip_confirmation_notification": {
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// 	Description: "If this is true, do not notify the user with changes to their subscription.",
			// },
			// "components": {
			// 	Type:        schema.TypeSet,
			// 	Optional:    true,
			// 	Description: "The components for which the subscriber has elected to receive updates.",
			// 	Set:         schema.HashString,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeString,
			// 	},
			// },
		},
	}
}
