package duo

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/duosecurity/duo_api_golang"
	admin "github.com/duosecurity/duo_api_golang/admin"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIntegration() *schema.Resource {
	return &schema.Resource{
		Create: resourceIntegrationCreate,
		Read:   resourceIntegrationRead,
		Update: resourceIntegrationUpdate,
		Delete: resourceIntegrationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ikey": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

type Integration struct {
	Name string `json:"name"`
	Type string `json:"type"`
	IKey string `json:"integration_key"`
	SKey string `json:"secret_key"`
}

type IntegrationResult struct {
	duoapi.StatResult
	Response Integration
}

func resourceIntegrationCreate(d *schema.ResourceData, meta interface{}) error {
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	params := url.Values{}
	params.Set("name", d.Get("name").(string))
	params.Set("type", d.Get("type").(string))

	_, body, err := duoAdminClient.SignedCall("POST", "/admin/v1/integrations", params, duoapi.UseTimeout)
	if err != nil {
		return err
	}

	result := &IntegrationResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	if result.Stat != "OK" {
		return fmt.Errorf("could not create duo integration %s %s", result.Stat, *result.Message)
	}
	d.SetId(result.Response.IKey)
	d.Set("name", result.Response.Name)
	d.Set("type", result.Response.Type)
	d.Set("ikey", result.Response.IKey)
	d.Set("skey", result.Response.SKey)

	return resourceIntegrationRead(d, meta)
}

func resourceIntegrationRead(d *schema.ResourceData, meta interface{}) error {
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	iKey := d.Id()
	_, body, err := duoAdminClient.SignedCall("GET", fmt.Sprintf("/admin/v1/integrations/%s", iKey), nil, duoapi.UseTimeout)
	if err != nil {
		return err
	}

	result := &IntegrationResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}

	if result.Stat != "OK" {
		if *result.Message == "Resource not found" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read integration from duo %s, %s", result.Stat, *result.Message)
	}

	d.Set("name", result.Response.Name)
	d.Set("type", result.Response.Type)
	d.Set("ikey", result.Response.IKey)
	d.Set("skey", result.Response.SKey)

	return nil
}

func resourceIntegrationUpdate(d *schema.ResourceData, meta interface{}) error {
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	iKey := d.Id()

	d.Partial(true)
	if d.HasChange("name") {
		params := url.Values{}
		params.Set("name", d.Get("name").(string))
		_, body, err := duoAdminClient.SignedCall("POST", fmt.Sprintf("/admin/v1/integrations/%s", iKey), params, duoapi.UseTimeout)
		if err != nil {
			return err
		}
		result := &IntegrationResult{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return err
		}
		if result.Stat != "OK" {
			return fmt.Errorf("there was a problem updating integration %s's name: %s", iKey, *result.Message)
		}
	}
	d.Partial(false)
	return resourceIntegrationRead(d, meta)
}

func resourceIntegrationDelete(d *schema.ResourceData, meta interface{}) error {
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	iKey := d.Id()
	_, body, err := duoAdminClient.SignedCall("DELETE", fmt.Sprintf("/admin/v1/integrations/%s", iKey), nil, duoapi.UseTimeout)
	if err != nil {
		return err
	}
	var deleteResult struct {
		duoapi.StatResult
		Response string
	}
	err = json.Unmarshal(body, &deleteResult)
	if err != nil {
		return err
	}
	if deleteResult.Stat != "OK" {
		return fmt.Errorf("there was a problem deleting ikey %s: %s", iKey, *deleteResult.Message)
	}
	return nil
}
