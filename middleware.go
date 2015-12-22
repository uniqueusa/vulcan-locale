package locale

import (
	"fmt"
	"net/http"
)

// host struct represents a single configuration for an origin.
type domain struct {
	Locales    []string
	Currencies []string
}

// Middleware struct holds configuration parameters.
type Middleware struct {
	Domains map[string]*domain
}

// NewHandler initializes a new handler from the middleware config and adds it to the middleware chain.
func (m *Middleware) NewHandler(next http.Handler) (http.Handler, error) {
	return &Handler{next: next, cfg: *m}, nil
}

// String() will be called by loggers inside Vulcand and command line tool.
func (m *Middleware) String() string {
	return fmt.Sprintf("domains=%v", m.Domains)
}

// Looks for the given origin or "*" if present.
func (m *Middleware) findDomain(origin string) *domain {
	domain := m.Domains[origin]

	return domain
}
