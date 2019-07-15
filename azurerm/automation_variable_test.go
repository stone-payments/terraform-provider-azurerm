package azurerm

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestParseAzureRmAutomationVariableValue(t *testing.T) {
	type ExpectFunc func(interface{}) bool
	cases := []struct {
		Name        string
		Resource    string
		IsNil       bool
		Value       string
		HasError    bool
		ExpectValue interface{}
		Expect      ExpectFunc
	}{
		{
			Name:        "string variable",
			Resource:    "azurerm_automation_variable_string",
			Value:       "\"Test String\"",
			HasError:    false,
			ExpectValue: "Test String",
			Expect:      func(v interface{}) bool { return v.(string) == "Test String" },
		},
		{
			Name:        "integer variable",
			Resource:    "azurerm_automation_variable_int",
			Value:       "135",
			HasError:    false,
			ExpectValue: 135,
			Expect:      func(v interface{}) bool { return v.(int32) == 135 },
		},
		{
			Name:        "boolean variable",
			Resource:    "azurerm_automation_variable_bool",
			Value:       "true",
			HasError:    false,
			ExpectValue: true,
			Expect:      func(v interface{}) bool { return v.(bool) == true },
		},
		{
			Name:        "datetime variable",
			Resource:    "azurerm_automation_variable_datetime",
			Value:       "\"\\/Date(1556142054074)\\/\"",
			HasError:    false,
			ExpectValue: time.Date(2019, time.April, 24, 21, 40, 54, 74000000, time.UTC),
			Expect: func(v interface{}) bool {
				return v.(time.Time) == time.Date(2019, time.April, 24, 21, 40, 54, 74000000, time.UTC)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			value := &tc.Value
			if tc.IsNil {
				value = nil
			}
			actual, err := parseAzureAutomationVariableValue(tc.Resource, value)
			if tc.HasError && err == nil {
				t.Fatalf("Expect parseAzureAutomationVariableValue to return error for resource %q and value %s", tc.Resource, tc.Value)
			}
			if !tc.HasError {
				if err != nil {
					t.Fatalf("Expect parseAzureAutomationVariableValue to return no error for resource %q and value %s, err: %+v", tc.Resource, tc.Value, err)
				} else if !tc.Expect(actual) {
					t.Fatalf("Expect parseAzureAutomationVariableValue to return %v instead of %v for resource %q and value %s", tc.ExpectValue, actual, tc.Resource, tc.Value)
				}
			}
		})
	}
}

func testCheckAzureRMAutomationVariableExists(resourceName string, varType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Automation %s Variable not found: %s", varType, resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		accountName := rs.Primary.Attributes["automation_account_name"]

		client := testAccProvider.Meta().(*ArmClient).automation.VariableClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.Get(ctx, resourceGroup, accountName, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Automation %s Variable %q (Automation Account Name %q / Resource Group %q) does not exist", varType, name, accountName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on automationVariableClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMAutomationVariableDestroy(s *terraform.State, varType string) error {
	client := testAccProvider.Meta().(*ArmClient).automation.VariableClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	resourceName := fmt.Sprintf("azurerm_automation_variable_%s", strings.ToLower(varType))

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourceName {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		accountName := rs.Primary.Attributes["automation_account_name"]

		if resp, err := client.Get(ctx, resourceGroup, accountName, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on automationVariableClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}
