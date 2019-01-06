package duo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/duosecurity/duo_api_golang"
	admin "github.com/duosecurity/duo_api_golang/admin"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAdminAuthFactors() *schema.Resource {
	return &schema.Resource{
		Create: resourceAdminAuthFactorsCreate,
		Read:   resourceAdminAuthFactorsRead,
		Update: resourceAdminAuthFactorsCreate,
		Delete: resourceAdminAuthFactorsDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"hardware_token_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"mobile_otp_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"push_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"sms_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"voice_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"yubikey_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

type AuthFactors struct {
	HardwareToken bool `json:"hardware_token_enabled"`
	MobileOTP     bool `json:"mobile_otp_enabled"`
	Push          bool `json:"push_enabled"`
	SMS           bool `json:"sms_enabled"`
	Voice         bool `json:"voice_enabled"`
	Yubikey       bool `json:"yubikey_enabled"`
}

type AdminAuthFactorsResult struct {
	duoapi.StatResult
	Response AuthFactors
}

func boolParser(input interface{}) string {
	s := strconv.FormatBool(input.(bool))
	return s
}

func resourceAdminAuthFactorsCreate(d *schema.ResourceData, meta interface{}) error {
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	params := url.Values{}
	params.Set("hardware_token_enabled", boolParser(d.Get("hardware_token_enabled")))
	params.Set("mobile_otp_enabled", boolParser(d.Get("mobile_otp_enabled")))
	params.Set("push_enabled", boolParser(d.Get("push_enabled")))
	params.Set("sms_enabled", boolParser(d.Get("sms_enabled")))
	params.Set("voice_enabled", boolParser(d.Get("voice_enabled")))
	params.Set("yubikey_enabled", boolParser(d.Get("yubikey_enabled")))

	_, body, err := duoAdminClient.SignedCall("POST", "/admin/v1/admins/allowed_auth_methods", params, duoapi.UseTimeout)
	if err != nil {
		return err
	}

	result := &AdminAuthFactorsResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	if result.Stat != "OK" {
		return fmt.Errorf("could not set admin auth methods %s", *result.Message)
	}
	d.SetId("admin_auth_factors")

	d.Set("hardware_token_enabled", result.Response.HardwareToken)
	d.Set("mobile_otp_enabled", result.Response.MobileOTP)
	d.Set("push_enabled", result.Response.Push)
	d.Set("sms_enabled", result.Response.SMS)
	d.Set("voice_enabled", result.Response.Voice)
	d.Set("yubikey_enabled", result.Response.Yubikey)

	return resourceAdminAuthFactorsRead(d, meta)
}

func resourceAdminAuthFactorsRead(d *schema.ResourceData, meta interface{}) error {
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	_, body, err := duoAdminClient.SignedCall("GET", "/admin/v1/admins/allowed_auth_methods", nil, duoapi.UseTimeout)
	if err != nil {
		return err
	}

	result := &AdminAuthFactorsResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	if result.Stat != "OK" {
		return fmt.Errorf("could not read allowed auth methods from duo %s, %s", result.Stat, *result.Message)
	}

	d.Set("hardware_token_enabled", result.Response.HardwareToken)
	d.Set("mobile_otp_enabled", result.Response.MobileOTP)
	d.Set("push_enabled", result.Response.Push)
	d.Set("sms_enabled", result.Response.SMS)
	d.Set("voice_enabled", result.Response.Voice)
	d.Set("yubikey_enabled", result.Response.Yubikey)

	return nil
}

func resourceAdminAuthFactorsDelete(d *schema.ResourceData, meta interface{}) error {
	//Auth Factors aren't a resoruce, so we'll revert to a state of push only
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	params := url.Values{}
	params.Set("hardware_token_enabled", "false")
	params.Set("mobile_otp_enabled", "false")
	params.Set("push_enabled", "true")
	params.Set("sms_enabled", "false")
	params.Set("voice_enabled", "false")
	params.Set("yubikey_enabled", "false")

	_, body, err := duoAdminClient.SignedCall("POST", "/admin/v1/admins/allowed_auth_methods", params, duoapi.UseTimeout)
	if err != nil {
		return err
	}

	result := &AdminAuthFactorsResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	if result.Stat != "OK" {
		return fmt.Errorf("could not set admin auth methods %s", *result.Message)
	}
	d.SetId("")
	return nil
}
