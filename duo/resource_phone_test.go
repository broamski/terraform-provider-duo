package duo

import (
	"fmt"
	"testing"

	"github.com/duosecurity/duo_api_golang"
	admin "github.com/duosecurity/duo_api_golang/admin"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPhone_Basic(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPhoneDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckPhoneConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPhoneExists("duo_phone.test"),
					resource.TestCheckResourceAttr(
						"duo_phone.test", "name", fmt.Sprintf("test-phone-%d", rInt)),
					resource.TestCheckResourceAttr(
						"duo_phone.test", "number", "+18005551234"),
					resource.TestCheckResourceAttr(
						"duo_phone.test", "type", "Mobile"),
					resource.TestCheckResourceAttr(
						"duo_phone.test", "platform", "Unknown"),
				),
			},
			resource.TestStep{
				Config: testAccCheckPhoneConfigUpdated(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPhoneExists("duo_phone.test"),
					resource.TestCheckResourceAttr(
						"duo_phone.test", "name", fmt.Sprintf("test-updated-phone-%d", rInt)),
					resource.TestCheckResourceAttr(
						"duo_phone.test", "number", "+18005551235"),
					resource.TestCheckResourceAttr(
						"duo_phone.test", "type", "Mobile"),
					resource.TestCheckResourceAttr(
						"duo_phone.test", "platform", "Apple iOS"),
				),
			},
		},
	})
}

func TestAccPhone_import(t *testing.T) {
	resourceName := "duo_phone.test"
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPhoneDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckPhoneConfig(rInt),
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPhoneDestroy(s *terraform.State) error {
	duoclient := testAccProvider.Meta().(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	for _, r := range s.RootModule().Resources {
		if r.Type != "duo_phone" {
			continue
		}

		result, err := duoAdminClient.GetPhone(r.Primary.ID)
		if err != nil {
			return err
		}

		if result.Stat == "OK" {
			return fmt.Errorf("Found phone when it should have been deleted: %+v", result.Response)
		}
	}
	return nil
}

func testAccCheckPhoneExists(n string) resource.TestCheckFunc {
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

		result, err := duoAdminClient.GetPhone(rs.Primary.ID)
		if err != nil {
			return err
		}
		if result.Stat != "OK" {
			return fmt.Errorf("Could not find phone %s %s", result.Stat, *result.Message)
		}

		if result.Response.PhoneID != rs.Primary.ID {
			return fmt.Errorf("Phone not found: %v - %v", rs.Primary.ID, result.Response.PhoneID)
		}
		return nil
	}
}

func testAccCheckPhoneConfig(rInt int) string {
	return fmt.Sprintf(`
resource "duo_phone" "test" {
  name = "test-phone-%d"
  number = "+18005551234"
  type = "Mobile"
  platform = "Unknown"
}
`, rInt)
}

func testAccCheckPhoneConfigUpdated(rInt int) string {
	return fmt.Sprintf(`
resource "duo_phone" "test" {
  name = "test-updated-phone-%d"
  number = "+18005551235"
  type = "Mobile"
  platform = "Apple iOS"
}
`, rInt)
}
