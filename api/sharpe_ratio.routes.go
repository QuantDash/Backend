package api

import (
	"net/http"
	"strconv"

	"github.com/Mannan-Ali/Backend/calc"
	"github.com/Mannan-Ali/Backend/httputils"
	"github.com/go-chi/chi"
)

func SharpeRatioRoutes(apiCgf *APIConfig) chi.Router {
	rRoute := chi.NewRouter()

	rRoute.Get("/", apiCgf.handlerCalcSharpeRatio)
	return rRoute
}

func (apiCgf *APIConfig) handlerCalcSharpeRatio(w http.ResponseWriter, r *http.Request) {
	riskFreeRate := 0.07

	// If custom rate is provided, parse it and overwrite the default
	riskFreeRateStr := r.URL.Query().Get("risk_free_rate")
	if riskFreeRateStr != "" {
		parsedRate, err := strconv.ParseFloat(riskFreeRateStr, 64)
		if err != nil {
			httputils.ResponseWithError(w, 400, "Invalid risk_free_rate parameter. Must be a valid float.")
			return
		}
		riskFreeRate = parsedRate
	}

	cvDbPrices, dbDates, err := apiCgf.HandlerFetchPrice(r)
	if err != nil {
		httputils.ResponseWithError(w, 400, err.Error())
		return
	}
	if len(cvDbPrices) == 0 || len(dbDates) == 0 {
		httputils.ResponseWithError(w, 400, "Length of Prices/Dates are missing")
		return
	}

	sharpeRatio, err := calc.CalculateSharpeRatio(cvDbPrices, dbDates, riskFreeRate)
	if err != nil {
		httputils.ResponseWithError(w, 400, err.Error())
		return
	}
	httputils.ResponseWithJSON(w, 200, sharpeRatio)
}
