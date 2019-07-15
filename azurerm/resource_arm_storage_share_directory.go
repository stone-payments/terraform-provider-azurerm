package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/file/directories"
)

func resourceArmStorageShareDirectory() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageShareDirectoryCreate,
		Read:   resourceArmStorageShareDirectoryRead,
		Update: resourceArmStorageShareDirectoryUpdate,
		Delete: resourceArmStorageShareDirectoryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageShareDirectoryName,
			},
			"share_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
			"storage_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"metadata": storage.MetaDataSchema(),
		},
	}
}

func resourceArmStorageShareDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	storageClient := meta.(*ArmClient).storage

	accountName := d.Get("storage_account_name").(string)
	shareName := d.Get("share_name").(string)
	directoryName := d.Get("name").(string)

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := storage.ExpandMetaData(metaDataRaw)

	resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group: %s", err)
	}

	client, err := storageClient.FileShareDirectoriesClient(ctx, *resourceGroup, accountName)
	if err != nil {
		return fmt.Errorf("Error building File Share Client: %s", err)
	}

	if requireResourcesToBeImported {
		existing, err := client.Get(ctx, accountName, shareName, directoryName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Directory %q (File Share %q / Storage Account %q / Resource Group %q): %s", directoryName, shareName, accountName, *resourceGroup, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			id := client.GetResourceID(accountName, shareName, directoryName)
			return tf.ImportAsExistsError("azurerm_storage_share_directory", id)
		}
	}

	if _, err := client.Create(ctx, accountName, shareName, directoryName, metaData); err != nil {
		return fmt.Errorf("Error creating Directory %q (File Share %q / Account %q): %+v", directoryName, shareName, accountName, err)
	}

	resourceID := client.GetResourceID(accountName, shareName, directoryName)
	d.SetId(resourceID)

	return resourceArmStorageShareDirectoryRead(d, meta)
}

func resourceArmStorageShareDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	storageClient := meta.(*ArmClient).storage

	id, err := directories.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := storage.ExpandMetaData(metaDataRaw)

	resourceGroup, err := storageClient.FindResourceGroup(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group: %s", err)
	}

	client, err := storageClient.FileShareDirectoriesClient(ctx, *resourceGroup, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error building File Share Client: %s", err)
	}

	if _, err := client.SetMetaData(ctx, id.AccountName, id.ShareName, id.DirectoryName, metaData); err != nil {
		return fmt.Errorf("Error updating MetaData for Directory %q (File Share %q / Account %q): %+v", id.DirectoryName, id.ShareName, id.AccountName, err)
	}

	return resourceArmStorageShareDirectoryRead(d, meta)
}

func resourceArmStorageShareDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	storageClient := meta.(*ArmClient).storage

	id, err := directories.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup, err := storageClient.FindResourceGroup(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group for Storage Account %q: %s", id.AccountName, err)
	}
	if resourceGroup == nil {
		log.Printf("[DEBUG] Unable to locate Resource Group for Storage Account %q - assuming removed & removing from state", id.AccountName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.FileShareDirectoriesClient(ctx, *resourceGroup, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error building File Share Client for Storage Account %q (Resource Group %q): %s", id.AccountName, *resourceGroup, err)
	}

	props, err := client.Get(ctx, id.AccountName, id.ShareName, id.DirectoryName)
	if err != nil {
		return fmt.Errorf("Error retrieving Storage Share %q (File Share %q / Account %q / Resource Group %q): %s", id.DirectoryName, id.ShareName, id.AccountName, *resourceGroup, err)
	}

	d.Set("name", id.DirectoryName)
	d.Set("share_name", id.ShareName)
	d.Set("storage_account_name", id.AccountName)

	if err := d.Set("metadata", storage.FlattenMetaData(props.MetaData)); err != nil {
		return fmt.Errorf("Error setting `metadata`: %s", err)
	}

	return nil
}

func resourceArmStorageShareDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	storageClient := meta.(*ArmClient).storage

	id, err := directories.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup, err := storageClient.FindResourceGroup(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group for Storage Account %q: %s", id.AccountName, err)
	}
	if resourceGroup == nil {
		log.Printf("[DEBUG] Unable to locate Resource Group for Storage Account %q - assuming removed already", id.AccountName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.FileShareDirectoriesClient(ctx, *resourceGroup, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error building File Share Client for Storage Account %q (Resource Group %q): %s", id.AccountName, *resourceGroup, err)
	}

	if _, err := client.Delete(ctx, id.AccountName, id.ShareName, id.DirectoryName); err != nil {
		return fmt.Errorf("Error deleting Storage Share %q (File Share %q / Account %q / Resource Group %q): %s", id.DirectoryName, id.ShareName, id.AccountName, *resourceGroup, err)
	}

	return nil
}
