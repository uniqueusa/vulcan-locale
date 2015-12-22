package locale

import (
	"net/http"
	"regexp"
)

// Handler executes CORS and handles the middleware chain to the next in stack
type Handler struct {
	cfg  Middleware
	next http.Handler
}

type localeRequest struct {
	locale   string
	currency string
	hostName string
}

// Runs the CORS specification on the request before passing it to the next middleware chain
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handleRequest(w, r)

	h.next.ServeHTTP(w, r)
}

// Runs the CORS specification for standard requests
func (h *Handler) handleRequest(w http.ResponseWriter, r *http.Request) {
	re, _ := regexp.Compile("(.*):")
	match := re.FindAllStringSubmatch(r.Host, -1)
	domain := h.cfg.findDomain(match[0][1])
	query := r.URL.Query()
	matchedLocale := domain.Locales[0]
	matchedCurrency := domain.Currencies[0]

	if query.Get("_l") != "" && stringInSlice(query["_l"][0], domain.Locales) {
		matchedLocale = query["_l"][0]
	}
	if query.Get("_c") != "" && stringInSlice(query["_c"][0], domain.Currencies) {
		matchedCurrency = query["_c"][0]
	}
	w.Header().Set(acceptLanguageHeader, matchedLocale)
	w.Header().Set(acceptCurrencyHeader, matchedCurrency)
}

func (h *Handler) extractLocaleInfo(r *http.Request) *localeRequest {

	return &localeRequest{
		locale:   r.Header.Get(acceptLanguageHeader),
		currency: r.Header.Get(acceptCurrencyHeader),
		hostName: r.URL.Host,
	}
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
