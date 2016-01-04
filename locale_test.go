package locale

import (
	"github.com/vulcand/vulcand/Godeps/_workspace/src/github.com/vulcand/oxy/testutils"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vulcand/vulcand/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/vulcand/vulcand/plugin"
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

func setupTestServer(testfile string, host string, handler http.HandlerFunc) *httptest.Server {
	data, _ := readConfigFile(testfile)
	config := map[string]*domain{host: data[host]}
	middleware, _ := New(config)

	h := http.HandlerFunc(handler)

	loc, _ := middleware.NewHandler(h)
	return httptest.NewServer(loc)
}

func setupTestRequest(method string, url string) *http.Request {
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Set("Referer", "http://esalerugs.com")

	return req
}

func setupTestHandler(header string, expectedResult string, t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val := r.Header.Get(header)

		if val != expectedResult {
			t.Errorf("Expected %v to be %v", header, expectedResult)
		}
		io.WriteString(w, "treasure")
	}
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
	server := setupTestServer("test1.yml", "127.0.0.1", setupTestHandler(acceptLanguageHeader, locale, t))
	defer server.Close()

	_, _, err := testutils.Get(server.URL, testutils.Header("Origin", "127.0.0.1"))

	if err != nil {
		t.Errorf("Error while processing request: %+v", err)
	}
}

func TestCurrencyHeader(t *testing.T) {
	t.Log("Set currency header when provided")
	currency := "usd"
	server := setupTestServer("test1.yml", "127.0.0.1", setupTestHandler(acceptCurrencyHeader, currency, t))
	defer server.Close()

	_, _, err := testutils.Get(server.URL, testutils.Header("Origin", "127.0.0.1"))

	if err != nil {
		t.Errorf("Error while processing request: %+v", err)
	}
}

func TestDefaultWithMultipleLanguages(t *testing.T) {
	t.Log("Set default when there are multiple languages")

	language := "en_GB"
	server := setupTestServer("test_multiple.yml", "127.0.0.1", setupTestHandler(acceptLanguageHeader, language, t))
	defer server.Close()

	testutils.Get(server.URL, testutils.Header("Origin", "127.0.0.1"))
}

func TestDefaultWithMultipleCurrencies(t *testing.T) {
	t.Log("Set default when there are multiple Currencies")

	currency := "eur"
	server := setupTestServer("test_multiple.yml", "127.0.0.1", setupTestHandler(acceptCurrencyHeader, currency, t))
	defer server.Close()

	testutils.Get(server.URL, testutils.Header("Origin", "127.0.0.1"))
}

func TestSpecifiedWithMultipleLanguages(t *testing.T) {
	t.Log("Set specified language when there are multiple languages")

	language := "fr_FR"
	server := setupTestServer("test_multiple.yml", "127.0.0.1", setupTestHandler(acceptLanguageHeader, language, t))
	defer server.Close()
	testutils.Get(server.URL+"?_l=fr_FR", testutils.Header("Origin", "127.0.0.1"))

}

func TestSpecifiedWithMultipleCurrencies(t *testing.T) {
	t.Log("Set specified when there are multiple Currencies")

	currency := "franc"
	server := setupTestServer("test_multiple.yml", "127.0.0.1", setupTestHandler(acceptCurrencyHeader, currency, t))
	defer server.Close()

	testutils.Get(server.URL+"?_c=franc", testutils.Header("Origin", "127.0.0.1"))
}

func TestDefaultWhenSpecifiedLanguageNotInConfig(t *testing.T) {
	t.Log("Set default language when the specified language is not in config")

	language := "en_GB"
	server := setupTestServer("test_multiple.yml", "127.0.0.1", setupTestHandler(acceptLanguageHeader, language, t))
	defer server.Close()
	testutils.Get(server.URL+"?_l=es_SP", testutils.Header("Origin", "127.0.0.1"))
}

func TestDefaultWhenSpecifiedCurrencyNotInConfig(t *testing.T) {
	t.Log("Set defaul when specified currency is not in config")

	currency := "eur"
	server := setupTestServer("test_multiple.yml", "127.0.0.1", setupTestHandler(acceptCurrencyHeader, currency, t))
	defer server.Close()

	testutils.Get(server.URL+"?_c=yen", testutils.Header("Origin", "127.0.0.1"))
}

func TestOriginNotSupplied(t *testing.T) {
	t.Log("Call next when Origin not supplied, set to en_US")

	server := setupTestServer("test_multiple.yml", "127.0.0.1", func(w http.ResponseWriter, r *http.Request) {
		val := r.Header.Get(acceptLanguageHeader)

		if val != "" {
			t.Errorf("Expected language header to be empty")
		}
		io.WriteString(w, "treasure")
	})
	defer server.Close()

	res, _, _ := testutils.Get(server.URL)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected quest to be ok")
	}

}
