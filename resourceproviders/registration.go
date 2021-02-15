package resourceproviders

import (
	"context"
	"terraform-provider-azurefahs/helpers"

	"github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/resources/mgmt/resources"
	"github.com/hashicorp/go-azure-helpers/resourceproviders"
)

func EnsureRegistered(ctx context.Context, client resources.ProvidersClient, availableRPs []resources.Provider, requiredRPs map[string]struct{}) error {
	helpers.DebugLog("[DEBUG] Determining which Resource Providers require Registration")
	providersToRegister := resourceproviders.DetermineResourceProvidersRequiringRegistration(availableRPs, requiredRPs)

	if len(providersToRegister) > 0 {
		helpers.DebugLog("[DEBUG] Registering %d Resource Providers", len(providersToRegister))
		if err := resourceproviders.RegisterForSubscription(ctx, client, providersToRegister); err != nil {
			return err
		}
	} else {
		helpers.DebugLog("[DEBUG] All required Resource Providers are registered")
	}

	return nil
}
