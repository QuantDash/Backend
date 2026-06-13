package api

import (
	"net/http"

	"github.com/Mannan-Ali/Backend/calc"
	"github.com/go-chi/chi"
)

func ReturnRoutes(apiCgf *APIConfig) chi.Router {
	rRoute := chi.NewRouter()

	rRoute.Get("/cumulative_return", apiCgf.handlerCalcReturns)
	rRoute.Get("/cagr", apiCgf.handlerCalcCagr)
	rRoute.Get("/periodic_return", apiCgf.handlerCalcPeriodicReturn)
	return rRoute
}
func (apiCgf *APIConfig) handlerCalcReturns(w http.ResponseWriter, r *http.Request) {
	HandleCalcRequest(apiCgf, w, r, calc.CalculateCumulativeReturns)
}

func (apiCgf *APIConfig) handlerCalcCagr(w http.ResponseWriter, r *http.Request) {
	HandleCalcRequest(apiCgf, w, r, calc.CalculateCAGR)
}

func (apiCgf *APIConfig) handlerCalcPeriodicReturn(w http.ResponseWriter, r *http.Request) {
	HandleCalcRequest(apiCgf, w, r, calc.CalculatePeriodicReturns)
}
