package duo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/duosecurity/duo_api_golang"
	admin "github.com/duosecurity/duo_api_golang/admin"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePhone() *schema.Resource {
	return &schema.Resource{
		Create: resourcePhoneCreate,
		Read:   resourcePhoneRead,
		Update: resourcePhoneUpdate,
		Delete: resourcePhoneDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"number": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"extension": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"platform": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"predelay": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"postdelay": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"phone_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourcePhoneCreate(d *schema.ResourceData, meta interface{}) error {
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	params := url.Values{}
	if number, ok := d.GetOk("number"); ok {
		params.Set("number", number.(string))
	}
	if name, ok := d.GetOk("name"); ok {
		params.Set("name", name.(string))
	}
	if ptype, ok := d.GetOk("type"); ok {
		params.Set("type", ptype.(string))
	}
	if extension, ok := d.GetOk("extension"); ok {
		params.Set("extension", extension.(string))
	}
	if platform, ok := d.GetOk("platform"); ok {
		params.Set("platform", platform.(string))
	}
	if predelay, ok := d.GetOk("predelay"); ok {
		params.Set("predelay", predelay.(string))
	}
	if postdelay, ok := d.GetOk("postdelay"); ok {
		params.Set("postdelay", postdelay.(string))
	}
	_, body, err := duoAdminClient.SignedCall("POST", "/admin/v1/phones", params, duoapi.UseTimeout)
	if err != nil {
		return err
	}
	result := &admin.GetPhoneResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}

	if result.Stat != "OK" {
		code := strconv.Itoa(int(*result.Code))
		if strings.HasPrefix(code, "400") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not create phone %s %s", result.Stat, *result.Message)
	}

	d.SetId(result.Response.PhoneID)
	return resourcePhoneRead(d, meta)
}

func resourcePhoneRead(d *schema.ResourceData, meta interface{}) error {
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	pid := d.Id()

	result, err := duoAdminClient.GetPhone(pid)
	if err != nil {
		return err
	}

	if result.Stat != "OK" {
		if *result.Message == "Resource not found" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read phone from duo %s, %s", result.Stat, *result.Message)
	}

	d.Set("number", result.Response.Number)
	d.Set("name", result.Response.Name)
	d.Set("type", result.Response.Type)
	d.Set("extension", result.Response.Extension)
	d.Set("platform", result.Response.Platform)
	d.Set("predelay", result.Response.Predelay)
	d.Set("postdelay", result.Response.Postdelay)
	d.Set("phone_id", result.Response.PhoneID)
	return nil
}

func resourcePhoneUpdate(d *schema.ResourceData, meta interface{}) error {
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	pid := d.Id()
	params := url.Values{}

	d.Partial(true)
	if d.HasChange("number") {
		params.Set("number", d.Get("number").(string))
	}
	if d.HasChange("name") {
		params.Set("name", d.Get("name").(string))
	}
	if d.HasChange("type") {
		params.Set("type", d.Get("type").(string))
	}
	if d.HasChange("extension") {
		params.Set("extension", d.Get("extension").(string))
	}
	if d.HasChange("platform") {
		params.Set("platform", d.Get("platform").(string))
	}
	if d.HasChange("predelay") {
		params.Set("predelay", d.Get("predelay").(string))
	}
	if d.HasChange("postdelay") {
		params.Set("postdelay", d.Get("postdelay").(string))
	}

	_, body, err := duoAdminClient.SignedCall("POST", fmt.Sprintf("/admin/v1/phones/%s", pid), params, duoapi.UseTimeout)
	if err != nil {
		return err
	}
	result := &admin.GetPhoneResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	if result.Stat != "OK" {
		return fmt.Errorf("there was a problem updating phone %s: %s", pid, *result.Message)
	}
	d.Partial(false)
	return resourcePhoneRead(d, meta)
}

func resourcePhoneDelete(d *schema.ResourceData, meta interface{}) error {
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	pid := d.Id()
	_, body, err := duoAdminClient.SignedCall("DELETE", fmt.Sprintf("/admin/v1/phones/%s", pid), nil, duoapi.UseTimeout)
	if err != nil {
		return err
	}
	var result deleteResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	if result.Stat != "OK" {
		return fmt.Errorf("there was a problem deleting phone %s: %s", pid, *result.Message)
	}
	return nil
}
