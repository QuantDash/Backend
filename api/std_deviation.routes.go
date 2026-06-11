package api

import (
	"net/http"

	"github.com/Mannan-Ali/Backend/calc"
	"github.com/Mannan-Ali/Backend/httputils"
	"github.com/go-chi/chi"
)

func StandardDeviationRoutes(apiCgf *APIConfig) chi.Router {
	rRoute := chi.NewRouter()

	rRoute.Get("/std_deviation", apiCgf.handlerCalcStandardDeviation)
	return rRoute
}

func (apiCgf *APIConfig) handlerCalcStandardDeviation(w http.ResponseWriter, r *http.Request) {
	cvDbPrices, dbDates, err := apiCgf.FetchPricehelper(r)
	if err != nil {
		httputils.ResponseWithError(w, 500, err.Error())
		return
	}
	if len(cvDbPrices) == 0 || len(dbDates) == 0 {
		httputils.ResponseWithError(w, 400, "Length of Prices/Dates are missing")
		return
	}

	calculateStdDev, err := calc.CalculateStandardDeviation(cvDbPrices, dbDates[0], dbDates[len(dbDates)-1])
	if err != nil {
		httputils.ResponseWithError(w, 400, err.Error())
		return
	}
	httputils.ResponseWithJSON(w, 200, calculateStdDev)

}
