package web

import (
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/Azure/go-autorest/autorest"
)

type KeyInfoProperties struct {
	autorest.Response `json:"-"`
	// SiteProperties - Site resource specific properties
	*web.KeyInfo `json:"properties"`
}
