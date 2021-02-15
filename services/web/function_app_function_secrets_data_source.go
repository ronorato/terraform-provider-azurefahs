package web

import (
	"fmt"

	"terraform-provider-azurefahs/clients"
	"terraform-provider-azurefahs/helpers"
	"terraform-provider-azurefahs/helpers/azure"
	"terraform-provider-azurefahs/timeouts"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DataSourceFunctionSecrets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFunctionAppHostKeysRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"function_secrets": {
				Type:      schema.TypeMap,
				Computed:  true,
				Sensitive: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func dataSourceFunctionAppHostKeysRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	functionSettings, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if helpers.ResponseWasNotFound(functionSettings.Response) {
			return fmt.Errorf("Error: AzureRM Function App %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on AzureRM Function App %q: %+v", name, err)
	}

	if functionSettings.ID == nil {
		return fmt.Errorf("cannot read ID for AzureRM Function App %q (Resource Group %q)", name, resourceGroup)
	}
	d.SetId(*functionSettings.ID)

	res, err := client.ListHostKeys(ctx, resourceGroup, name)
	if err != nil {
		if helpers.ResponseWasNotFound(res.Response) {
			return fmt.Errorf("Error: AzureRM Function App %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on AzureRM Function App Hostkeys %q: %+v", name, err)
	}

	d.Set("function_secrets", res.FunctionKeys)

	return nil
}
