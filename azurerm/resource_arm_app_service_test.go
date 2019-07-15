package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAzureRMAppServiceName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "ab",
			ErrCount: 0,
		},
		{
			Value:    "abc",
			ErrCount: 0,
		},
		{
			Value:    "webapp1",
			ErrCount: 0,
		},
		{
			Value:    "hello-world",
			ErrCount: 0,
		},
		{
			Value:    "hello_world",
			ErrCount: 1,
		},
		{
			Value:    "helloworld21!",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateAppServiceName(tc.Value, "azurerm_app_service")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the App Service Name to trigger a validation error for '%s'", tc.Value)
		}
	}
}

func TestAccAzureRMAppService_basic(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "outbound_ip_addresses"),
					resource.TestCheckResourceAttrSet(resourceName, "possible_outbound_ip_addresses"),
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

func TestAccAzureRMAppService_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMAppService_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_app_service"),
			},
		},
	})
}

func TestAccAzureRMAppService_movingAppService(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMAppService_moved(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_freeTier(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_freeTier(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
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

func TestAccAzureRMAppService_sharedTier(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_sharedTier(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
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

func TestAccAzureRMAppService_32Bit(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_32Bit(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.use_32_bit_worker_process", "true"),
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

func TestAccAzureRMAppService_http2Enabled(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_http2Enabled(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.http2_enabled", "true"),
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

func TestAccAzureRMAppService_alwaysOn(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_alwaysOn(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.always_on", "true"),
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

func TestAccAzureRMAppService_appCommandLine(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_appCommandLine(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.app_command_line", "/sbin/myserver -b 0.0.0.0"),
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

func TestAccAzureRMAppService_httpsOnly(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_httpsOnly(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "https_only", "true"),
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

func TestAccAzureRMAppService_clientCertEnabled(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	configClientCertEnabled := testAccAzureRMAppService_clientCertEnabled(ri, testLocation())
	configClientCertEnabledNotSet := testAccAzureRMAppService_clientCertEnabledNotSet(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: configClientCertEnabled,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "client_cert_enabled", "true"),
				),
			},
			{
				Config: configClientCertEnabledNotSet,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "client_cert_enabled", "false"),
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

func TestAccAzureRMAppService_appSettings(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_appSettings(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "app_settings.foo", "bar"),
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

func TestAccAzureRMAppService_clientAffinityEnabled(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_clientAffinityEnabled(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "client_affinity_enabled", "true"),
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

func TestAccAzureRMAppService_clientAffinityDisabled(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_clientAffinityDisabled(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "client_affinity_enabled", "false"),
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

func TestAccAzureRMAppService_virtualNetwork(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_virtualNetwork(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.virtual_network_name", fmt.Sprintf("acctestvn-%d", ri)),
				),
			},
			{
				Config: testAccAzureRMAppService_virtualNetworkUpdated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.virtual_network_name", fmt.Sprintf("acctestvn2-%d", ri)),
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

func TestAccAzureRMAppService_enableManageServiceIdentity(t *testing.T) {

	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_mangedServiceIdentity(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(resourceName, "identity.0.principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(resourceName, "identity.0.tenant_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_updateResourceByEnablingManageServiceIdentity(t *testing.T) {

	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()

	basicResourceNoManagedIdentity := testAccAzureRMAppService_basic(ri, testLocation())
	managedIdentityEnabled := testAccAzureRMAppService_mangedServiceIdentity(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: basicResourceNoManagedIdentity,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity.#", "0"),
				),
			},
			{
				Config: managedIdentityEnabled,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(resourceName, "identity.0.principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(resourceName, "identity.0.tenant_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_userAssignedIdentity(t *testing.T) {

	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_userAssignedIdentity(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity.0.type", "UserAssigned"),
					resource.TestCheckResourceAttr(resourceName, "identity.0.identity_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "identity.0.principal_id", ""),
					resource.TestCheckResourceAttr(resourceName, "identity.0.tenant_id", ""),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_multipleAssignedIdentities(t *testing.T) {

	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_multipleAssignedIdentities(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity.0.type", "SystemAssigned, UserAssigned"),
					resource.TestCheckResourceAttr(resourceName, "identity.0.identity_ids.#", "1"),
					resource.TestMatchResourceAttr(resourceName, "identity.0.principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(resourceName, "identity.0.tenant_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_clientAffinityUpdate(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_clientAffinity(ri, testLocation(), true)
	updatedConfig := testAccAzureRMAppService_clientAffinity(ri, testLocation(), false)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "client_affinity_enabled", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "client_affinity_enabled", "false"),
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

func TestAccAzureRMAppService_connectionStrings(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_connectionStrings(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "connection_string.3173438943.name", "First"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.3173438943.value", "first-connection-string"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.3173438943.type", "Custom"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.2442860602.name", "Second"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.2442860602.value", "some-postgresql-connection-string"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.2442860602.type", "PostgreSQL"),
				),
			},
			{
				Config: testAccAzureRMAppService_connectionStringsUpdated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "connection_string.3173438943.name", "First"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.3173438943.value", "first-connection-string"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.3173438943.type", "Custom"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.2442860602.name", "Second"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.2442860602.value", "some-postgresql-connection-string"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.2442860602.type", "PostgreSQL"),
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

func TestAccAzureRMAppService_storageAccounts(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_storageAccounts(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "storage_account.#", "1"),
				),
			},
			{
				Config: testAccAzureRMAppService_storageAccountsUpdated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "storage_account.#", "2"),
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

func TestAccAzureRMAppService_oneIpRestriction(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_oneIpRestriction(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.0.ip_address", "10.10.10.10"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.0.subnet_mask", "255.255.255.255"),
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

func TestAccAzureRMAppService_zeroedIpRestriction(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_oneIpRestriction(ri, testLocation())
	noBlocksConfig := testAccAzureRMAppService_basic(ri, testLocation())
	blocksEmptyConfig := testAccAzureRMAppService_zeroedIpRestriction(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				// This configuration includes a single explicit ip_restriction
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.#", "1"),
				),
			},
			{
				// This configuration has no site_config blocks at all.
				Config: noBlocksConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.#", "1"),
				),
			},
			{
				// This configuration explicitly sets ip_restriction to [] using attribute syntax.
				Config: blocksEmptyConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_manyIpRestrictions(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_manyIpRestrictions(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.0.ip_address", "10.10.10.10"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.0.subnet_mask", "255.255.255.255"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.1.ip_address", "20.20.20.0"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.1.subnet_mask", "255.255.255.0"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.2.ip_address", "30.30.0.0"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.2.subnet_mask", "255.255.0.0"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.3.ip_address", "192.168.1.2"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.3.subnet_mask", "255.255.255.0"),
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

func TestAccAzureRMAppService_defaultDocuments(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_defaultDocuments(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.default_documents.0", "first.html"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.default_documents.1", "second.jsp"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.default_documents.2", "third.aspx"),
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

func TestAccAzureRMAppService_enabled(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_enabled(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
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

func TestAccAzureRMAppService_localMySql(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_localMySql(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.local_mysql_enabled", "true"),
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

func TestAccAzureRMAppService_applicationBlobStorageLogs(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_applicationBlobStorageLogs(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "logs.0.application_logs.0.azure_blob_storage.0.level", "Information"),
					resource.TestCheckResourceAttr(resourceName, "logs.0.application_logs.0.azure_blob_storage.0.sas_url", "http://x.com/"),
					resource.TestCheckResourceAttr(resourceName, "logs.0.application_logs.0.azure_blob_storage.0.retention_in_days", "3"),
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

func TestAccAzureRMAppService_managedPipelineMode(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_managedPipelineMode(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.managed_pipeline_mode", "Classic"),
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

func TestAccAzureRMAppService_tagsUpdate(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_tags(ri, testLocation())
	updatedConfig := testAccAzureRMAppService_tagsUpdated(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Hello", "World"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.Hello", "World"),
					resource.TestCheckResourceAttr(resourceName, "tags.Terraform", "AcceptanceTests"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_remoteDebugging(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_remoteDebugging(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.remote_debugging_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.remote_debugging_version", "VS2015"),
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

func TestAccAzureRMAppService_windowsDotNet2(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_windowsDotNet(ri, testLocation(), "v2.0")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.dotnet_framework_version", "v2.0"),
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

func TestAccAzureRMAppService_windowsDotNet4(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_windowsDotNet(ri, testLocation(), "v4.0")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.dotnet_framework_version", "v4.0"),
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

func TestAccAzureRMAppService_windowsDotNetUpdate(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_windowsDotNet(ri, testLocation(), "v2.0")
	updatedConfig := testAccAzureRMAppService_windowsDotNet(ri, testLocation(), "v4.0")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.dotnet_framework_version", "v2.0"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.dotnet_framework_version", "v4.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_windowsJava7Jetty(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_windowsJava(ri, testLocation(), "1.7", "JETTY", "9.3")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_version", "1.7"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container", "JETTY"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container_version", "9.3"),
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

func TestAccAzureRMAppService_windowsJava8Jetty(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_windowsJava(ri, testLocation(), "1.8", "JETTY", "9.3")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_version", "1.8"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container", "JETTY"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container_version", "9.3"),
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
func TestAccAzureRMAppService_windowsJava11Jetty(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_windowsJava(ri, testLocation(), "11", "JETTY", "9.3")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_version", "11"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container", "JETTY"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container_version", "9.3"),
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
func TestAccAzureRMAppService_windowsJava7Tomcat(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_windowsJava(ri, testLocation(), "1.7", "TOMCAT", "9.0")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_version", "1.7"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container", "TOMCAT"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container_version", "9.0"),
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

func TestAccAzureRMAppService_windowsJava8Tomcat(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_windowsJava(ri, testLocation(), "1.8", "TOMCAT", "9.0")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_version", "1.8"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container", "TOMCAT"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container_version", "9.0"),
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
func TestAccAzureRMAppService_windowsJava11Tomcat(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_windowsJava(ri, testLocation(), "11", "TOMCAT", "9.0")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_version", "11"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container", "TOMCAT"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container_version", "9.0"),
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

func TestAccAzureRMAppService_windowsPHP7(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_windowsPHP(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.php_version", "7.2"),
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

func TestAccAzureRMAppService_windowsPython(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_windowsPython(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.python_version", "3.4"),
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

func TestAccAzureRMAppService_webSockets(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_webSockets(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.websockets_enabled", "true"),
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

func TestAccAzureRMAppService_scmType(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_scmType(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.scm_type", "LocalGit"),
					resource.TestCheckResourceAttr(resourceName, "source_control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "site_credential.#", "1"),
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

func TestAccAzureRMAppService_ftpsState(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_ftpsState(ri, testLocation())
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ftps_state", "AllAllowed"),
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

func TestAccAzureRMAppService_linuxFxVersion(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_linuxFxVersion(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.always_on", "true"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.linux_fx_version", "DOCKER|(golang:latest)"),
					resource.TestCheckResourceAttr(resourceName, "app_settings.WEBSITES_ENABLE_APP_SERVICE_STORAGE", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_minTls(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_minTls(ri, testLocation(), "1.0")
	updatedConfig := testAccAzureRMAppService_minTls(ri, testLocation(), "1.1")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.min_tls_version", "1.0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.min_tls_version", "1.1"),
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

func TestAccAzureRMAppService_corsSettings(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_corsSettings(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.cors.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.cors.0.support_credentials", "true"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.cors.0.allowed_origins.#", "3"),
				)},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppService_authSettingsAdditionalLoginParams(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	tenantID := os.Getenv("ARM_TENANT_ID")
	config := testAccAzureRMAppService_authSettingsAdditionalLoginParams(ri, testLocation(), tenantID)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.additional_login_params.test_key", "test_value"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
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

func TestAccAzureRMAppService_authSettingsAdditionalAllowedExternalRedirectUrls(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	tenantID := os.Getenv("ARM_TENANT_ID")
	config := testAccAzureRMAppService_authSettingsAdditionalAllowedExternalRedirectUrls(ri, testLocation(), tenantID)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.allowed_external_redirect_urls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.allowed_external_redirect_urls.0", "https://terra.form"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
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

func TestAccAzureRMAppService_authSettingsRuntimeVersion(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	tenantID := os.Getenv("ARM_TENANT_ID")
	config := testAccAzureRMAppService_authSettingsRuntimeVersion(ri, testLocation(), tenantID)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.runtime_version", "1.0"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
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

func TestAccAzureRMAppService_authSettingsTokenRefreshExtensionHours(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	tenantID := os.Getenv("ARM_TENANT_ID")
	config := testAccAzureRMAppService_authSettingsTokenRefreshExtensionHours(ri, testLocation(), tenantID)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.token_refresh_extension_hours", "75"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
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

func TestAccAzureRMAppService_authSettingsUnauthenticatedClientAction(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	tenantID := os.Getenv("ARM_TENANT_ID")
	config := testAccAzureRMAppService_authSettingsUnauthenticatedClientAction(ri, testLocation(), tenantID)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.unauthenticated_client_action", "RedirectToLoginPage"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
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

func TestAccAzureRMAppService_authSettingsTokenStoreEnabled(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	tenantID := os.Getenv("ARM_TENANT_ID")
	config := testAccAzureRMAppService_authSettingsTokenStoreEnabled(ri, testLocation(), tenantID)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.token_store_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
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

func TestAccAzureRMAppService_aadAuthSettings(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	tenantID := os.Getenv("ARM_TENANT_ID")
	config := testAccAzureRMAppService_aadAuthSettings(ri, testLocation(), tenantID)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
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

func TestAccAzureRMAppService_facebookAuthSettings(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_facebookAuthSettings(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.facebook.0.app_id", "facebookappid"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.facebook.0.app_secret", "facebookappsecret"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.facebook.0.oauth_scopes.#", "1"),
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

func TestAccAzureRMAppService_googleAuthSettings(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_googleAuthSettings(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.google.0.client_id", "googleclientid"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.google.0.client_secret", "googleclientsecret"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.google.0.oauth_scopes.#", "1"),
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

func TestAccAzureRMAppService_microsoftAuthSettings(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_microsoftAuthSettings(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.microsoft.0.client_id", "microsoftclientid"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.microsoft.0.client_secret", "microsoftclientsecret"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.microsoft.0.oauth_scopes.#", "1"),
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

func TestAccAzureRMAppService_twitterAuthSettings(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_twitterAuthSettings(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.twitter.0.consumer_key", "twitterconsumerkey"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.twitter.0.consumer_secret", "twitterconsumersecret"),
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

func TestAccAzureRMAppService_multiAuthSettings(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	tenantID := os.Getenv("ARM_TENANT_ID")
	config1 := testAccAzureRMAppService_aadAuthSettings(ri, testLocation(), tenantID)
	config2 := testAccAzureRMAppService_aadMicrosoftAuthSettings(ri, testLocation(), tenantID)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.microsoft.0.client_id", "microsoftclientid"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.microsoft.0.client_secret", "microsoftclientsecret"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.microsoft.0.oauth_scopes.#", "1"),
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

func TestAccAzureRMAppService_basicWindowsContainer(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppService_basicWindowsContainer(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.windows_fx_version", "DOCKER|mcr.microsoft.com/azure-app-service/samples/aspnethelloworld:latest"),
					resource.TestCheckResourceAttr(resourceName, "app_settings.DOCKER_REGISTRY_SERVER_URL", "https://mcr.microsoft.com"),
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

func testCheckAzureRMAppServiceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).appServicesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return nil
	}

	return nil
}

func testCheckAzureRMAppServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		appServiceName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Service: %s", appServiceName)
		}

		client := testAccProvider.Meta().(*ArmClient).appServicesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, appServiceName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: App Service %q (resource group: %q) does not exist", appServiceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on appServicesClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMAppService_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_requiresImport(rInt int, location string) string {
	template := testAccAzureRMAppService_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service" "import" {
  name                = "${azurerm_app_service.test.name}"
  location            = "${azurerm_app_service.test.location}"
  resource_group_name = "${azurerm_app_service.test.resource_group_name}"
  app_service_plan_id = "${azurerm_app_service.test.app_service_plan_id}"
}
`, template)
}

func testAccAzureRMAppService_freeTier(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Free"
    size = "F1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    use_32_bit_worker_process = true
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_moved(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service_plan" "other" {
  name                = "acctestASP2-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.other.id}"
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppService_sharedTier(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Free"
    size = "F1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    use_32_bit_worker_process = true
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_alwaysOn(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    always_on = true
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_appCommandLine(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    app_command_line = "/sbin/myserver -b 0.0.0.0"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_httpsOnly(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  https_only          = true
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_clientCertEnabled(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  client_cert_enabled = true
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_clientCertEnabledNotSet(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_32Bit(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    use_32_bit_worker_process = true
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_http2Enabled(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    http2_enabled = true
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_appSettings(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  app_settings = {
    "foo" = "bar"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_clientAffinityEnabled(rInt int, location string) string {
	return testAccAzureRMAppService_clientAffinity(rInt, location, true)
}

func testAccAzureRMAppService_clientAffinityDisabled(rInt int, location string) string {
	return testAccAzureRMAppService_clientAffinity(rInt, location, false)
}

func testAccAzureRMAppService_clientAffinity(rInt int, location string, clientAffinity bool) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                    = "acctestAS-%d"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  app_service_plan_id     = "${azurerm_app_service_plan.test.id}"
  client_affinity_enabled = %t
}
`, rInt, location, rInt, rInt, clientAffinity)
}

func testAccAzureRMAppService_virtualNetwork(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  subnet {
    name           = "internal"
    address_prefix = "10.0.1.0/24"
  }
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    virtual_network_name = "${azurerm_virtual_network.test.name}"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppService_virtualNetworkUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  subnet {
    name           = "internal"
    address_prefix = "10.0.1.0/24"
  }
}

resource "azurerm_virtual_network" "second" {
  name                = "acctestvn2-%d"
  address_space       = ["172.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  subnet {
    name           = "internal"
    address_prefix = "172.0.1.0/24"
  }
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    virtual_network_name = "${azurerm_virtual_network.second.name}"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMAppService_mangedServiceIdentity(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  identity {
    type = "SystemAssigned"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_userAssignedIdentity(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  identity {
    type         = "UserAssigned"
    identity_ids = ["${azurerm_user_assigned_identity.test.id}"]
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppService_multipleAssignedIdentities(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = ["${azurerm_user_assigned_identity.test.id}"]
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppService_connectionStrings(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  connection_string {
    name  = "First"
    value = "first-connection-string"
    type  = "Custom"
  }

  connection_string {
    name  = "Second"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_connectionStringsUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  connection_string {
    name  = "Second"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  connection_string {
    name  = "First"
    value = "first-connection-string"
    type  = "Custom"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_storageAccounts(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acct%d"
  location                 = "${azurerm_resource_group.test.location}"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestcontainer"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  storage_account {
    name         = "blobs"
    type         = "AzureBlob"
    account_name = "${azurerm_storage_account.test.name}"
    share_name   = "${azurerm_storage_container.test.name}"
    access_key   = "${azurerm_storage_account.test.primary_access_key}"
    mount_path   = "/blobs"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppService_storageAccountsUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acct%d"
  location                 = "${azurerm_resource_group.test.location}"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestcontainer"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
}

resource "azurerm_storage_share" "test" {
  name                 = "acctestshare"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  storage_account {
    name         = "blobs"
    type         = "AzureBlob"
    account_name = "${azurerm_storage_account.test.name}"
    share_name   = "${azurerm_storage_container.test.name}"
    access_key   = "${azurerm_storage_account.test.primary_access_key}"
    mount_path   = "/blobs"
  }

  storage_account {
    name         = "files"
    type         = "AzureFiles"
    account_name = "${azurerm_storage_account.test.name}"
    share_name   = "${azurerm_storage_share.test.name}"
    access_key   = "${azurerm_storage_account.test.primary_access_key}"
    mount_path   = "/files"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppService_oneIpRestriction(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10"
    }
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_manyIpRestrictions(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10"
    }

    ip_restriction {
      ip_address  = "20.20.20.0"
      subnet_mask = "255.255.255.0"
    }

    ip_restriction {
      ip_address  = "30.30.0.0"
      subnet_mask = "255.255.0.0"
    }

    ip_restriction {
      ip_address  = "192.168.1.2"
      subnet_mask = "255.255.255.0"
    }
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_zeroedIpRestriction(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    ip_restriction = []
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_defaultDocuments(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    default_documents = [
      "first.html",
      "second.jsp",
      "third.aspx",
    ]
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_enabled(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  enabled             = false
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_localMySql(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    local_mysql_enabled = true
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_applicationBlobStorageLogs(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  logs {
    application_logs {
      azure_blob_storage {
        level             = "Information"
        sas_url           = "http://x.com/"
        retention_in_days = 3
      }
    }
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_managedPipelineMode(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    managed_pipeline_mode = "Classic"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_remoteDebugging(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    remote_debugging_enabled = true
    remote_debugging_version = "VS2015"
  }

  tags = {
    Hello = "World"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_tags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  tags = {
    Hello = "World"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_tagsUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  tags = {
    "Hello"     = "World"
    "Terraform" = "AcceptanceTests"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_windowsDotNet(rInt int, location, version string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    dotnet_framework_version = "%s"
  }
}
`, rInt, location, rInt, rInt, version)
}

func testAccAzureRMAppService_windowsJava(rInt int, location, javaVersion, container, containerVersion string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    java_version           = "%s"
    java_container         = "%s"
    java_container_version = "%s"
  }
}
`, rInt, location, rInt, rInt, javaVersion, container, containerVersion)
}

func testAccAzureRMAppService_windowsPHP(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    php_version = "7.2"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_windowsPython(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    python_version = "3.4"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_webSockets(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    websockets_enabled = true
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_scmType(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    scm_type = "LocalGit"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_ftpsState(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    ftps_state = "AllAllowed"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_linuxFxVersion(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    always_on        = true
    linux_fx_version = "DOCKER|(golang:latest)"
  }

  app_settings = {
    "WEBSITES_ENABLE_APP_SERVICE_STORAGE" = "false"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_minTls(rInt int, location string, tlsVersion string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    min_tls_version = "%s"
  }
}
`, rInt, location, rInt, rInt, tlsVersion)
}

func testAccAzureRMAppService_corsSettings(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    cors {
      allowed_origins = [
        "http://www.contoso.com",
        "www.contoso.com",
        "contoso.com",
      ]

      support_credentials = true
    }
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_authSettingsAdditionalLoginParams(rInt int, location string, tenantID string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    always_on = true
  }

  auth_settings {
    enabled = true
    issuer  = "https://sts.windows.net/%s"

    additional_login_params = {
      test_key = "test_value"
    }

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }
  }
}
`, rInt, location, rInt, rInt, tenantID)
}

func testAccAzureRMAppService_authSettingsAdditionalAllowedExternalRedirectUrls(rInt int, location string, tenantID string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    always_on = true
  }

  auth_settings {
    enabled = true
    issuer  = "https://sts.windows.net/%s"

    allowed_external_redirect_urls = [
      "https://terra.form",
    ]

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }
  }
}
`, rInt, location, rInt, rInt, tenantID)
}

func testAccAzureRMAppService_authSettingsRuntimeVersion(rInt int, location string, tenantID string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    always_on = true
  }

  auth_settings {
    enabled         = true
    issuer          = "https://sts.windows.net/%s"
    runtime_version = "1.0"

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }
  }
}
`, rInt, location, rInt, rInt, tenantID)
}

func testAccAzureRMAppService_authSettingsTokenRefreshExtensionHours(rInt int, location string, tenantID string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    always_on = true
  }

  auth_settings {
    enabled                       = true
    issuer                        = "https://sts.windows.net/%s"
    token_refresh_extension_hours = 75

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }
  }
}
`, rInt, location, rInt, rInt, tenantID)
}

func testAccAzureRMAppService_authSettingsTokenStoreEnabled(rInt int, location string, tenantID string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    always_on = true
  }

  auth_settings {
    enabled             = true
    issuer              = "https://sts.windows.net/%s"
    token_store_enabled = true

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }
  }
}
`, rInt, location, rInt, rInt, tenantID)
}

func testAccAzureRMAppService_authSettingsUnauthenticatedClientAction(rInt int, location string, tenantID string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    always_on = true
  }

  auth_settings {
    enabled                       = true
    issuer                        = "https://sts.windows.net/%s"
    unauthenticated_client_action = "RedirectToLoginPage"

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }
  }
}
`, rInt, location, rInt, rInt, tenantID)
}

func testAccAzureRMAppService_aadAuthSettings(rInt int, location string, tenantID string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    always_on = true
  }

  auth_settings {
    enabled = true
    issuer  = "https://sts.windows.net/%s"

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }
  }
}
`, rInt, location, rInt, rInt, tenantID)
}

func testAccAzureRMAppService_facebookAuthSettings(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    always_on = true
  }

  auth_settings {
    enabled = true

    facebook {
      app_id     = "facebookappid"
      app_secret = "facebookappsecret"

      oauth_scopes = [
        "facebookscope",
      ]
    }
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_googleAuthSettings(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    always_on = true
  }

  auth_settings {
    enabled = true

    google {
      client_id     = "googleclientid"
      client_secret = "googleclientsecret"

      oauth_scopes = [
        "googlescope",
      ]
    }
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_microsoftAuthSettings(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    always_on = true
  }

  auth_settings {
    enabled = true

    microsoft {
      client_id     = "microsoftclientid"
      client_secret = "microsoftclientsecret"

      oauth_scopes = [
        "microsoftscope",
      ]
    }
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_twitterAuthSettings(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    always_on = true
  }

  auth_settings {
    enabled = true

    twitter {
      consumer_key    = "twitterconsumerkey"
      consumer_secret = "twitterconsumersecret"
    }
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_aadMicrosoftAuthSettings(rInt int, location string, tenantID string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    always_on = true
  }

  auth_settings {
    enabled          = true
    issuer           = "https://sts.windows.net/%s"
    default_provider = "%s"

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }

    microsoft {
      client_id     = "microsoftclientid"
      client_secret = "microsoftclientsecret"

      oauth_scopes = [
        "microsoftscope",
      ]
    }
  }
}
`, rInt, location, rInt, rInt, tenantID, web.AzureActiveDirectory)
}

func testAccAzureRMAppService_basicWindowsContainer(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  is_xenon            = true
  kind                = "xenon"

  sku {
    tier = "PremiumContainer"
    size = "PC2"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    windows_fx_version = "DOCKER|mcr.microsoft.com/azure-app-service/samples/aspnethelloworld:latest"
  }

  app_settings = {
    "DOCKER_REGISTRY_SERVER_URL"      = "https://mcr.microsoft.com"
    "DOCKER_REGISTRY_SERVER_USERNAME" = ""
    "DOCKER_REGISTRY_SERVER_PASSWORD" = ""
  }
}
`, rInt, location, rInt, rInt)
}
