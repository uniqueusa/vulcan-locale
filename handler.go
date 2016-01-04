package locale

import (
	"fmt"
	"net/http"
	"regexp"
)

// Handler executes Locale and handles the middleware chain to the next in stack
type Handler struct {
	cfg  Middleware
	next http.Handler
}

type localeRequest struct {
	locale   string
	currency string
	hostName string
}

// Runs the Locale specification on the request before passing it to the next middleware chain
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handleRequest(w, r)

	h.next.ServeHTTP(w, r)
}

// Runs the Locale specification for standard requests
func (h *Handler) handleRequest(w http.ResponseWriter, r *http.Request) {
	domain := h.getDomain(r.Header.Get("Origin"))
	if domain == nil {
		return
	}
	matchedLocale := defaultLanguage
	matchedCurrency := defaultCurrency
	if domain != nil {
		fmt.Println("** Domain found ** ")
		query := r.URL.Query()
		matchedLocale = domain.Locales[0]
		matchedCurrency = domain.Currencies[0]

		if loc := query.Get(queryStringLanguage); loc != "" && stringInSlice(loc, domain.Locales) {
			matchedLocale = loc
		}

		if cur := query.Get(queryStringCurrency); cur != "" && stringInSlice(cur, domain.Currencies) {
			matchedCurrency = cur
		}
	}
	fmt.Printf("** setting header %v** ", matchedLocale)
	r.Header.Set(acceptLanguageHeader, matchedLocale)
	r.Header.Set(acceptCurrencyHeader, matchedCurrency)
}

func (h *Handler) getDomain(host string) *domain {
	re, _ := regexp.Compile("/?/?(.*):?")
	match := re.FindAllStringSubmatch(host, -1)
	return h.cfg.findDomain(match[0][1])
}

// Shares common functionality for prefilght and standard requests
func (h *Handler) handleCommon(w http.ResponseWriter, r *http.Request, requestedInfo *localeRequest) {
	hostConfig := h.cfg.findDomain(requestedInfo.hostName)
	if hostConfig == nil {
		return
	}

	h.buildResponse(w, r, requestedInfo)
}

// Writes the Access Control response headers
func (h *Handler) buildResponse(w http.ResponseWriter, r *http.Request, localeInfo *localeRequest) {
	w.Header().Set(acceptLanguageHeader, localeInfo.locale)
	w.Header().Set(acceptCurrencyHeader, localeInfo.currency)
}
