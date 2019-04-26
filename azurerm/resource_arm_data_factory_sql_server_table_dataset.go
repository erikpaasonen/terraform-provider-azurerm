package azurerm

import (
	"fmt"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataFactorySQLServerTableDataset() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDataFactorySQLServerTableDatasetCreateOrUpdate,
		Read:   resourceArmDataFactorySQLServerTableDatasetRead,
		Update: resourceArmDataFactorySQLServerTableDatasetCreateOrUpdate,
		Delete: resourceArmDataFactorySQLServerTableDatasetDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`),
					`Invalid name for Data Factory, see https://docs.microsoft.com/en-us/azure/data-factory/naming-rules`,
				),
			},

			"data_factory_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"tags": tagsSchema(),
		},
	}
}

func resourceArmDataFactorySQLServerTableDatasetCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataFactoryDatasetClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	dataFactoryName := d.Get("data_factory_name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Data Factory SQL Server Table Dataset %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_factory_sql", *existing.ID)
		}
	}

	dataset := datafactory.DatasetResource{}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, dataFactoryName, name, dataset, ""); err != nil {
		return fmt.Errorf("Error creating/updating Data Factory SQL Server Table Dataset %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Data Factory SQL Server Table Dataset %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Data Factory SQL Server Table Dataset %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmDataFactorySQLServerTableDatasetRead(d, meta)
}

func resourceArmDataFactorySQLServerTableDatasetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataFactoryDatasetClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	dataFactoryName := id.Path["factories"]
	name := id.Path["datasets"]

	resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Data Factory SQL Server Table Dataset %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)

	return nil
}

func resourceArmDataFactorySQLServerTableDatasetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataFactoryDatasetClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	dataFactoryName := id.Path["factories"]
	name := id.Path["datasets"]

	response, err := client.Delete(ctx, resourceGroup, dataFactoryName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("Error deleting Data Factory SQL Server Table Dataset %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
		}
	}

	return nil
}
