package locale

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mailgun/vulcand/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/mailgun/vulcand/plugin"
)

// Helper method to read the test configuration file.
func readConfigFile(testfile string) (map[string]*domain, error) {
	configFile, err := ioutil.ReadFile(testfile)
	if err != nil {
		return nil, err
	}

	var config map[string]*domain
	yaml.Unmarshal(configFile, &config)

	return config, nil
}

func setupTestServer(testfile string, host string) *httptest.Server {
	data, _ := readConfigFile(testfile)
	config := map[string]*domain{host: data[host]}
	locale, _ := New(config)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	handler, _ := locale.NewHandler(next)

	return httptest.NewServer(handler)
}

func setupTestRequest(method string, url string) *http.Request {
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Set("Referer", "http://esalerugs.com")

	return req
}

func TestSpecIsOK(t *testing.T) {
	t.Log("Add Locale Middleware spec to Vulcan registry")

	err := plugin.NewRegistry().AddSpec(GetSpec())
	if err != nil {
		t.Errorf("Expected to be able to add spec but got error %+v", err)
	}
}

func TestNew(t *testing.T) {
	t.Log("Creating Locale Middleware with New method")

	config, err := readConfigFile("test.yml")
	if err != nil {
		t.Errorf("Received error while processing config file: %+v", err)
	}

	cm, err := New(config)
	if err != nil {
		t.Errorf("Expected to create middleware but got error: %+v", err)
	}

	if cm == nil {
		t.Errorf("Expected a Locale Middleware instance but got %+v", cm)
	}

	if cm.String() == "" {
		t.Errorf("Expected middleware string %+v but got empty string", cm)
	}

	handler, err := cm.NewHandler(nil)
	if err != nil {
		t.Errorf("Expected to received a handler but got error: %+v", err)
	}

	if handler == nil {
		t.Errorf("Expected a Locale Handler instance but got %+v", handler)
	}
}

func TestNewInvalid(t *testing.T) {
	t.Log("Creating Locale Middleware with invalid data")

	_, err := New(map[string]*domain{})
	if err == nil {
		t.Errorf("Expected to receive an error but got %+v", err)
	}
}

func TestFromOther(t *testing.T) {
	t.Log("Creating Locale Middleware from other Locale Middleware")

	config, err := readConfigFile("test.yml")
	if err != nil {
		t.Errorf("Received error while processing config file: %+v", err)
	}

	cm, err := New(config)
	if err != nil {
		t.Errorf("Expected to create middleware but got error: %+v", err)
	}

	if cm == nil {
		t.Errorf("Expected a Locale Middleware instance but got %+v", cm)
	}

	other, err := FromOther(*cm)
	if err != nil {
		t.Errorf("Expected to create other middleware but got error: %+v", err)
	}

	if other == nil {
		t.Errorf("Expected other middleware to equal %+v but got nil", cm)
	}
}

func TestFromCli(t *testing.T) {
	t.Log("Create Locale Middleware from command line")

	app := cli.NewApp()
	app.Name = "Locale Middleware Test"
	executed := false
	app.Action = func(ctx *cli.Context) {
		executed = true
		cm, err := FromCli(ctx)
		if err != nil {
			t.Errorf("Expected to create middleware but got error: %+v", err)
		}

		if cm == nil {
			t.Errorf("Expected Locale Middleware instance but got %+v", cm)
		}

		originCount := len((cm.(*Middleware)).Domains)
		if originCount != 4 {
			t.Errorf("Expected 4 domains but got %v", originCount)
		}
	}

	app.Flags = CliFlags()
	app.Run([]string{"Locale Middleware Test", "--configFile=test.yml"})
	if !executed {
		t.Errorf("Expected CLI app to run but it did not.")
	}
}

func TestLanguageHeader(t *testing.T) {
	t.Log("Set language header when provided")

	locale := "en_US"
	server := setupTestServer("test1.yml", "127.0.0.1")
	defer server.Close()

	req := setupTestRequest("GET", server.URL)
	res, err := (&http.Client{}).Do(req)

	if err != nil {
		t.Errorf("Error while processing request: %+v", err)
	}

	code := res.StatusCode
	if code != http.StatusOK {
		t.Errorf("Expected HTTP status %v but it was %v", http.StatusOK, code)
	}

	resLocale := res.Header.Get(acceptLanguageHeader)
	if resLocale != locale {
		t.Errorf("Expected Language header %v but it was %v", locale, resLocale)
	}
}

func TestCurrencyHeader(t *testing.T) {
	t.Log("Set currency header when provided")

	currency := "usd"
	server := setupTestServer("test1.yml", "127.0.0.1")
	defer server.Close()

	req := setupTestRequest("GET", server.URL)
	res, err := (&http.Client{}).Do(req)

	if err != nil {
		t.Errorf("Error while processing request: %+v", err)
	}

	code := res.StatusCode
	if code != http.StatusOK {
		t.Errorf("Expected HTTP status %v but it was %v", http.StatusOK, code)
	}

	resCurrency := res.Header.Get(acceptCurrencyHeader)
	if resCurrency != currency {
		t.Errorf("Expected Currency header %v but it was %v", currency, resCurrency)
	}
}

func TestDefaultWithMultipleLanguages(t *testing.T) {
	t.Log("Set default when there are multiple languages")

	language := "en_GB"
	server := setupTestServer("test_multiple.yml", "127.0.0.1")
	defer server.Close()

	req := setupTestRequest("GET", server.URL)
	res, err := (&http.Client{}).Do(req)

	if err != nil {
		t.Errorf("Error while processing request: %+v", err)
	}

	code := res.StatusCode
	if code != http.StatusOK {
		t.Errorf("Expected HTTP status %v but it was %v", http.StatusOK, code)
	}

	resLocale := res.Header.Get(acceptLanguageHeader)
	if resLocale != language {
		t.Errorf("Expected Language header %v but it was %v", language, resLocale)
	}
}

func TestDefaultWithMultipleCurrencies(t *testing.T) {
	t.Log("Set default when there are multiple Currencies")

	currency := "eur"
	server := setupTestServer("test_multiple.yml", "127.0.0.1")
	defer server.Close()

	req := setupTestRequest("GET", server.URL)
	res, err := (&http.Client{}).Do(req)

	if err != nil {
		t.Errorf("Error while processing request: %+v", err)
	}

	code := res.StatusCode
	if code != http.StatusOK {
		t.Errorf("Expected HTTP status %v but it was %v", http.StatusOK, code)
	}

	resCurrency := res.Header.Get(acceptCurrencyHeader)
	if resCurrency != currency {
		t.Errorf("Expected Currency header %v but it was %v", currency, resCurrency)
	}
}

func TestSpecifiedWithMultipleLanguages(t *testing.T) {
	t.Log("Set specified language when there are multiple languages")

	language := "fr_FR"
	server := setupTestServer("test_multiple.yml", "127.0.0.1")
	defer server.Close()

	req := setupTestRequest("GET", server.URL+"?_l=fr_FR")
	res, err := (&http.Client{}).Do(req)

	if err != nil {
		t.Errorf("Error while processing request: %+v", err)
	}

	code := res.StatusCode
	if code != http.StatusOK {
		t.Errorf("Expected HTTP status %v but it was %v", http.StatusOK, code)
	}

	resLocale := res.Header.Get(acceptLanguageHeader)
	if resLocale != language {
		t.Errorf("Expected Language header %v but it was %v", language, resLocale)
	}
}

func TestSpecifiedWithMultipleCurrencies(t *testing.T) {
	t.Log("Set specified when there are multiple Currencies")

	currency := "franc"
	server := setupTestServer("test_multiple.yml", "127.0.0.1")
	defer server.Close()

	req := setupTestRequest("GET", server.URL+"?_c=franc")
	res, err := (&http.Client{}).Do(req)

	if err != nil {
		t.Errorf("Error while processing request: %+v", err)
	}

	code := res.StatusCode
	if code != http.StatusOK {
		t.Errorf("Expected HTTP status %v but it was %v", http.StatusOK, code)
	}

	resCurrency := res.Header.Get(acceptCurrencyHeader)
	if resCurrency != currency {
		t.Errorf("Expected Currency header %v but it was %v", currency, resCurrency)
	}
}

func TestDefaultWhenSpecifiedLanguageNotInConfig(t *testing.T) {
	t.Log("Set default language when the specified language is not in config")

	language := "en_GB"
	server := setupTestServer("test_multiple.yml", "127.0.0.1")
	defer server.Close()

	req := setupTestRequest("GET", server.URL+"?_l=es_SP")
	res, err := (&http.Client{}).Do(req)

	if err != nil {
		t.Errorf("Error while processing request: %+v", err)
	}

	code := res.StatusCode
	if code != http.StatusOK {
		t.Errorf("Expected HTTP status %v but it was %v", http.StatusOK, code)
	}

	resLocale := res.Header.Get(acceptLanguageHeader)
	if resLocale != language {
		t.Errorf("Expected Language header %v but it was %v", language, resLocale)
	}
}

func TestDefaultWhenSpecifiedCurrencyNotInConfig(t *testing.T) {
	t.Log("Set defaul when specified currency is not in config")

	currency := "eur"
	server := setupTestServer("test_multiple.yml", "127.0.0.1")
	defer server.Close()

	req := setupTestRequest("GET", server.URL+"?_c=yen")
	res, err := (&http.Client{}).Do(req)

	if err != nil {
		t.Errorf("Error while processing request: %+v", err)
	}

	code := res.StatusCode
	if code != http.StatusOK {
		t.Errorf("Expected HTTP status %v but it was %v", http.StatusOK, code)
	}

	resCurrency := res.Header.Get(acceptCurrencyHeader)
	if resCurrency != currency {
		t.Errorf("Expected Currency header %v but it was %v", currency, resCurrency)
	}
}
