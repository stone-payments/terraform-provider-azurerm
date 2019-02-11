package azurerm

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// Default Authorization Rule/Policy created by Azure, used to populate the
// default connection strings and keys
var serviceBusGeoDRDefaultAuthorizationRule = "RootManageSharedAccessKey"

func resourceArmServiceBusGeoDR() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceBusGeoDRCreateUpdate,
		Read:   resourceArmServiceBusGeoDRRead,
		Update: resourceArmServiceBusGeoDRUpdate,
		Delete: resourceArmServiceBusGeoDRDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]{0,100}[a-zA-Z0-9]$"),
					"The Geo-DR Alias can contain only letters, numbers, and hyphens. The alias must start with a letter, and it must end with a letter or number.",
				),
			},

			"primary_resource_group_name": resourceGroupNameSchema(),

			"primary_namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateServiceBusNamespaceName(),
			},

			"secondary_namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: azure.ValidateServiceBusNamespaceName(),
			},

			"secondary_resource_group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: resourceGroupNameSchema().ValidateFunc,
			},

			"default_primary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_primary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_namespace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmServiceBusGeoDRCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusGeoDRClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM ServiceBus Geo-DR Alias creation.")

	name := d.Get("name").(string)
	primaryResourceGroup := d.Get("primary_resource_group_name").(string)
	primaryNamespaceName := d.Get("primary_namespace_name").(string)
	secondaryNamespaceName := d.Get("secondary_namespace_name").(string)
	secondaryResourceGroup := d.Get("secondary_resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, primaryResourceGroup, primaryNamespaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing ServiceBus Geo-DR Alias %q (resource group %q) ID", name, primaryResourceGroup)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_servicebus_disaster_recovery_config", *existing.ID)
		}
	}

	secondaryNamespaceArmId := fmt.Sprintf("/subscriptions/%v/resourceGroups/%v/providers/%v/namespaces/%v", client.SubscriptionID, secondaryResourceGroup, "Microsoft.ServiceBus", secondaryNamespaceName)
	parameters := servicebus.ArmDisasterRecovery{
		ArmDisasterRecoveryProperties: &servicebus.ArmDisasterRecoveryProperties{
			PartnerNamespace: &secondaryNamespaceArmId,
		},
	}

	disasterRecoveryProperties, err := client.CreateOrUpdate(ctx, primaryResourceGroup, primaryNamespaceName, name, parameters)
	if err != nil {
		return err
	}

	if disasterRecoveryProperties.ID == nil {
		return fmt.Errorf("Cannot read ServiceBus Geo-DR Alias %q (resource group %q) ID", name, primaryResourceGroup)
	}

	d.Set("secondary_namespace_id", secondaryNamespaceArmId)
	d.SetId(*disasterRecoveryProperties.ID)
	return resourceArmServiceBusGeoDRRead(d, meta)
}

func resourceArmServiceBusGeoDRUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusGeoDRClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM ServiceBus Geo-DR Alias update.")

	name := d.Get("name").(string)
	primaryResourceGroup := d.Get("primary_resource_group_name").(string)
	primaryNamespaceName := d.Get("primary_namespace_name").(string)
	oldSecondaryNamespaceName, newSecondaryNamespaceName := d.GetChange("secondary_namespace_name")
	oldSecondaryResourceGroup, newSecondaryResourceGroup := d.GetChange("secondary_resource_group_name")

	resp, err := client.Get(ctx, primaryResourceGroup, primaryNamespaceName, name)
	if err != nil {
		return fmt.Errorf("Error making Read request while updating Azure ServiceBus Geo-DR Alias %q: %+v", name, err)
	}

	// If is paired and secondary namespace changed we must break the pairing.
	if *resp.ArmDisasterRecoveryProperties.PartnerNamespace != "" {
		if oldSecondaryNamespaceName != newSecondaryNamespaceName || oldSecondaryResourceGroup != newSecondaryResourceGroup {
			err := breakPairing(meta, primaryNamespaceName, primaryResourceGroup, name)
			if err != nil {
				return err
			}
		}
	}

	return resourceArmServiceBusGeoDRCreateUpdate(d, meta)
}

func resourceArmServiceBusGeoDRRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusGeoDRClient
	ctx := meta.(*ArmClient).StopContext

	// Parse ARM Id
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	primaryResourceGroup := id.ResourceGroup
	primaryNamespaceName := id.Path["namespaces"]
	name := id.Path["disasterRecoveryConfigs"]

	// Parse partner ARM Id
	partnerNamespaceId := d.Get("secondary_namespace_id").(string)
	partnerResourceGroup := ""
	partnerNamespaceName := ""
	if partnerNamespaceId != "" {
		partnerId, err := parseAzureResourceID(partnerNamespaceId)
		if err != nil {
			return err
		}
		partnerResourceGroup = partnerId.ResourceGroup
		partnerNamespaceName = partnerId.Path["namespaces"]
	}

	resp, err := getCurrentPrimaryNamespace(meta, primaryNamespaceName, primaryResourceGroup, partnerNamespaceName, partnerResourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure ServiceBus Geo-DR Alias %q: %+v", name, err)
	}

	// Parse current primary namespace ARM Id
	id, err = parseAzureResourceID(*resp.ID)
	if err != nil {
		return err
	}
	primaryResourceGroup = id.ResourceGroup
	primaryNamespaceName = id.Path["namespaces"]

	// Safe read secondary namespace info who may not exists after the pair is breaked.
	secondaryNamespace := ""
	secondaryResourceGroup := ""
	if *resp.ArmDisasterRecoveryProperties.PartnerNamespace != "" {
		secondaryId, err := parseAzureResourceID(*resp.ArmDisasterRecoveryProperties.PartnerNamespace)
		if err != nil {
			return err
		}
		secondaryNamespace = secondaryId.Path["namespaces"]
		secondaryResourceGroup = secondaryId.ResourceGroup
	}

	d.SetId(*resp.ID)
	d.Set("name", name)
	d.Set("primary_resource_group_name", primaryResourceGroup)
	d.Set("primary_namespace_name", primaryNamespaceName)
	d.Set("secondary_namespace_name", secondaryNamespace)
	d.Set("secondary_resource_group_name", secondaryResourceGroup)
	d.Set("secondary_namespace_id", *resp.ArmDisasterRecoveryProperties.PartnerNamespace)

	keys, err := client.ListKeys(ctx, primaryResourceGroup, primaryNamespaceName, name, serviceBusGeoDRDefaultAuthorizationRule)
	if err != nil {
		log.Printf("[WARN] Unable to List default keys for Geo-DR Alias %q (Resource Group %q): %+v", name, primaryResourceGroup, err)
	} else {
		d.Set("default_primary_connection_string", keys.AliasPrimaryConnectionString)
		d.Set("default_secondary_connection_string", keys.AliasSecondaryConnectionString)
		d.Set("default_primary_key", keys.PrimaryKey)
		d.Set("default_secondary_key", keys.SecondaryKey)
	}

	return nil
}

func resourceArmServiceBusGeoDRDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusGeoDRClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	primaryResourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	name := id.Path["disasterRecoveryConfigs"]

	// Must break the pairing before deleting the resource.
	secondaryNamespaceName := d.Get("secondary_namespace_name").(string)
	if secondaryNamespaceName != "" {
		err := breakPairing(meta, namespaceName, primaryResourceGroup, name)
		if err != nil {
			return err
		}
	}

	deleteResult, err := client.Delete(ctx, primaryResourceGroup, namespaceName, name)
	if err != nil {
		if !response.WasNotFound(deleteResult.Response) {
			return fmt.Errorf("Error deleting ServiceBus Geo-DR Alias %q: %+v", name, err)
		}
	}

	// Wait deletion confirmation.
	err = waitFinalProvisioningState(meta, namespaceName, primaryResourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error waiting ServiceBus Geo-DR Alias %q to be deleted: %+v", name, err)
	}

	// Give a little time to alias name become available again after deletion.
	var sleepSeconds time.Duration
	sleepSeconds = 50
	time.Sleep(time.Second * sleepSeconds)

	return nil
}

func breakPairing(meta interface{}, namespaceName string, resourceGroupName string, alias string) error {
	client := meta.(*ArmClient).serviceBusGeoDRClient
	ctx := meta.(*ArmClient).StopContext

	_, err := client.BreakPairing(ctx, resourceGroupName, namespaceName, alias)
	if err != nil {
		return fmt.Errorf("Error on break pairing operation at ServiceBus Geo-DR Alias: %q: %+v", alias, err)
	}

	// Wait pairing break confirmation.
	err = waitFinalProvisioningState(meta, namespaceName, resourceGroupName, alias)
	if err != nil {
		return fmt.Errorf("Error waiting ServiceBus Geo-DR Alias pairing break: %q: %+v", alias, err)
	}

	return nil
}

func waitFinalProvisioningState(meta interface{}, namespaceName string, resourceGroupName string, alias string) error {
	client := meta.(*ArmClient).serviceBusGeoDRClient
	ctx := meta.(*ArmClient).StopContext

	for succeded := 0; succeded < 1; {
		resp, err := client.Get(ctx, resourceGroupName, namespaceName, alias)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return fmt.Errorf("Error waiting ServiceBus Geo-DR provisioning state update: %q: %+v", alias, err)
		}

		var sleepSeconds time.Duration
		sleepSeconds = 30
		time.Sleep(time.Second * sleepSeconds)
		log.Printf("[INFO] waiting ServiceBus Geo-DR provisioning state update. New try in %v seconds...", sleepSeconds)

		// ProvisioningState is equal to Accepted while processing an state update.
		// Wait until the operation ends.
		if resp.ArmDisasterRecoveryProperties.ProvisioningState != servicebus.Accepted {
			succeded = 1
		}

	}

	return nil
}

func getCurrentPrimaryNamespace(meta interface{}, namespace string, resourceGroupName string, partnerNamespaceName string, partnerResourceGroup string, alias string) (resp servicebus.ArmDisasterRecovery, err error) {
	client := meta.(*ArmClient).serviceBusGeoDRClient
	ctx := meta.(*ArmClient).StopContext

	resp, err = client.Get(ctx, resourceGroupName, namespace, alias)
	if err == nil {
		return
	}

	if !utils.ResponseWasNotFound(resp.Response) {
		return
	}

	// Try get alias resource through secondary namespace.
	if partnerNamespaceName != "" && partnerResourceGroup != "" {
		resp, err = client.Get(ctx, partnerResourceGroup, partnerNamespaceName, alias)
	}

	return
}
