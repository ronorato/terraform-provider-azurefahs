package location

import (
	"context"
	"terraform-provider-azurefahs/helpers"

	"github.com/Azure/go-autorest/autorest/azure"
)

// supportedLocations can be (validly) nil - as such this shouldn't be relied on
var supportedLocations *[]string

// CacheSupportedLocations attempts to retrieve the supported locations from the Azure MetaData Service
// and caches them, for used in enhanced validation
func CacheSupportedLocations(ctx context.Context, env *azure.Environment) {
	locs, err := availableAzureLocations(ctx, env)
	if err != nil {
		helpers.DebugLog("[DEBUG] error retrieving locations: %s. Enhanced validation will be unavailable", err)
		return
	}

	supportedLocations = locs.Locations
}
