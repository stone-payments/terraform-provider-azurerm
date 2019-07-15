package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func testCheckAzureRMHDInsightClusterDestroy(terraformResourceName string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != terraformResourceName {
				continue
			}

			client := testAccProvider.Meta().(*ArmClient).hdinsight.ClustersClient
			ctx := testAccProvider.Meta().(*ArmClient).StopContext
			name := rs.Primary.Attributes["name"]
			resourceGroup := rs.Primary.Attributes["resource_group_name"]
			resp, err := client.Get(ctx, resourceGroup, name)

			if err != nil {
				if !utils.ResponseWasNotFound(resp.Response) {
					return err
				}
			}
		}

		return nil
	}
}

func testCheckAzureRMHDInsightClusterExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		clusterName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).hdinsight.ClustersClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, clusterName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: HDInsight Cluster %q (Resource Group: %q) does not exist", clusterName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on hdinsightClustersClient: %+v", err)
		}

		return nil
	}
}
