package duo

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/duosecurity/duo_api_golang"
	admin "github.com/duosecurity/duo_api_golang/admin"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"alias1": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"alias2": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"alias3": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"alias4": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"realname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "active",
				Optional: true,
			},
			"notes": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, meta interface{}) error {
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	params := url.Values{}
	params.Set("username", d.Get("username").(string))

	if d.Get("alias1") != "" {
		params.Set("alias1", d.Get("alias1").(string))
	}
	if d.Get("alias2") != "" {
		params.Set("alias2", d.Get("alias2").(string))
	}
	if d.Get("alias3") != "" {
		params.Set("alias3", d.Get("alias3").(string))
	}
	if d.Get("alias4") != "" {
		params.Set("alias4", d.Get("alias4").(string))
	}
	params.Set("realname", d.Get("realname").(string))
	params.Set("email", d.Get("email").(string))
	params.Set("status", d.Get("status").(string))
	params.Set("notes", d.Get("notes").(string))

	_, body, err := duoAdminClient.SignedCall("POST", "/admin/v1/users", params, duoapi.UseTimeout)
	if err != nil {
		return err
	}
	result := &admin.GetUserResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	if result.Stat != "OK" {
		return fmt.Errorf("could not create user %s %s", result.Stat, *result.Message)
	}

	user := result.Response
	d.SetId(user.UserID)
	return resourceUserRead(d, meta)
}

func resourceUserRead(d *schema.ResourceData, meta interface{}) error {
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	uid := d.Id()

	result, err := duoAdminClient.GetUser(uid)
	if err != nil {
		return err
	}

	if result.Stat != "OK" {
		if *result.Message == "Resource not found" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read user from duo %s, %s", result.Stat, *result.Message)
	}
	user := result.Response
	d.Set("username", user.Username)
	d.Set("alias1", user.Alias1)
	d.Set("alias2", user.Alias2)
	d.Set("alias3", user.Alias3)
	d.Set("alias4", user.Alias4)
	d.Set("realname", user.RealName)
	d.Set("email", user.Email)
	d.Set("status", user.Status)
	d.Set("notes", user.Notes)
	d.Set("user_id", user.UserID)
	return nil
}

func resourceUserUpdate(d *schema.ResourceData, meta interface{}) error {
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	userID := d.Id()
	params := url.Values{}

	d.Partial(true)
	if d.HasChange("username") {
		params.Set("username", d.Get("username").(string))
	}
	if d.HasChange("alias1") {
		params.Set("alias1", d.Get("alias1").(string))
	}
	if d.HasChange("alias2") {
		params.Set("alias2", d.Get("alias2").(string))
	}
	if d.HasChange("alias3") {
		params.Set("alias3", d.Get("alias3").(string))
	}
	if d.HasChange("alias4") {
		params.Set("alias4", d.Get("alias4").(string))
	}
	if d.HasChange("realname") {
		params.Set("realname", d.Get("realname").(string))
	}
	if d.HasChange("email") {
		params.Set("email", d.Get("email").(string))
	}
	if d.HasChange("status") {
		params.Set("status", d.Get("status").(string))
	}
	if d.HasChange("notes") {
		params.Set("notes", d.Get("notes").(string))
	}
	_, body, err := duoAdminClient.SignedCall("POST", fmt.Sprintf("/admin/v1/users/%s", userID), params, duoapi.UseTimeout)
	if err != nil {
		return err
	}
	result := &admin.GetUserResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	if result.Stat != "OK" {
		return fmt.Errorf("there was a problem updating user %s's name: %s", userID, *result.Message)
	}
	d.Partial(false)
	return resourceUserRead(d, meta)
}

func resourceUserDelete(d *schema.ResourceData, meta interface{}) error {
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	userID := d.Id()
	_, body, err := duoAdminClient.SignedCall("DELETE", fmt.Sprintf("/admin/v1/users/%s", userID), nil, duoapi.UseTimeout)
	if err != nil {
		return err
	}

	var result deleteResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	if result.Stat != "OK" {
		return fmt.Errorf("there was a problem deleting user %s: %s", userID, *result.Message)
	}
	return nil
}
