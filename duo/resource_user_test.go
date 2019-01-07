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

func TestAccUser_Basic(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckUserConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("duo_user.test"),
					resource.TestCheckResourceAttr(
						"duo_user.test", "username", fmt.Sprintf("test-user-%d", rInt)),
					resource.TestCheckResourceAttr(
						"duo_user.test", "alias1", "t1"),
					resource.TestCheckResourceAttr(
						"duo_user.test", "alias2", "t2"),
					resource.TestCheckResourceAttr(
						"duo_user.test", "alias3", "t3"),
					resource.TestCheckResourceAttr(
						"duo_user.test", "alias4", "t4"),
					resource.TestCheckResourceAttr(
						"duo_user.test", "realname", "Mister Sir"),
					resource.TestCheckResourceAttr(
						"duo_user.test", "email", "le1f@wut.what"),
					resource.TestCheckResourceAttr(
						"duo_user.test", "status", "active"),
				),
			},
			resource.TestStep{
				Config: testAccCheckUserConfigUpdated(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("duo_user.test"),
					resource.TestCheckResourceAttr(
						"duo_user.test", "username", fmt.Sprintf("test-user-updated-%d", rInt)),
					resource.TestCheckResourceAttr(
						"duo_user.test", "alias1", "t1"),
					resource.TestCheckResourceAttr(
						"duo_user.test", "alias2", "t2"),
					resource.TestCheckResourceAttr(
						"duo_user.test", "alias3", "t3"),
					resource.TestCheckResourceAttr(
						"duo_user.test", "alias4", "t4"),
					resource.TestCheckResourceAttr(
						"duo_user.test", "realname", "Mister Sir"),
					resource.TestCheckResourceAttr(
						"duo_user.test", "email", "le1f@wut.what"),
					resource.TestCheckResourceAttr(
						"duo_user.test", "status", "bypass"),
				),
			},
		},
	})
}

func TestAccUser_import(t *testing.T) {
	resourceName := "duo_user.test"
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckUserConfig(rInt),
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckUserDestroy(s *terraform.State) error {
	duoclient := testAccProvider.Meta().(*duoapi.DuoApi)
	duoAdminClient := admin.New(*duoclient)

	for _, r := range s.RootModule().Resources {
		if r.Type != "duo_user" {
			continue
		}

		result, err := duoAdminClient.GetUser(r.Primary.ID)
		if err != nil {
			return err
		}

		if result.Stat == "OK" {
			return fmt.Errorf("Found user when it should have been deleted: %+v", result.Response)
		}
	}
	return nil
}

func testAccCheckUserExists(n string) resource.TestCheckFunc {
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

		result, err := duoAdminClient.GetUser(rs.Primary.ID)
		if err != nil {
			return err
		}

		if result.Stat != "OK" {
			return fmt.Errorf("Could not find integration %s %s", result.Stat, *result.Message)
		}

		if result.Response.UserID != rs.Primary.ID {
			return fmt.Errorf("Integration not found: %v - %v", rs.Primary.ID, result.Response.UserID)
		}
		return nil
	}
}

func testAccCheckUserConfig(rInt int) string {
	return fmt.Sprintf(`
resource "duo_user" "test" {
  username = "test-user-%d"
  alias1 = "t1"
  alias2 = "t2"
  alias3 = "t3"
  alias4 = "t4"
  realname = "Mister Sir"
  email = "le1f@wut.what"
  status = "active"
}
`, rInt)
}

func testAccCheckUserConfigUpdated(rInt int) string {
	return fmt.Sprintf(`
resource "duo_user" "test" {
	username = "test-user-updated-%d"
	alias1 = "t1"
	alias2 = "t2"
	alias3 = "t3"
	alias4 = "t4"
	realname = "Mister Sir"
	email = "le1f@wut.what"
	status = "bypass"
}
`, rInt)
}
