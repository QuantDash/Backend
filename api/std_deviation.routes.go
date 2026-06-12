package api

import (
	"net/http"

	"github.com/Mannan-Ali/Backend/calc"
	"github.com/go-chi/chi"
)

func StandardDeviationRoutes(apiCgf *APIConfig) chi.Router {
	rRoute := chi.NewRouter()

	rRoute.Get("/std_deviation", apiCgf.handlerCalcStandardDeviation)
	return rRoute
}

func (apiCgf *APIConfig) handlerCalcStandardDeviation(w http.ResponseWriter, r *http.Request) {
	HandleCalcRequest(apiCgf, w, r, calc.CalculateStandardDeviation)

}
