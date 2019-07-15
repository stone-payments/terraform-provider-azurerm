---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_zone"
sidebar_current: "docs-azurerm-resource-private-dns-zone"
description: |-
  Manages a Private DNS Zone.
---

# azurerm_private_dns_zone

Enables you to manage Private DNS zones within Azure DNS. These zones are hosted on Azure's name servers.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acceptanceTestResourceGroup1"
  location = "West US"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "mydomain.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Private DNS Zone. Must be a valid domain name.

* `resource_group_name` - (Required) Specifies the resource group where the resource exists. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Prviate DNS Zone ID.
* `number_of_record_sets` - The current number of record sets in this Private DNS zone.
* `max_number_of_record_sets` - The maximum number of record sets that can be created in this Private DNS zone.
* `max_number_of_virtual_network_links` - The maximum number of virtual networks that can be linked to this Private DNS zone.
* `max_number_of_virtual_network_links_with_registration` - The maximum number of virtual networks that can be linked to this Private DNS zone with registration enabled.

## Import

Private DNS Zones can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_private_dns_zone.zone1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/privateDnsZones/zone1
```
