package api

import (
	"net/http"

	"github.com/Mannan-Ali/Backend/calc"
	"github.com/Mannan-Ali/Backend/httputils"
	"github.com/go-chi/chi"
)

func ReturnRoutes(apiCgf *APIConfig) chi.Router {
	rRoute := chi.NewRouter()

	rRoute.Get("/cumulative_return", apiCgf.handlerCalcReturns)
	rRoute.Get("/cagr", apiCgf.handlerCalcCagr)
	return rRoute
}
func (apiCgf *APIConfig) handlerCalcReturns(w http.ResponseWriter, r *http.Request) {
	cvDbPrices, dbDates, err := apiCgf.FetchPricehelper(r)
	if err != nil {
		httputils.ResponseWithError(w, 500, err.Error())
		return
	}
	if len(cvDbPrices) == 0 || len(dbDates) == 0 {
		httputils.ResponseWithError(w, 400, "Length of Prices/Dates are missing")
		return
	}

	calculatedReturns, err := calc.CalculateCumulativeReturns(cvDbPrices, dbDates[0], dbDates[len(dbDates)-1])
	if err != nil {
		httputils.ResponseWithError(w, 400, err.Error())
		return
	}
	httputils.ResponseWithJSON(w, 201, calculatedReturns)
}

func (apiCgf *APIConfig) handlerCalcCagr(w http.ResponseWriter, r *http.Request) {
	cvDbPrices, dbDates, err := apiCgf.FetchPricehelper(r)
	if err != nil {
		httputils.ResponseWithError(w, 500, err.Error())
		return
	}
	if len(dbDates) == 0 {
		httputils.ResponseWithError(w, 400, "no price data found")
		return
	}

	calculateCagr, err := calc.CalculateCAGR(cvDbPrices, dbDates[0], dbDates[len(dbDates)-1])
	if err != nil {
		httputils.ResponseWithError(w, 400, err.Error())
		return
	}
	httputils.ResponseWithJSON(w, 201, calculateCagr)
}
