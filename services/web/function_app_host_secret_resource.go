package web

import (
	"fmt"
	"terraform-provider-azurefahs/clients"
	"terraform-provider-azurefahs/helpers"
	"terraform-provider-azurefahs/helpers/azure"
	"terraform-provider-azurefahs/services/web/parse"
	"terraform-provider-azurefahs/timeouts"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const HostKeyType = "functionKeys"

func ResourceFunctionAppHostSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceFunctionAppHostSecretCreate,
		Read:   resourceFunctionAppHostSecretRead,
		Update: resourceFunctionAppHostSecretUpdate,
		Delete: resourceFunctionAppHostSecretDelete,
		Exists: resourceFunctionHostSecretExists,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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

func resourceFunctionAppHostSecretCreate(d *schema.ResourceData, meta interface{}) error {
	helpers.DebugLog("resourceFunctionAppHostSecretCreate - ResourceData: %+v", d)
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	keyName := d.Get("key_name").(string)

	functionSettings, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !helpers.ResponseWasNotFound(functionSettings.Response) {
			return fmt.Errorf("Error checking for presence of existing Function App %q (Resource Group %q): %s", name, resourceGroup, err)
		}
	}

	if functionSettings.ID == nil {
		return fmt.Errorf("cannot read ID for AzureRM Function App %q (Resource Group %q)", name, resourceGroup)
	}
	d.SetId(*functionSettings.ID)

	res, err := client.ListHostKeys(ctx, resourceGroup, name)
	helpers.DebugLog("ListHostKeys - res: %+v, err: %+v", functionSettings, err)
	if err != nil {
		if helpers.ResponseWasNotFound(res.Response) {
			return fmt.Errorf("Error: AzureRM Function App %q (Resource Group %q) was not found", name, resourceGroup)
		}

		helpers.DebugLog("Error making Read request on AzureRM Function App Hostkeys %q: %+v", name, err)
	}

	if _, ok := res.FunctionKeys[keyName]; ok {
		return helpers.ImportAsExistsError("dppazure_function_app_function_key", keyName)
	}

	webAppClient := NewWebAppClient(client)

	result, err := webAppClient.CreateOrUpdateHostSecret(ctx, resourceGroup, name, HostKeyType, keyName, web.KeyInfo{Name: &keyName})

	if err != nil {
		return err
	}

	helpers.DebugLog("Result = Name: %s, Value: %s", *result.Name, *result.Value)

	d.Set("key_name", result.Name)
	d.Set("host_secret", result.Value)

	return nil
}

func resourceFunctionAppHostSecretRead(d *schema.ResourceData, meta interface{}) error {
	helpers.DebugLog("resourceFunctionAppHostSecretRead (enter) - ResourceData: %+v", d)

	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	helpers.DebugLog("resourceFunctionAppHostSecretRead - FunctionAppID: %+v", d.Id())

	id, err := parse.FunctionAppID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	keyName := d.Get("key_name").(string)

	res, err := client.ListHostKeys(ctx, resourceGroup, name)
	if err != nil {
		if helpers.ResponseWasNotFound(res.Response) {
			return fmt.Errorf("Error: dppAzure Function App %q (Resource Group %q) was not found", name, id.ResourceGroup)
		}

		return fmt.Errorf("Error making Read request on dppAzure Function App Hostkeys %q: %+v", name, err)
	}

	helpers.DebugLog("resourceFunctionAppHostSecretRead - res.FunctionKeys: %+v", res.FunctionKeys)
	if v, ok := res.FunctionKeys[keyName]; ok {
		helpers.DebugLog("resourceFunctionAppHostSecretRead - Found secret - Key: %s | Secret: %s", keyName, *v)
		d.Set("host_secret", *v)
	}

	helpers.DebugLog("resourceFunctionAppHostSecretRead (exit) - host Secret: %s", d.Get("host_secret").(string))
	return nil
}

func resourceFunctionAppHostSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	helpers.DebugLog("resourceFunctionAppHostSecretUpdate - ResourceData: %+v", d)
	helpers.DebugLog("State().Attributes: %+v", d.State().Attributes)

	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FunctionAppID(d.Id())
	if err != nil {
		return err
	}

	keyName := d.Get("key_name").(string)
	secret := d.Get("host_secret").(string)

	if d.HasChange("key_name") {
		prevName, newName := d.GetChange("key_name")
		helpers.DebugLog("Key Name has changed from %v to %v", prevName, newName)
		keyName = prevName.(string)

		if _, err := client.DeleteHostSecret(ctx, id.ResourceGroup, id.SiteName, HostKeyType, keyName); err != nil {
			return err
		}
	}

	keyInfo := web.KeyInfo{
		Name:  &keyName,
		Value: &secret,
	}
	webAppClient := NewWebAppClient(client)
	result, err := webAppClient.CreateOrUpdateHostSecret(ctx, id.ResourceGroup, id.SiteName, HostKeyType, keyName, keyInfo)

	if err != nil {
		return err
	}

	d.Set("key_value", result.Value)

	return nil
}

func resourceFunctionAppHostSecretDelete(d *schema.ResourceData, meta interface{}) error {
	helpers.DebugLog("resourceFunctionAppHostSecretDelete - ResourceData: %+v", d)

	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FunctionAppID(d.Id())
	if err != nil {
		return err
	}

	keyName := d.Get("key_name").(string)

	res, err := client.ListHostKeys(ctx, id.ResourceGroup, id.SiteName)

	if err != nil {
		if helpers.ResponseWasNotFound(res.Response) {
			return fmt.Errorf("Error: AzureRM Function App %q (Resource Group %q) was not found", id.SiteName, id.ResourceGroup)
		}

		return fmt.Errorf("Error making Read request on AzureRM Function App Hostkeys %q: %+v", id.SiteName, err)
	}

	if _, ok := res.FunctionKeys[keyName]; !ok {
		return nil
	}

	if _, err := client.DeleteHostSecret(ctx, id.ResourceGroup, id.SiteName, HostKeyType, keyName); err != nil {
		return err
	}

	return nil
}

func resourceFunctionHostSecretExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	helpers.DebugLog("resourceFunctionHostSecretExists - ResourceData: %+v", d)

	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FunctionAppID(d.Id())
	if err != nil {
		return false, err
	}

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	keyName := d.Get("key_name").(string)

	res, err := client.ListHostKeys(ctx, resourceGroup, name)
	if err != nil {
		if helpers.ResponseWasNotFound(res.Response) {
			return false, fmt.Errorf("Error: dppAzure Function App %q (Resource Group %q) was not found", name, id.ResourceGroup)
		}

		return false, fmt.Errorf("Error making Read request on dppAzure Function App Hostkeys %q: %+v", name, err)
	}

	if _, ok := res.FunctionKeys[keyName]; ok {
		return true, nil
	}

	return false, nil
}
