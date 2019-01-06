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

func TestAccIntegration_Basic(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIntegrationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIntegrationConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationExists("duo_integration.test"),
					resource.TestCheckResourceAttr(
						"duo_integration.test", "name", fmt.Sprintf("test-integration-%d", rInt)),
					resource.TestCheckResourceAttr(
						"duo_integration.test", "type", "authapi"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIntegrationConfigUpdated(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationExists("duo_integration.test"),
					resource.TestCheckResourceAttr(
						"duo_integration.test", "name", fmt.Sprintf("test-updated-integration-%d", rInt)),
					resource.TestCheckResourceAttr(
						"duo_integration.test", "type", "authapi"),
				),
			},
		},
	})
}

func TestAccIntegration_import(t *testing.T) {
	resourceName := "duo_integration.test"
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAdminDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIntegrationConfig(rInt),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIntegrationDestroy(s *terraform.State) error {
	duoclient := testAccProvider.Meta().(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)
	for _, r := range s.RootModule().Resources {
		if r.Type != "duo_integration" {
			continue
		}

		_, body, err := duoAdminClient.SignedCall("GET", fmt.Sprintf("/admin/v1/integrations/%s", r.Primary.ID), nil, duoapi.UseTimeout)
		if err != nil {
			return err
		}

		result := &IntegrationResult{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return err
		}

		if result.Stat == "OK" {
			return fmt.Errorf("Found integration when it should have been deleted: %s", result.Response)
		}
	}
	return nil
}

func testAccCheckIntegrationExists(n string) resource.TestCheckFunc {
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

		_, body, err := duoAdminClient.SignedCall("GET", fmt.Sprintf("/admin/v1/integrations/%s", rs.Primary.ID), nil, duoapi.UseTimeout)
		if err != nil {
			return err
		}

		result := &IntegrationResult{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return err
		}

		if result.Stat != "OK" {
			return fmt.Errorf("Could not find integration %s %s", result.Stat, *result.Message)
		}

		if result.Response.IKey != rs.Primary.ID {
			return fmt.Errorf("Integration not found: %v - %v", rs.Primary.ID, result.Response.IKey)
		}
		return nil
	}
}

func testAccCheckIntegrationConfig(rInt int) string {
	return fmt.Sprintf(`
resource "duo_integration" "test" {
  name = "test-integration-%d"
  type = "authapi"
}
`, rInt)
}

func testAccCheckIntegrationConfigUpdated(rInt int) string {
	return fmt.Sprintf(`
resource "duo_integration" "test" {
  name = "test-updated-integration-%d"
  type = "authapi"
}
`, rInt)
}
