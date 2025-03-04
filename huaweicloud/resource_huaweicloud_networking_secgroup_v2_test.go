package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/security/groups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccNetworkingV2SecGroup_basic(t *testing.T) {
	var secGroup groups.SecGroup
	name := fmt.Sprintf("seg-acc-test-%s", acctest.RandString(5))
	updatedName := fmt.Sprintf("%s-updated", name)
	resourceName := "huaweicloud_networking_secgroup.secgroup_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(resourceName, &secGroup),
					testAccCheckNetworkingV2SecGroupRuleCount(&secGroup, 6),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "6"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSecGroup_update(updatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceName, "id", &secGroup.ID),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroup_withEpsId(t *testing.T) {
	var secGroup groups.SecGroup
	name := fmt.Sprintf("seg-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_networking_secgroup.secgroup_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecGroup_epsId(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(resourceName, &secGroup),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroup_noDefaultRules(t *testing.T) {
	var secGroup groups.SecGroup
	name := fmt.Sprintf("seg-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_networking_secgroup.secgroup_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecGroup_noDefaultRules(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(resourceName, &secGroup),
					testAccCheckNetworkingV2SecGroupRuleCount(&secGroup, 0),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "0"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2SecGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV2Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_networking_secgroup" {
			continue
		}

		_, err := groups.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Security group still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2SecGroupExists(n string, secGroup *groups.SecGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		networkingClient, err := config.NetworkingV2Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
		}

		found, err := groups.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Security group not found")
		}

		*secGroup = *found

		return nil
	}
}

func testAccCheckNetworkingV2SecGroupRuleCount(
	sg *groups.SecGroup, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(sg.Rules) == count {
			return nil
		}

		return fmtp.Errorf("Unexpected number of rules in group %s. Expected %d, got %d",
			sg.ID, count, len(sg.Rules))
	}
}

func testAccSecGroup_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name        = "%s"
  description = "security group acceptance test"
}
`, name)
}

func testAccSecGroup_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name        = "%s"
  description = "security group acceptance test updated"
}
`, name)
}

func testAccSecGroup_epsId(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name                  = "%s"
  description           = "ecurity group acceptance test with eps ID"
  enterprise_project_id = "%s"
}
`, name, HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccSecGroup_noDefaultRules(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name                 = "%s"
  description          = "security group acceptance test without default rules"
  delete_default_rules = true
}
`, name)
}
