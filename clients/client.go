package clients

import (
	"context"

	"terraform-provider-azurefahs/common"
	"terraform-provider-azurefahs/services/authorization"
	"terraform-provider-azurefahs/services/resource"
	web "terraform-provider-azurefahs/services/web/client"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/validation"
)

type Client struct {
	// StopContext is used for propagating control from Terraform Core (e.g. Ctrl/Cmd+C)
	StopContext context.Context

	Account       *ResourceManagerAccount
	Authorization *authorization.Client
	Resource      *resource.Client
	Web           *web.Client
}

// NOTE: it should be possible for this method to become Private once the top level Client's removed

func (client *Client) Build(ctx context.Context, o *common.ClientOptions) error {
	autorest.Count429AsRetry = false
	// Disable the Azure SDK for Go's validation since it's unhelpful for our use-case
	validation.Disabled = true

	client.StopContext = ctx

	client.Authorization = authorization.NewClient(o)
	client.Resource = resource.NewClient(o)
	client.Web = web.NewClient(o)

	return nil
}
