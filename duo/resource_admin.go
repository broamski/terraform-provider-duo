package duo

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/duosecurity/duo_api_golang"
	admin "github.com/duosecurity/duo_api_golang/admin"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAdmin() *schema.Resource {
	return &schema.Resource{
		Create: resourceAdminCreate,
		Read:   resourceAdminRead,
		Update: resourceAdminUpdate,
		Delete: resourceAdminDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"phone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

type Admin struct {
	AdminID string `json:"admin_id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Role    string `json:"role"`
}

type AdminResult struct {
	duoapi.StatResult
	Response Admin
}

func resourceAdminCreate(d *schema.ResourceData, meta interface{}) error {
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	params := url.Values{}
	params.Set("email", d.Get("email").(string))
	if d.Get("password").(string) == "" {
		b := make([]byte, 40)
		_, err := rand.Read(b)
		if err != nil {
			return err
		}
		h := sha1.New()
		h.Write(b)
		rnd := hex.EncodeToString(h.Sum(nil))
		log.Print("[DEBUG] User's temporary password is", rnd)
		params.Set("password", rnd)
	} else {
		params.Set("password", d.Get("password").(string))
	}
	params.Set("name", d.Get("name").(string))
	params.Set("phone", d.Get("phone").(string))
	if d.Get("role").(string) != "" {
		params.Set("role", d.Get("role").(string))
	}

	_, body, err := duoAdminClient.SignedCall("POST", "/admin/v1/admins", params, duoapi.UseTimeout)
	if err != nil {
		return err
	}

	result := &AdminResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	if result.Stat != "OK" {
		return fmt.Errorf("could not create duo user %s: %s", result.Stat, *result.Message)
	}
	adminID := result.Response.AdminID
	d.SetId(adminID)
	return resourceAdminRead(d, meta)
}

func resourceAdminRead(d *schema.ResourceData, meta interface{}) error {
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	_, body, err := duoAdminClient.SignedCall("GET", fmt.Sprintf("/admin/v1/admins/%s", d.Id()), nil, duoapi.UseTimeout)
	if err != nil {
		return err
	}

	result := &AdminResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	if result.Stat != "OK" {
		if *result.Message == "Resource not found" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read user from duo %s: %s", result.Stat, *result.Message)
	}

	d.Set("email", result.Response.Email)
	d.Set("name", result.Response.Name)
	d.Set("phone", result.Response.Phone)
	d.Set("role", result.Response.Role)
	return nil
}

func resourceAdminUpdate(d *schema.ResourceData, meta interface{}) error {
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	d.Partial(true)

	if d.HasChange("name") && !d.IsNewResource() {
		params := url.Values{}
		params.Set("name", d.Get("name").(string))
		_, body, err := duoAdminClient.SignedCall("POST", fmt.Sprintf("/admin/v1/admins/%s", d.Id()), params, duoapi.UseTimeout)
		if err != nil {
			return err
		}
		result := &AdminResult{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return err
		}
		if result.Stat != "OK" {
			if *result.Message == "Resource not found" {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("there was a problem updating user %s's name: %s", d.Id(), *result.Message)
		}
		d.SetPartial("name")
	}
	if d.HasChange("phone") && !d.IsNewResource() {
		params := url.Values{}
		params.Set("phone", d.Get("phone").(string))
		_, body, err := duoAdminClient.SignedCall("POST", fmt.Sprintf("/admin/v1/admins/%s", d.Id()), params, duoapi.UseTimeout)
		if err != nil {
			return err
		}
		result := &AdminResult{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return err
		}
		if result.Stat != "OK" {
			return fmt.Errorf("there was a problem updating user %s's phone: %s", d.Id(), *result.Message)
		}
		d.SetPartial("phone")
	}
	if d.HasChange("role") && !d.IsNewResource() {
		params := url.Values{}
		params.Set("role", d.Get("role").(string))
		_, body, err := duoAdminClient.SignedCall("POST", fmt.Sprintf("/admin/v1/admins/%s", d.Id()), params, duoapi.UseTimeout)
		if err != nil {
			return err
		}
		result := &AdminResult{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return err
		}
		if result.Stat != "OK" {
			return fmt.Errorf("there was a problem updating user %s's role: %s", d.Id(), *result.Message)
		}
		d.SetPartial("role")
	}
	d.Partial(false)
	return resourceAdminRead(d, meta)
}

func resourceAdminDelete(d *schema.ResourceData, meta interface{}) error {
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	adminID := d.Id()
	_, body, err := duoAdminClient.SignedCall("DELETE", fmt.Sprintf("/admin/v1/admins/%s", adminID), nil, duoapi.UseTimeout)
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
		return fmt.Errorf("there was a problem deleting user %s: %s", adminID, *deleteResult.Message)
	}
	d.SetId("")
	return nil
}
