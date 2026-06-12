package api

import (
	"net/http"

	"github.com/Mannan-Ali/Backend/calc"
	"github.com/go-chi/chi"
)

func DrawdownRoute(apiCgf *APIConfig) chi.Router {
	rRoute := chi.NewRouter()

	rRoute.Get("/maximum_drawdown", apiCgf.handlerCalcMaximumDrawdown)
	rRoute.Get("/maximum_drawdown_duration", apiCgf.handlerCalcMaxDurationDrawdown)
	return rRoute
}

func (apiCgf *APIConfig) handlerCalcMaximumDrawdown(w http.ResponseWriter, r *http.Request) {
	HandleCalcRequest(apiCgf, w, r, calc.CalculateMaximumDrawDown)
}

func (apiCgf *APIConfig) handlerCalcMaxDurationDrawdown(w http.ResponseWriter, r *http.Request) {
	HandleCalcRequest(apiCgf, w, r, calc.CalculateMaximumDrawDownDuration)
}
