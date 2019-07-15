package azurerm

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMRoleDefinition_basic(t *testing.T) {
	dataSourceName := "data.azurerm_role_definition.test"

	id := uuid.New().String()
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRoleDefinition_basic(id, ri),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "type"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.0", "*"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.#", "3"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.0", "Microsoft.Authorization/*/Delete"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.1", "Microsoft.Authorization/*/Write"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.2", "Microsoft.Authorization/elevateAccess/Action"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMRoleDefinition_basicByName(t *testing.T) {
	dataSourceName := "data.azurerm_role_definition.test"

	id := uuid.New().String()
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRoleDefinition_byName(id, ri),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "type"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.0", "*"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.#", "3"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.0", "Microsoft.Authorization/*/Delete"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.1", "Microsoft.Authorization/*/Write"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.2", "Microsoft.Authorization/elevateAccess/Action"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMRoleDefinition_builtIn_contributor(t *testing.T) {
	dataSourceName := "data.azurerm_role_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMRoleDefinition_builtIn("Contributor"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "/providers/Microsoft.Authorization/roleDefinitions/b24988ac-6180-42a0-ab88-20f7382dd24c"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "type"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.0", "*"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.#", "5"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.0", "Microsoft.Authorization/*/Delete"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.1", "Microsoft.Authorization/*/Write"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.2", "Microsoft.Authorization/elevateAccess/Action"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.3", "Microsoft.Blueprint/blueprintAssignments/write"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.4", "Microsoft.Blueprint/blueprintAssignments/delete"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMRoleDefinition_builtIn_owner(t *testing.T) {
	dataSourceName := "data.azurerm_role_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMRoleDefinition_builtIn("Owner"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "/providers/Microsoft.Authorization/roleDefinitions/8e3af657-a8ff-443c-a75c-2fe8c4bcb635"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "type"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.0", "*"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMRoleDefinition_builtIn_reader(t *testing.T) {
	dataSourceName := "data.azurerm_role_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMRoleDefinition_builtIn("Reader"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "type"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.0", "*/read"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMRoleDefinition_builtIn_virtualMachineContributor(t *testing.T) {
	dataSourceName := "data.azurerm_role_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMRoleDefinition_builtIn("VirtualMachineContributor"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "/providers/Microsoft.Authorization/roleDefinitions/9980e02c-c2be-4d73-94e8-173b1dc7cf3c"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "type"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.#", "38"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.0", "Microsoft.Authorization/*/read"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.15", "Microsoft.Network/networkSecurityGroups/join/action"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMRoleDefinition_builtIn(name string) string {
	return fmt.Sprintf(`
data "azurerm_role_definition" "test" {
  name = "%s"
}
`, name)
}

func testAccDataSourceRoleDefinition_basic(id string, rInt int) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

resource "azurerm_role_definition" "test" {
  role_definition_id = "%s"
  name               = "acctestrd-%d"
  scope              = "${data.azurerm_subscription.primary.id}"
  description        = "Created by the Data Source Role Definition Acceptance Test"

  permissions {
    actions = ["*"]

    not_actions = [
      "Microsoft.Authorization/*/Delete",
      "Microsoft.Authorization/*/Write",
      "Microsoft.Authorization/elevateAccess/Action",
    ]
  }

  assignable_scopes = [
    "${data.azurerm_subscription.primary.id}",
  ]
}

data "azurerm_role_definition" "test" {
  role_definition_id = "${azurerm_role_definition.test.role_definition_id}"
  scope              = "${data.azurerm_subscription.primary.id}"
}
`, id, rInt)
}

func testAccDataSourceRoleDefinition_byName(id string, rInt int) string {
	return fmt.Sprintf(`
%s

data "azurerm_role_definition" "byName" {
  name  = "${azurerm_role_definition.test.name}"
  scope = "${data.azurerm_subscription.primary.id}"
}
`, testAccDataSourceRoleDefinition_basic(id, rInt))
}
