package api

import (
	"github.com/go-chi/chi"
)

// NewV1Router initializes and returns the base v1 router
func NewV1Router() chi.Router {
	r := chi.NewRouter()

	// Mount your domain-specific feature routers here
	r.Mount("/returns", ReturnRoutes())
	// r.Mount("/kelly", KellyRouter())       // Create kelly.go later and uncomment
	// r.Mount("/drawdown", DrawdownRouter()) // Create drawdown.go later and uncomment

	return r
}
