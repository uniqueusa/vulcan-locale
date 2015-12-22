package locale

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
	"io/ioutil"

	"github.com/vulcand/vulcand/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/vulcand/vulcand/plugin"
)

// Type represents the type of Vulcan middleware.
const Type string = "locale"

// GetSpec is part of the Vulcan middleware interface.
func GetSpec() *plugin.MiddlewareSpec {
	return &plugin.MiddlewareSpec{
		Type:      Type,       // A short name for the middleware
		FromOther: FromOther,  // Tells Vulcan how to create middleware from another one
		FromCli:   FromCli,    // Tell Vulcan how to create middleware from command line tool
		CliFlags:  CliFlags(), // Vulcan will add this flags to middleware specific command line tool
	}
}

// New checks input paramters and initializes the middleware
func New(domains map[string]*domain) (*Middleware, error) {
	_, err := validateConfig(domains)
	if err != nil {
		return nil, err
	}

	return &Middleware{domains}, nil
}

// FromOther Will be called by Vulcand when engine or API will read the middleware from the serialized format.
// It's important that the signature of the function will be exactly the same, otherwise Vulcand will fail to register this middleware.
// The first and the only parameter should be the struct itself, no pointers and other variables.
// Function should return middleware interface and error in case if the parameters are wrong.
func FromOther(m Middleware) (plugin.Middleware, error) {
	return New(m.Domains)
}

// FromCli constructs the middleware from the command line.
func FromCli(c *cli.Context) (plugin.Middleware, error) {
	var suppliedConfig map[string]*domain

	configFile := c.String(configFile)
	if configFile != "" {
		yamlFile, err := ioutil.ReadFile(configFile)
		if err != nil {
			fmt.Println(errorFileIO)
		}

		yaml.Unmarshal(yamlFile, &suppliedConfig)
	}

	return New(suppliedConfig)
}

// CliFlags will be used by Vulcan construct help and CLI command for `vctl`
func CliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{"configFile, cf", "", "YAML configuration file", ""},
	}
}

// Validates the configuration file.
func validateConfig(domains map[string]*domain) (bool, error) {
	if len(domains) == 0 {
		return false, errors.New(errorConfigDomains)
	}

	for dom, cfg := range domains {
		if dom == "" || cfg == nil {
			return false, errors.New(errorConfigBadDomain)
		}

		if len(cfg.Locales) == 0 {
			return false, errors.New(errorConfigLocales)
		}

		if len(cfg.Currencies) == 0 {
			return false, errors.New(errorConfigCurrencies)
		}

	}

	return true, nil
}
