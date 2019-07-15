---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_api_schema"
sidebar_current: "docs-azurerm-resource-api-management-api-schema"
description: |-
  Manages an API Schema within an API Management Service.
---

# azurerm_api_management_api_schema

Manages an API Schema within an API Management Service.

## Example Usage

```hcl
data "azurerm_api_management_api" "example" {
  name                = "search-api"
  api_management_name = "search-api-management"
  resource_group_name = "search-service"
  revision            = "2"
}

resource "azurerm_api_management_api_schema" "example" {
  api_name            = "${data.azurerm_api_management_api.example.name}"
  api_management_name = "${data.azurerm_api_management_api.example.api_management_name}"
  resource_group_name = "${data.azurerm_api_management_api.example.resource_group_name}"
  schema_id           = "example-sche,a"
  content_type        = "application/vnd.ms-azure-apim.xsd+xml"
  value               = "${file("api_management_api_schema.xml")}"
}
```

## Argument Reference

The following arguments are supported:

* `schema_id` - (Required) A unique identifier for this API Schema. Changing this forces a new resource to be created.

* `api_name` - (Required) The name of the API within the API Management Service where this API Schema should be created. Changing this forces a new resource to be created.

* `api_management_name` - (Required) The Name of the API Management Service where the API exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The Name of the Resource Group in which the API Management Service exists. Changing this forces a new resource to be created.

* `content_type` - (Required) The content type of the API Schema.

* `value` - (Required) The JSON escaped string defining the document representing the Schema.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the API Management API Schema.

## Import

API Management API Schema's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_api_schema.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1/schemas/schema1
```
