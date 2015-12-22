package locale

const (
	// Response Headers
	acceptLanguageHeader string = "Accept-Language"
	acceptCurrencyHeader string = "Accept-Currency"

	// Querystring Parameters
	queryStringLanguage = "_l"
	queryStringCurrency = "_c"

	// Error Messages
	errorConfigDomains    string = "No Domains supplied in config"
	errorConfigBadDomain  string = "Invalid Domain supplied in config"
	errorConfigLocales    string = "No Locales for domain"
	errorConfigCurrencies string = "No Currencies for domain"
	errorFileIO           string = "file error"

	// Common
	configFile string = "configFile"
)
