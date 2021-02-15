package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"
	"terraform-provider-azurefahs/helpers/azure"
)

type FunctionAppId struct {
	SubscriptionId string
	ResourceGroup  string
	SiteName       string
}

func NewFunctionAppID(subscriptionId, resourceGroup, siteName string) FunctionAppId {
	return FunctionAppId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SiteName:       siteName,
	}
}

func (id FunctionAppId) String() string {
	segments := []string{
		fmt.Sprintf("Site Name %q", id.SiteName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Function App", segmentsStr)
}

func (id FunctionAppId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SiteName)
}

// FunctionAppID parses a FunctionApp ID into an FunctionAppId struct
func FunctionAppID(input string) (*FunctionAppId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FunctionAppId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SiteName, err = id.PopSegment("sites"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
