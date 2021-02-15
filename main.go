package main

import (
	"log"

	"terraform-provider-azurefahs/provider"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	// remove date and time stamp from log output as the plugin SDK already adds its own
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.AzureProvider,
	})
}
