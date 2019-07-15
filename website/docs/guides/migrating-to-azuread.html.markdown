---
layout: "azurerm"
page_title: "Azure Active Directory: Migrating to the AzureAD Provider"
sidebar_current: "docs-azurerm-migrating-to-azuread"
description: |-
  This page documents how to migrate from using the AzureAD resources within this repository to the resources in the new split-out repository.

---

# Azure Active Directory: Migrating to the AzureAD Provider

In v1.21 of the AzureRM Provider the Azure Active Directory Data Sources and Resources have been split out into a new Provider specifically for Azure Active Directory.

This guide covers how to migrate from using the following Data Sources and Resources in the AzureRM Provider to using them in the new AzureAD Provider:

* Data Source: `azurerm_azuread_application`
* Data Source: `azurerm_azuread_service_principal`
* Resource: `azurerm_azuread_application`
* Resource: `azurerm_azuread_service_principal`
* Resource: `azurerm_azuread_service_principal_password`

## Updating the Provider block

As the AzureAD and AzureRM Provider support the same authentication methods - it's possible to update the Provider block by setting the new Provider name and version, for example:

```hcl
provider "azurerm" {
  version = "=1.28.0"
}
```

can become:

```hcl
provider "azuread" {
  version = "=0.1.0"
}
```

## Updating the Terraform Configurations

The Azure Active Directory Data Sources and Resources have been split out into the new Provider - which means the name of the Data Sources and Resources has changed slightly.

The main difference in naming is that the `azurerm_` prefix has been removed from the names of the Data Sources and Resources - the following table explains the new name for each of the Azure Active Directory resources:


| Type        | Old Name                                   | New Name                           |
| ----------- | ------------------------------------------ | ---------------------------------- |
| Data Source | azurerm_azuread_application                | azuread_application                |
| Data Source | azurerm_azuread_service_principal          | azuread_service_principal          |
| Resource    | azurerm_azuread_application                | azuread_application                |
| Resource    | azurerm_azuread_service_principal          | azuread_service_principal          |
| Resource    | azurerm_azuread_service_principal_password | azuread_service_principal_password |

---

Once the Provider blocks have been updated, it should be possible to replace the `azurerm_` prefix in your Terraform Configuration from each of the AzureAD resources (and any interpolations) so that the new resources in the AzureAD Provider are used instead.

For example the following Terraform Configuration:

```hcl
resource "azurerm_azuread_application" "test" {
  name = "my-application"
}

resource "azurerm_azuread_service_principal" "test" {
  application_id = "${azurerm_azuread_application.test.application_id}"
}

resource "azurerm_azuread_service_principal_password" "test" {
  service_principal_id = "${azurerm_azuread_service_principal.test.id}"
  value                = "bd018069-622d-4b46-bcb9-2bbee49fe7d9"
  end_date             = "2020-01-01T01:02:03Z"
}
```

we can remove the `azurerm_` prefix from each of the resource names and interpolations to use the `AzureAD` provider instead by making this:

```hcl
resource "azuread_application" "test" {
  name = "my-application"
}

resource "azuread_service_principal" "test" {
  application_id = "${azuread_application.test.application_id}"
}

resource "azuread_service_principal_password" "test" {
  service_principal_id = "${azuread_service_principal.test.id}"
  value                = "bd018069-622d-4b46-bcb9-2bbee49fe7d9"
  end_date             = "2020-01-01T01:02:03Z"
}
```

At this point it should be possible to run `terraform init`, which will download the new AzureAD Provider.

## Migrating Resources in the State

Now that we've updated the Provider Block and the Terraform Configuration we need to update the names of the resources in the state.

Firstly, let's list the existing items in the state - we can do this by running `terraform state list`, for example:

```bash
$ terraform state list
azurerm_azuread_application.test
azurerm_azuread_service_principal.test
azurerm_azuread_service_principal_password.import
azurerm_azuread_service_principal_password.test
```

As the Terraform Configuration has been updated - we can move each of the resources in the state using the `terraform state mv` command, for example:

```shell
$ terraform state mv azurerm_azuread_application.test azuread_application.test
Moved azurerm_azuread_application.test to azuread_application.test
```

This needs to be repeated for each of the Azure Active Directory resources which exist in the state.

Once this has been done, running `terraform plan` should show no changes:

```shell
$ terraform plan
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.


------------------------------------------------------------------------

No changes. Infrastructure is up-to-date.

This means that Terraform did not detect any differences between your
configuration and real physical resources that exist. As a result, no
actions need to be performed.
```

At this point you've switched over to using [the new Azure Active Directory provider](http://terraform.io/docs/providers/azuread/index.html)! You can stay up to date with Releases (and file Feature Requests/Bugs) [on the Github repository](https://github.com/terraform-providers/terraform-provider-azuread).
