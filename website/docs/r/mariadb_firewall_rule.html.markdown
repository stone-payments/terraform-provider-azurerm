---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mariadb_firewall_rule"
sidebar_current: "docs-azurerm-resource-database-mariadb-firewall-rule"
description: |-
  Manages a Firewall Rule for a MariaDB Server.
---

# azurerm_mariadb_firewall_rule

Manages a Firewall Rule for a MariaDB Server

## Example Usage (Single IP Address)

```hcl
resource "azurerm_mariadb_firewall_rule" "test" {
  name                = "test-rule"
  resource_group_name = "test-rg"
  server_name         = "test-server"
  start_ip_address    = "40.112.8.12"
  end_ip_address      = "40.112.8.12"
}
```

## Example Usage (IP Range)

```hcl
resource "azurerm_mariadb_firewall_rule" "test" {
  name                = "test-rule"
  resource_group_name = "test-rg"
  server_name         = "test-server"
  start_ip_address    = "40.112.0.0"
  end_ip_address      = "40.112.255.255"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the MariaDB Firewall Rule. Changing this forces a new resource to be created.

* `server_name` - (Required) Specifies the name of the MariaDB Server. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the MariaDB Server exists. Changing this forces a new resource to be created.

* `start_ip_address` - (Required) Specifies the Start IP Address associated with this Firewall Rule. Changing this forces a new resource to be created.

* `end_ip_address` - (Required) Specifies the End IP Address associated with this Firewall Rule. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MariaDB Firewall Rule.

## Import

MariaDB Firewall rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mariadb_firewall_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforMariaDB/servers/server1/firewallRules/rule1
```
