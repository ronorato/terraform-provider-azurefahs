package web

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Web"
}

// // WebsiteCategories returns a list of categories which can be used for the sidebar
// func (r Registration) WebsiteCategories() []string {
// 	return []string{
// 		"App Service (Web Apps)",
// 	}
// }

// // SupportedDataSources returns the supported Data Sources supported by this Service
// func (r Registration) SupportedDataSources() map[string]*schema.Resource {
// 	return map[string]*schema.Resource{
// 		"azurerm_function_app": dataSourceFunctionApp(),
// 		// "azurerm_function_app_host_keys":        dataSourceFunctionAppHostKeys(),
// 	}
// }

// // SupportedResources returns the supported Resources supported by this Service
// func (r Registration) SupportedResources() map[string]*schema.Resource {
// 	return map[string]*schema.Resource{
// 		// "azurerm_function_app":                                      resourceFunctionApp(),
// 	}
// }
