package api

import (
	"net/http"

	"github.com/Mannan-Ali/Backend/calc"
	"github.com/go-chi/chi"
)

func CAGRToDrawdownRatioRoutes(apiCgf *APIConfig) chi.Router {
	rRoute := chi.NewRouter()

	rRoute.Get("/", apiCgf.handlerCalcCAGRToDrawdownRatio)
	return rRoute
}

func (apiCgf *APIConfig) handlerCalcCAGRToDrawdownRatio(w http.ResponseWriter, r *http.Request) {
	//we dont want to call the func here so we send it wihout ()
	HandleCalcRequest(apiCgf, w, r, calc.CalculateCagrToDrawdownRatio)
}
