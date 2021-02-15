package provider

import (
	"terraform-provider-azurefahs/sdk"
)

//go:generate go run ../tools/generator-services/main.go -path=../../../

func SupportedTypedServices() []sdk.TypedServiceRegistration {
	return []sdk.TypedServiceRegistration{}
}

// func SupportedUntypedServices() []sdk.UntypedServiceRegistration {
// 	return []sdk.UntypedServiceRegistration{
// 		web.Registration{},
// 	}
// }
