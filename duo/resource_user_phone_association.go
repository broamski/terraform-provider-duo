package duo

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/duosecurity/duo_api_golang"
	admin "github.com/duosecurity/duo_api_golang/admin"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUserPhoneAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserPhoneAssociationCreate,
		Read:   resourceUserPhoneAssociationRead,
		Delete: resourceUserPhoneAssociationDelete,

		Schema: map[string]*schema.Schema{
			"user_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"phone_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

type AssociationResult struct {
	duoapi.StatResult
}

func resourceUserPhoneAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	pid := d.Get("phone_id").(string)
	uid := d.Get("user_id").(string)

	params := url.Values{}

	params.Set("phone_id", pid)

	_, body, err := duoAdminClient.SignedCall("POST", fmt.Sprintf("/admin/v1/users/%s/phones", uid), params, duoapi.UseTimeout)
	if err != nil {
		return err
	}

	result := &AssociationResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	if result.Stat != "OK" {
		return fmt.Errorf("could not associate phone to user %s %s", result.Stat, *result.Message)
	}
	d.SetId(fmt.Sprintf("%s-%s", d.Get("user_id").(string), d.Get("phone_id").(string)))
	return resourceUserPhoneAssociationRead(d, meta)
}

func resourceUserPhoneAssociationRead(d *schema.ResourceData, meta interface{}) error {
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)
	pid := d.Get("phone_id").(string)
	uid := d.Get("user_id").(string)
	result, err := duoAdminClient.GetPhone(pid)
	if err != nil {
		return err
	}
	if result.Stat != "OK" {
		return fmt.Errorf(fmt.Sprintf("could not find phone %s", pid))
	}

	var found bool
	var foundUser string
	for _, v := range result.Response.Users {
		if v.UserID == d.Get("user_id").(string) {
			found = true
			foundUser = v.UserID
		}
	}
	if !found {
		return fmt.Errorf("could not find phone %s attached to user %s", pid, uid)
	}
	d.Set("phone_id", pid)
	d.Set("user_id", foundUser)
	return nil
}

func resourceUserPhoneAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	duoclient := meta.(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	pid := d.Get("phone_id").(string)
	uid := d.Get("user_id").(string)
	_, body, err := duoAdminClient.SignedCall("DELETE", fmt.Sprintf("/admin/v1/users/%s/phones/%s", uid, pid), nil, duoapi.UseTimeout)
	if err != nil {
		return err
	}

	result := &AssociationResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	if result.Stat != "OK" {
		return fmt.Errorf("could not disassociate phone %s from user %s: %+v", pid, uid, *result.Message)
	}
	return nil
}
