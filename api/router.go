package api

import (
	"github.com/Mannan-Ali/Backend/db/internaldb"
	"github.com/go-chi/chi"
)

// APIConfig is capitalized so it is exported and accessible from main.go
type APIConfig struct {
	DB *internaldb.Queries
}

// NewV1Router initializes and returns the base v1 router
func NewV1Router(apiCfg *APIConfig) chi.Router {
	r := chi.NewRouter()

	// Mount your domain-specific feature routers here
	r.Mount("/returns", ReturnRoutes(apiCfg))
	r.Mount("/dd", DrawdownRoute(apiCfg))
	r.Mount("/std", StandardDeviationRoutes(apiCfg))
	r.Mount("/cagr_to_drawdown_ratio", CAGRToDrawdownRatioRoutes(apiCfg))
	// r.Mount("/kelly", KellyRouter())       // Create kelly.go later and uncomment
	// r.Mount("/drawdown", DrawdownRouter()) // Create drawdown.go later and uncomment

	return r
}
