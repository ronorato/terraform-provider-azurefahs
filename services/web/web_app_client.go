package web

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
)

const fqdn = "github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"

type WebAppsClient struct {
	appsClient *web.AppsClient
}

func NewWebAppClient(client *web.AppsClient) *WebAppsClient {
	return &WebAppsClient{appsClient: client}
}

// CreateOrUpdateHostSecret prepares the CreateHostSecret request.
func CreateOrUpdateHostSecret(client *web.AppsClient, ctx context.Context, resourceGroupName string, name string, keyType string, keyName string, key KeyInfoProperties) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"keyName":           autorest.Encode("path", keyName),
		"keyType":           autorest.Encode("path", keyType),
		"name":              autorest.Encode("path", name),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/sites/{name}/host/default/{keyType}/{keyName}", pathParameters),
		autorest.WithJSON(key),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateOrUpdateHostSecretResponder handles the response to the CreateOrUpdateHostSecret request. The method always
// closes the http.Response Body.
func CreateOrUpdateHostSecretResponder(resp *http.Response) (result KeyInfoProperties, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// CreateOrUpdateHostSecret description for Add or update a host level secret.
// Parameters:
// resourceGroupName - name of the resource group to which the resource belongs.
// name - site name.
// keyType - the type of host key.
// keyName - the name of the key.
// key - the key to create or update
func (client WebAppsClient) CreateOrUpdateHostSecret(ctx context.Context, resourceGroupName string, name string, keyType string, keyName string, key web.KeyInfo) (result web.KeyInfo, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/AppsClient.CreateOrUpdateHostSecret")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+[^\.]$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("web.AppsClient", "CreateOrUpdateHostSecret", err.Error())
	}

	req, err := client.CreateOrUpdateHostSecretPreparer(ctx, resourceGroupName, name, keyType, keyName, KeyInfoProperties{KeyInfo: &key})

	if err != nil {
		err = autorest.NewErrorWithError(err, "web.AppsClient", "CreateOrUpdateHostSecret", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateOrUpdateHostSecretSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "web.AppsClient", "CreateOrUpdateHostSecret", resp, "Failure sending request")
		return
	}

	r, err := client.CreateOrUpdateHostSecretResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "web.AppsClient", "CreateOrUpdateHostSecret", resp, "Failure responding to request")
		return
	}

	result = web.KeyInfo{Name: r.Name, Value: r.Value}
	return
}

// CreateOrUpdateHostSecretPreparer prepares the CreateOrUpdateHostSecret request.
func (client WebAppsClient) CreateOrUpdateHostSecretPreparer(ctx context.Context, resourceGroupName string, name string, keyType string, keyName string, key KeyInfoProperties) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"keyName":           autorest.Encode("path", keyName),
		"keyType":           autorest.Encode("path", keyType),
		"name":              autorest.Encode("path", name),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.appsClient.SubscriptionID),
	}

	const APIVersion = "2020-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.appsClient.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/sites/{name}/host/default/{keyType}/{keyName}", pathParameters),
		autorest.WithJSON(key),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateOrUpdateHostSecretSender sends the CreateOrUpdateHostSecret request. The method will close the
// http.Response Body if it receives an error.
func (client WebAppsClient) CreateOrUpdateHostSecretSender(req *http.Request) (*http.Response, error) {
	return client.appsClient.Send(req, azure.DoRetryWithRegistration(client.appsClient.Client))
}

// CreateOrUpdateHostSecretResponder handles the response to the CreateOrUpdateHostSecret request. The method always
// closes the http.Response Body.
func (client WebAppsClient) CreateOrUpdateHostSecretResponder(resp *http.Response) (result KeyInfoProperties, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
