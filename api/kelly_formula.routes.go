package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Mannan-Ali/Backend/calc"
	"github.com/Mannan-Ali/Backend/db"
	"github.com/Mannan-Ali/Backend/db/internaldb"
	"github.com/Mannan-Ali/Backend/httputils"
	"github.com/go-chi/chi"
)

// read point 26 in about.txt
type MultiKellyResponse struct {
	FullKelly map[string]float64 `json:"full_kelly"`
	HalfKelly map[string]float64 `json:"half_kelly"`
}

func KellyFormulaRoutes(apiCgf *APIConfig) chi.Router {
	rRoute := chi.NewRouter()

	rRoute.Get("/single", apiCgf.handlerCalcSingleKelly)
	rRoute.Get("/multi", apiCgf.handlerCalcMultiKelly)
	return rRoute
}

func (apiCgf *APIConfig) handlerCalcSingleKelly(w http.ResponseWriter, r *http.Request) {
	riskFreeRate := 0.07 // Default risk-free rate

	riskFreeRateStr := r.URL.Query().Get("risk_free_rate")
	if riskFreeRateStr != "" {
		parsedRate, err := strconv.ParseFloat(riskFreeRateStr, 64)
		if err != nil {
			httputils.ResponseWithError(w, 400, "Invalid risk_free_rate parameter. Must be a valid float.")
			return
		}
		riskFreeRate = parsedRate
	}

	prices, dates, err := apiCgf.HandlerFetchPrice(r)
	if err != nil {
		httputils.ResponseWithError(w, 400, err.Error())
		return
	}

	result, err := calc.CalculateSingleKellyPercentage(prices, dates, riskFreeRate)
	if err != nil {
		httputils.ResponseWithError(w, 400, err.Error())
		return
	}

	httputils.ResponseWithJSON(w, 200, result)
}

func (apiCgf *APIConfig) handlerCalcMultiKelly(w http.ResponseWriter, r *http.Request) {
	// Step 1: Parse Parameters
	params := r.URL.Query()

	startDateStr := params.Get("start_date")
	endDateStr := params.Get("end_date")
	riskFreeRateStr := params.Get("risk_free_rate")

	riskFreeRate := 0.07 // Default
	if riskFreeRateStr != "" {
		parsedRate, err := strconv.ParseFloat(riskFreeRateStr, 64)
		if err != nil {
			httputils.ResponseWithError(w, 400, "Invalid risk_free_rate parameter")
			return
		}
		riskFreeRate = parsedRate
	}

	stockIDs, err := apiCgf.ResolveSymbols(r)
	if err != nil {
		httputils.ResponseWithError(w, 400, err.Error())
		return
	}

	if len(stockIDs) < 2 {
		httputils.ResponseWithError(w, 400, "Multi-Kelly requires at least 2 valid symbols")
		return
	}

	startDate, endDate, err := HandlerQueryParams(startDateStr, endDateStr)
	if err != nil {
		httputils.ResponseWithError(w, 400, err.Error())
		return
	}

	// Step 2: Fetch Data
	dbRows, err := apiCgf.DB.GetMultiStockPricesInDateRange(r.Context(), internaldb.GetMultiStockPricesInDateRangeParams{
		StockIds:      stockIDs,
		FromTradeDate: startDate,
		ToTradeDate:   endDate,
	})
	if err != nil {
		httputils.ResponseWithError(w, 500, "failed to fetch price data")
		return
	}

	// Step 3: Align Data
	alignedPrices, alignedIDs, err := db.AlignMultiStockPrices(dbRows)
	if err != nil {
		httputils.ResponseWithError(w, 400, err.Error()) // e.g., "no common trading dates"
		return
	}
	// Warn if some requested stocks were dropped
	// later we would return it somehow after bulding frontend
	if len(alignedIDs) != len(stockIDs) {
		// Log it server-side and/or include in response
		fmt.Printf("Warning: requested %d stocks but only %d had data\n", len(stockIDs), len(alignedIDs))
	}
	// Step 4: Calculate
	result, err := calc.CalculateMultiKellyPercentage(alignedPrices, riskFreeRate)
	if err != nil {
		httputils.ResponseWithError(w, 400, err.Error()) // e.g., "failed to invert matrix"
		return
	}

	// Step 5: Respond with a clear map of stock IDs to their allocations
	fullKellyMap := make(map[string]float64)
	halfKellyMap := make(map[string]float64)
	for i, id := range alignedIDs {
		idStr := strconv.Itoa(int(id))
		fullAllocation := result.OptimalAllocations[i]
		fullKellyMap[idStr] = fullAllocation
		halfKellyMap[idStr] = fullAllocation / 2.0
	}

	response := MultiKellyResponse{
		FullKelly: fullKellyMap,
		HalfKelly: halfKellyMap,
	}
	httputils.ResponseWithJSON(w, 200, response)
}
