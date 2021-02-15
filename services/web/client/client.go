package client

import (
	"terraform-provider-azurefahs/common"

	azureweb "github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
)

type Client struct {
	AppServiceEnvironmentsClient *azureweb.AppServiceEnvironmentsClient
	AppServicePlansClient        *azureweb.AppServicePlansClient
	AppServicesClient            *azureweb.AppsClient
	BaseClient                   *azureweb.BaseClient
	CertificatesClient           *azureweb.CertificatesClient
	CertificatesOrderClient      *azureweb.AppServiceCertificateOrdersClient
}

func NewClient(o *common.ClientOptions) *Client {
	appServiceEnvironmentsClient := azureweb.NewAppServiceEnvironmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&appServiceEnvironmentsClient.Client, o.ResourceManagerAuthorizer)

	appServicePlansClient := azureweb.NewAppServicePlansClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&appServicePlansClient.Client, o.ResourceManagerAuthorizer)

	appServicesClient := azureweb.NewAppsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&appServicesClient.Client, o.ResourceManagerAuthorizer)

	baseClient := azureweb.NewWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&baseClient.Client, o.ResourceManagerAuthorizer)

	certificatesClient := azureweb.NewCertificatesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&certificatesClient.Client, o.ResourceManagerAuthorizer)

	certificatesOrderClient := azureweb.NewAppServiceCertificateOrdersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&certificatesOrderClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AppServiceEnvironmentsClient: &appServiceEnvironmentsClient,
		AppServicePlansClient:        &appServicePlansClient,
		AppServicesClient:            &appServicesClient,
		BaseClient:                   &baseClient,
		CertificatesClient:           &certificatesClient,
		CertificatesOrderClient:      &certificatesOrderClient,
	}
}
