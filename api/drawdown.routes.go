package api

import (
	"net/http"

	"github.com/Mannan-Ali/Backend/calc"
	"github.com/Mannan-Ali/Backend/httputils"
	"github.com/go-chi/chi"
)

func DrawdownRoute(apiCgf *APIConfig) chi.Router {
	rRoute := chi.NewRouter()

	rRoute.Get("/maximum_drawdown", apiCgf.handlerCalcMaximumDrawdown)
	rRoute.Get("/maximum_drawdown_duration", apiCgf.handlerCalcMaxDurationDrawdown)
	return rRoute
}

func (apicgf *APIConfig) handlerCalcMaximumDrawdown(w http.ResponseWriter, r *http.Request) {

	cvDbPrices, dbDates, err := apicgf.FetchPricehelper(r)
	if err != nil {
		httputils.ResponseWithError(w, 400, err.Error())
		return
	}
	if len(cvDbPrices) == 0 || len(dbDates) == 0 {
		httputils.ResponseWithError(w, 400, "Length of Prices/Dates are missing")
		return
	}
	calculateMaxDrawDown, err := calc.CalculateMaximumDrawDown(cvDbPrices, dbDates)
	if err != nil {
		httputils.ResponseWithError(w, 400, err.Error())
		return
	}
	httputils.ResponseWithJSON(w, 201, calculateMaxDrawDown)
}

func (apicgf *APIConfig) handlerCalcMaxDurationDrawdown(w http.ResponseWriter, r *http.Request) {
	cvDbPrices, dbDates, err := apicgf.FetchPricehelper(r)
	if err != nil {
		httputils.ResponseWithError(w, 400, err.Error())
		return
	}
	if len(cvDbPrices) == 0 || len(dbDates) == 0 {
		httputils.ResponseWithError(w, 400, "Length of Prices/Dates are missing")
		return
	}
	calculateMaxDrawDownDuration, err := calc.CalculateMaximumDrawDownDuration(cvDbPrices, dbDates)
	if err != nil {
		httputils.ResponseWithError(w, 400, err.Error())
		return
	}
	httputils.ResponseWithJSON(w, 201, calculateMaxDrawDownDuration)
}
