package duo

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/duosecurity/duo_api_golang"
	admin "github.com/duosecurity/duo_api_golang/admin"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAdmin_Basic(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAdminDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAdminConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdminExists("duo_admin.test"),
					resource.TestCheckResourceAttr(
						"duo_admin.test", "name", fmt.Sprintf("test-%d", rInt)),
					resource.TestCheckResourceAttr(
						"duo_admin.test", "email", "le1f@wut.wut"),
					resource.TestCheckResourceAttr(
						"duo_admin.test", "phone", "+12813308004"),
				),
			},
			resource.TestStep{
				Config: testAccCheckAdminConfigUpdated(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdminExists("duo_admin.test"),
					resource.TestCheckResourceAttr(
						"duo_admin.test", "name", fmt.Sprintf("test-updated-%d", rInt)),
					resource.TestCheckResourceAttr(
						"duo_admin.test", "email", "le1f@wut.wut"),
					resource.TestCheckResourceAttr(
						"duo_admin.test", "phone", "+12813308004"),
				),
			},
		},
	})
}

func TestAccAdmin_import(t *testing.T) {
	resourceName := "duo_admin.test"
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAdminDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAdminConfig(rInt),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAdminDestroy(s *terraform.State) error {
	duoclient := testAccProvider.Meta().(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)
	for _, r := range s.RootModule().Resources {
		if r.Type != "duo_admin" {
			continue
		}
		_, body, err := duoAdminClient.SignedCall("GET", fmt.Sprintf("/admin/v1/admins/%s", r.Primary.ID), nil, duoapi.UseTimeout)
		if err != nil {
			return err
		}

		result := &AdminResult{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return err
		}
		if result.Stat == "OK" {
			return fmt.Errorf("Found undeleted admin: %s", result.Response)
		}
	}
	return nil
}
func testAccCheckAdminExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		duoclient := testAccProvider.Meta().(*duoapi.DuoApi)
		duoAdminClient := admin.New(*duoclient)

		_, body, err := duoAdminClient.SignedCall("GET", fmt.Sprintf("/admin/v1/admins/%s", rs.Primary.ID), nil, duoapi.UseTimeout)
		if err != nil {
			return err
		}

		result := &AdminResult{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return err
		}
		if result.Stat != "OK" {
			return fmt.Errorf("Could not find admin %s %s", result.Stat, *result.Message)
		}

		if result.Response.AdminID != rs.Primary.ID {
			return fmt.Errorf("Admin not found: %v - %v", rs.Primary.ID, result.Response.AdminID)
		}
		return nil
	}
}

func testAccCheckAdminConfig(rInt int) string {
	return fmt.Sprintf(`
resource "duo_admin" "test" {
  name = "test-%d"
  email = "le1f@wut.wut"
  phone = "+12813308004"	
}
`, rInt)
}

func testAccCheckAdminConfigUpdated(rInt int) string {
	return fmt.Sprintf(`
resource "duo_admin" "test" {
  name = "test-updated-%d"
  email = "le1f@wut.wut"
  phone = "+12813308004"	
}
`, rInt)
}
