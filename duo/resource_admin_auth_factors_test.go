package duo

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/duosecurity/duo_api_golang"
	admin "github.com/duosecurity/duo_api_golang/admin"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAdminAuthFactors_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAdminAuthFactorsDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAdminAuthFactorsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"duo_admin_auth_factors.test", "push_enabled", "true"),
					resource.TestCheckResourceAttr(
						"duo_admin_auth_factors.test", "mobile_otp_enabled", "false"),
					resource.TestCheckResourceAttr(
						"duo_admin_auth_factors.test", "hardware_token_enabled", "false"),
					resource.TestCheckResourceAttr(
						"duo_admin_auth_factors.test", "sms_enabled", "false"),
					resource.TestCheckResourceAttr(
						"duo_admin_auth_factors.test", "voice_enabled", "false"),
					resource.TestCheckResourceAttr(
						"duo_admin_auth_factors.test", "yubikey_enabled", "true"),
				),
			},
			resource.TestStep{
				Config: testAccCheckAdminAuthFactorsConfigUpdated(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"duo_admin_auth_factors.test", "push_enabled", "true"),
					resource.TestCheckResourceAttr(
						"duo_admin_auth_factors.test", "mobile_otp_enabled", "true"),
					resource.TestCheckResourceAttr(
						"duo_admin_auth_factors.test", "hardware_token_enabled", "true"),
					resource.TestCheckResourceAttr(
						"duo_admin_auth_factors.test", "sms_enabled", "true"),
					resource.TestCheckResourceAttr(
						"duo_admin_auth_factors.test", "voice_enabled", "true"),
					resource.TestCheckResourceAttr(
						"duo_admin_auth_factors.test", "yubikey_enabled", "true"),
				),
			},
		},
	})
}

func TestAccAdminAuthFactors_import(t *testing.T) {
	resourceName := "duo_admin_auth_factors.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAdminAuthFactorsDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAdminAuthFactorsConfig(),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAdminAuthFactorsDestroy(s *terraform.State) error {
	duoclient := testAccProvider.Meta().(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)
	for _, r := range s.RootModule().Resources {
		if r.Type != "duo_admin_auth_factors" {
			continue
		}
		_, body, err := duoAdminClient.SignedCall("GET", "/admin/v1/admins/allowed_auth_methods", nil, duoapi.UseTimeout)
		if err != nil {
			return err
		}

		result := &AdminAuthFactorsResult{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return err
		}

		if result.Response.Push != true && (result.Response.HardwareToken == false && result.Response.Yubikey == false &&
			result.Response.MobileOTP == false && result.Response.SMS == false && result.Response.Voice == false) {
			return fmt.Errorf("Undesired admin auth state: %+v", result.Response)
		}
	}
	return nil
}

func testAccCheckAdminAuthFactorsConfig() string {
	return `
resource "duo_admin_auth_factors" "test" {
  push_enabled = true
  mobile_otp_enabled = false
  hardware_token_enabled = false
  sms_enabled = false
  voice_enabled = false
  yubikey_enabled = true
}
`
}

func testAccCheckAdminAuthFactorsConfigUpdated() string {
	return `
resource "duo_admin_auth_factors" "test" {
	push_enabled = true
	mobile_otp_enabled = true
	hardware_token_enabled = true
	sms_enabled = true
	voice_enabled = true
	yubikey_enabled = true	
}
`
}
