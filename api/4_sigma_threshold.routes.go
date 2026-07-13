package api

import (
	"net/http"

	"github.com/Mannan-Ali/Backend/calc"
	"github.com/go-chi/chi"
)

func SigmaThreshholdRoutes(apiCgf *APIConfig) chi.Router {
	rRoute := chi.NewRouter()

	rRoute.Get("/", apiCgf.handlerCalc4SigmaThreshold)
	return rRoute
}
func (apiCgf *APIConfig) handlerCalc4SigmaThreshold(w http.ResponseWriter, r *http.Request) {
	HandleCalcRequest(apiCgf, w, r, calc.Calculate4SigmaThreshold)
}
