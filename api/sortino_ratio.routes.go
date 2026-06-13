package api

import (
	"net/http"
	"strconv"

	"github.com/Mannan-Ali/Backend/calc"
	"github.com/Mannan-Ali/Backend/httputils"
	"github.com/go-chi/chi"
)

func SortinoRatioRoutes(apiCgf *APIConfig) chi.Router {
	rRoute := chi.NewRouter()

	rRoute.Get("/", apiCgf.handlerCalcSortinoRatio)
	return rRoute
}

func (apiCgf *APIConfig) handlerCalcSortinoRatio(w http.ResponseWriter, r *http.Request) {
	target := 0.07

	// If custom rate is provided, parse it and overwrite the default
	targetStr := r.URL.Query().Get("target")
	if targetStr != "" {
		parsedRate, err := strconv.ParseFloat(targetStr, 64)
		if err != nil {
			httputils.ResponseWithError(w, 400, "Invalid risk_free_rate parameter. Must be a valid float.")
			return
		}
		target = parsedRate
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

	sortinoRatio, err := calc.CalculateSortinoRatio(cvDbPrices, dbDates, target)
	if err != nil {
		httputils.ResponseWithError(w, 400, err.Error())
		return
	}
	httputils.ResponseWithJSON(w, 200, sortinoRatio)
}
