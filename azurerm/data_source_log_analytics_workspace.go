package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceLogAnalyticsWorkspace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLogAnalyticsWorkspaceRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location": locationForDataSourceSchema(),

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"sku": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"retention_in_days": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"workspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"portal_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_shared_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_shared_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tagsForDataSourceSchema(),
		},
	}

}

func dataSourceLogAnalyticsWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).logAnalytics.WorkspacesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Log Analytics workspaces %q (Resource Group %q) was not found", name, resGroup)
		}
		return fmt.Errorf("Error making Read request on AzureRM Log Analytics workspaces '%s': %+v", name, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	d.Set("workspace_id", resp.CustomerID)
	d.Set("portal_url", resp.PortalURL)
	if sku := resp.Sku; sku != nil {
		d.Set("sku", sku.Name)
	}
	d.Set("retention_in_days", resp.RetentionInDays)

	sharedKeys, err := client.GetSharedKeys(ctx, resGroup, name)
	if err != nil {
		log.Printf("[ERROR] Unable to List Shared keys for Log Analytics workspaces %s: %+v", name, err)
	} else {
		d.Set("primary_shared_key", sharedKeys.PrimarySharedKey)
		d.Set("secondary_shared_key", sharedKeys.SecondarySharedKey)
	}

	flattenAndSetTags(d, resp.Tags)
	return nil
}
