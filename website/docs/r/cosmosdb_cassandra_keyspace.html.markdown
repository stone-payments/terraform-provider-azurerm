---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_cassandra_keyspace"
sidebar_current: "docs-azurerm-resource-cosmosdb-cassandra-keyspace"
description: |-
  Manages a Cassandra KeySpace within a Cosmos DB Account.
---

# azurerm_cosmosdb_cassandra_keyspace

Manages a Cassandra KeySpace within a Cosmos DB Account.

## Example Usage

```hcl
data "azurerm_cosmosdb_account" "example" {
  name                = "tfex-cosmosdb-account"
  resource_group_name = "tfex-cosmosdb-account-rg"
}

resource "azurerm_cosmosdb_cassandra_keyspace" "example" {
  name                = "tfex-cosmos-cassandra-keyspace"
  resource_group_name = "${data.azurerm_cosmosdb_account.example.resource_group_name}"
  account_name        = "${data.azurerm_cosmosdb_account.example.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cosmos DB Cassandra KeySpace. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Cosmos DB Cassandra KeySpace is created. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the Cosmos DB Cassandra KeySpace to create the table within. Changing this forces a new resource to be created.


## Attributes Reference

The following attributes are exported:

* `id` - the Cosmos DB Cassandra KeySpace ID.

## Import

Cosmos Cassandra KeySpace can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_cassandra_keyspace.ks1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.DocumentDB/databaseAccounts/account1/apis/cassandra/keyspaces/ks1
```
