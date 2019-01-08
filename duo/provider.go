package duo

import (
	"errors"
	"os"

	"github.com/duosecurity/duo_api_golang"
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"ikey": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DUO_IKEY", nil),
				Description: "Duo AdminAPI Integration ikey",
			},
			"skey": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Duo AdminAPI Integration skey",
			},
			"api_host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DUO_API_HOST", nil),
				Description: "Duo AdminAPI Integration API Server",
			},
		},
		ConfigureFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"duo_admin":                  resourceAdmin(),
			"duo_admin_auth_factors":     resourceAdminAuthFactors(),
			"duo_integration":            resourceIntegration(),
			"duo_user":                   resourceUser(),
			"duo_phone":                  resourcePhone(),
			"duo_user_phone_association": resourceUserPhoneAssociation(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	skey := d.Get("skey").(string)
	if skey != "" {
		return nil, errors.New("In order to keep the skey secret, you should NOT provide this value via config but rather the DUO_SKEY env var")
	}
	skey = os.Getenv("DUO_SKEY")
	if skey == "" {
		return nil, errors.New("DUO_SKEY is missing")
	}

	ikey := d.Get("ikey").(string)
	apiHost := d.Get("api_host").(string)

	duoClient := duoapi.NewDuoApi(
		ikey,
		skey,
		apiHost,
		"terraform-provider-duo",
	)
	return duoClient, nil
}

type deleteResult struct {
	duoapi.StatResult
	Response string
}
