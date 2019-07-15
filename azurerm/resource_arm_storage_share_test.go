package azurerm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMStorageShare_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()
	resourceName := "azurerm_storage_share.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShare_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMStorageShare_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()
	resourceName := "azurerm_storage_share.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShare_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMStorageShare_requiresImport(ri, rs, location),
				ExpectError: testRequiresImportError("azurerm_storage_share"),
			},
		},
	})
}

func TestAccAzureRMStorageShare_disappears(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()
	resourceName := "azurerm_storage_share.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShare_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(resourceName),
					testCheckAzureRMStorageShareDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMStorageShare_metaData(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()
	resourceName := "azurerm_storage_share.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShare_metaData(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMStorageShare_metaDataUpdated(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMStorageShare_acl(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()
	resourceName := "azurerm_storage_share.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShare_acl(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMStorageShare_aclUpdated(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMStorageShare_updateQuota(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()
	resourceName := "azurerm_storage_share.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShare_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMStorageShare_updateQuota(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "quota", "5"),
				),
			},
		},
	})
}

func testCheckAzureRMStorageShareExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		shareName := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		storageClient := testAccProvider.Meta().(*ArmClient).storage
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error finding Resource Group: %s", err)
		}

		client, err := storageClient.FileSharesClient(ctx, *resourceGroup, accountName)
		if err != nil {
			return fmt.Errorf("Error building FileShare Client: %s", err)
		}

		_, err = client.GetProperties(ctx, accountName, shareName)
		if err != nil {
			return fmt.Errorf("Bad: Share %q (Storage Account: %q) does not exist", shareName, accountName)
		}

		return nil
	}
}

func testCheckAzureRMStorageShareDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		shareName := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		storageClient := testAccProvider.Meta().(*ArmClient).storage
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error finding Resource Group: %s", err)
		}

		client, err := storageClient.FileSharesClient(ctx, *resourceGroup, accountName)
		if err != nil {
			return fmt.Errorf("Error building FileShare Client: %s", err)
		}

		if _, err := client.Delete(ctx, accountName, shareName, true); err != nil {
			return fmt.Errorf("Error deleting Share %q (Account %q): %v", shareName, accountName, err)
		}

		return nil
	}
}

func testCheckAzureRMStorageShareDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_share" {
			continue
		}

		shareName := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		storageClient := testAccProvider.Meta().(*ArmClient).storage
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error finding Resource Group: %s", err)
		}
		if resourceGroup == nil {
			return nil
		}

		client, err := storageClient.FileSharesClient(ctx, *resourceGroup, accountName)
		if err != nil {
			return fmt.Errorf("Error building FileShare Client: %s", err)
		}

		props, err := client.GetProperties(ctx, accountName, shareName)
		if err != nil {
			if utils.ResponseWasNotFound(props.Response) {
				return nil
			}

			return fmt.Errorf("Error retrieving Share %q: %s", shareName, accountName)
		}

		return fmt.Errorf("Bad: Share %q (storage account: %q) still exists", shareName, accountName)
	}

	return nil
}

func testAccAzureRMStorageShare_basic(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageShare_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
}
`, template, rString)
}

func testAccAzureRMStorageShare_metaData(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageShare_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"

  metadata = {
    hello = "world"
  }
}
`, template, rString)
}

func testAccAzureRMStorageShare_metaDataUpdated(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageShare_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"

  metadata = {
    hello = "world"
    happy = "birthday"
  }
}
`, template, rString)
}

func testAccAzureRMStorageShare_acl(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageShare_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"

  acl {
    id = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "rwd"
      start       = "2019-07-02T09:38:21.0000000Z"
      expiry      = "2019-07-02T10:38:21.0000000Z"
    }
  }
}
`, template, rString)
}

func testAccAzureRMStorageShare_aclUpdated(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageShare_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"

  acl {
    id = "AAAANDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "rwd"
      start       = "2019-07-02T09:38:21.0000000Z"
      expiry      = "2019-07-02T10:38:21.0000000Z"
    }
  }
  acl {
    id = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "rwd"
      start       = "2019-07-02T09:38:21.0000000Z"
      expiry      = "2019-07-02T10:38:21.0000000Z"
    }
  }
}
`, template, rString)
}
func testAccAzureRMStorageShare_requiresImport(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageShare_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "import" {
  name                 = "${azurerm_storage_share.test.name}"
  resource_group_name  = "${azurerm_storage_share.test.resource_group_name}"
  storage_account_name = "${azurerm_storage_share.test.storage_account_name}"
}
`, template)
}

func testAccAzureRMStorageShare_updateQuota(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageShare_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
  quota                = 5
}
`, template, rString)
}

func testAccAzureRMStorageShare_template(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}
`, rInt, location, rString)
}

func TestValidateArmStorageShareName(t *testing.T) {
	validNames := []string{
		"valid-name",
		"valid02-name",
	}
	for _, v := range validNames {
		_, errors := validateArmStorageShareName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Share Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"InvalidName1",
		"-invalidname1",
		"invalid_name",
		"invalid!",
		"double-hyphen--invalid",
		"ww",
		strings.Repeat("w", 65),
	}
	for _, v := range invalidNames {
		_, errors := validateArmStorageShareName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Share Name", v)
		}
	}
}
