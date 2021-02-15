package web

import (
	"fmt"
	"log"

	"terraform-provider-azurefahs/clients"
	"terraform-provider-azurefahs/helpers"
	"terraform-provider-azurefahs/helpers/azure"
	"terraform-provider-azurefahs/timeouts"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceFunctionHostKey() *schema.Resource {
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

			"key_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"host_secret": {
				Type:      schema.TypeString,
				Required:  false,
				Computed:  true,
				Sensitive: false,
			},
		},
	}
}

func dataSourceFunctionAppHostSecretRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("dataSourceFunctionAppHostSecretRead - ResourceData: %+v", d)
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	keyName := d.Get("key_name").(string)

	log.Printf("dataSourceFunctionAppHostSecretRead: name: %s | key: %s", name, keyName)
	functionSettings, err := client.Get(ctx, resourceGroup, name)

	if err != nil {
		if helpers.ResponseWasNotFound(functionSettings.Response) {
			return fmt.Errorf("Error: dppAzure Function App %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on dppAzure Function App %q: %+v", name, err)
	}

	if functionSettings.ID == nil {
		return fmt.Errorf("cannot read ID for dppAzure Function App %q (Resource Group %q)", name, resourceGroup)
	}
	d.SetId(*functionSettings.ID)

	res, err := client.ListHostKeys(ctx, resourceGroup, name)
	if err != nil {
		if helpers.ResponseWasNotFound(res.Response) {
			return fmt.Errorf("Error: dppAzure Function App %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on dppAzure Function App Hostkeys %q: %+v", name, err)
	}

	log.Printf("List of host Kyes rsp: %+v", res)
	log.Printf("res.FunctionKeys = %+v", res.FunctionKeys)
	v, ok := res.FunctionKeys[keyName]
	log.Printf("dataSourceFunctionAppHostSecretRead - res.FunctionKeys[keyName]: name: %s | key: %s, value: %+v | ok: %+v", name, keyName, v, ok)

	if !ok {
		return fmt.Errorf("Error: dataSourceFunctionAppHostSecretRead - dppAzure Function Secret %q (function app %q) was not found", keyName, name)
	}

	log.Printf("Host Secret: Key: %s | Secret: %+v", keyName, *v)
	d.Set("host_secret", *v)

	return nil
}
