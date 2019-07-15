---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_workspace_linked_service"
sidebar_current: "docs-azurerm-log-analytics-workspace-linked-service"
description: |-
  Manages a Log Analytics (formally Operational Insights) Linked Service.
---

# azurerm_log_analytics_workspace_linked_service

Links a Log Analytics (formally Operational Insights) Workspace to another resource. The (currently) only linkable service is an Azure Automation Account.

~> **NOTE:** This resource has been deprecated in favour of the `azurerm_log_analytics_linked_service` resource and will be removed in the next major version of the AzureRM Provider. The new resource shares the same fields as this one, and information on migrating across [can be found in this guide](../guides/migrating-between-renamed-resources.html).


## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "resourcegroup-01"
  location = "West Europe"
}

resource "azurerm_automation_account" "test" {
  name                = "automation-01"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name = "Basic"
  }

  tags = {
    environment = "development"
  }
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "workspace-01"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_log_analytics_workspace_linked_service" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  workspace_name      = "${azurerm_log_analytics_workspace.test.name}"
  resource_id         = "${azurerm_automation_account.test.id}"
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group in which the Log Analytics Linked Service is created. Changing this forces a new resource to be created.

* `workspace_name` - (Required) Name of the Log Analytics Workspace that will contain the linkedServices resource. Changing this forces a new resource to be created.

* `linked_service_name` - (Optional) Name of the type of linkedServices resource to connect to the Log Analytics Workspace specified in `workspace_name`. Currently it defaults to and only supports `automation` as a value. Changing this forces a new resource to be created.

* `resource_id` - (Optional) The ID of the Resource that will be linked to the workspace.

* `linked_service_properties` - (Optional **Deprecated**) A `linked_service_properties` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

`linked_service_properties` supports the following:

* `resource_id` - (Optional  **Deprecated**) The resource id of the resource that will be linked to the workspace. This field has been deprecated in favour of the top-level `resource_id` field and will be removed in v2.0 of the AzureRM Provider.

## Attributes Reference

The following attributes are exported:

* `id` - The Log Analytics Linked Service ID.

* `name` - The automatically generated name of the Linked Service. This cannot be specified. The format is always `<workspace_name>/<linked_service_name>` e.g. `workspace1/Automation`

## Import

Log Analytics Workspaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_workspace_linked_service.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/linkedservices/automation
```
