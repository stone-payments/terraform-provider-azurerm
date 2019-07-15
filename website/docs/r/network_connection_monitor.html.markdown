---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_connection_monitor"
sidebar_current: "docs-azurerm-resource-network-connection-monitor"
description: |-
  Configures a Network Connection Monitor to monitor communication between a Virtual Machine and an endpoint using a Network Watcher.

---

# azurerm_connection_monitor

Configures a Network Connection Monitor to monitor communication between a Virtual Machine and an endpoint using a Network Watcher.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "connection-monitor-rg"
  location = "West US"
}

resource "azurerm_network_watcher" "test" {
  name                = "network-watcher"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_virtual_network" "test" {
  name                = "production-network"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "cmtest-nic"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "cmtest-vm"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.test.id}"]
  vm_size               = "Standard_F2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osdisk"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "cmtest-vm"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}

resource "azurerm_virtual_machine_extension" "test" {
  name                       = "cmtest-vm-network-watcher"
  location                   = "${azurerm_resource_group.test.location}"
  resource_group_name        = "${azurerm_resource_group.test.name}"
  virtual_machine_name       = "${azurerm_virtual_machine.test.name}"
  publisher                  = "Microsoft.Azure.NetworkWatcher"
  type                       = "NetworkWatcherAgentLinux"
  type_handler_version       = "1.4"
  auto_upgrade_minor_version = true
}

resource "azurerm_network_connection_monitor" "test" {
  name                 = "cmtest-connectionmonitor"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  network_watcher_name = "${azurerm_network_watcher.test.name}"

  source {
    virtual_machine_id = "${azurerm_virtual_machine.test.id}"
  }

  destination {
    address = "terraform.io"
    port    = 80
  }

  depends_on = ["azurerm_virtual_machine_extension.test"]
}
```

~> **NOTE:** This Resource requires that [the Network Watcher Agent Virtual Machine Extension](https://docs.microsoft.com/en-us/azure/network-watcher/connection-monitor) is installed on the Virtual Machine before monitoring can be started. The extension can be installed via [the `azurerm_virtual_machine_extension` resource](virtual_machine_extension.html).

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Network Connection Monitor. Changing this forces a new resource to be created.

* `network_watcher_name` - (Required) The name of the Network Watcher. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Connection Monitor. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `auto_start` - (Optional) Specifies whether the connection monitor will start automatically once created. Defaults to `true`. Changing this forces a new resource to be created.

* `interval_in_seconds` - (Optional) Monitoring interval in seconds. Defaults to `60`.

* `source` - (Required) A `source` block as defined below.

* `destination` - (Required) A `destination` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `source` block contains:

* `virtual_machine_id` - (Required) The ID of the Virtual Machine to monitor connectivity from.

* `port` - (Optional) The port on the Virtual Machine to monitor connectivity from. Defaults to `0` (Dynamic Port Assignment).

A `destination` block contains:

* `virtual_machine_id` - (Optional) The ID of the Virtual Machine to monitor connectivity to.

* `address` - (Optional) IP address or domain name to monitor connectivity to.

* `port` - (Required) The port on the destination to monitor connectivity to.

~> **NOTE:** One of `virtual_machine_id` or `address` must be specified.

## Attributes Reference

The following attributes are exported:

* `id` - The Connection Monitor ID.

## Import

Connection Monitors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_connection_monitor.monitor1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/networkWatchers/watcher1/connectionMonitors/monitor1
```
