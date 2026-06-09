package api

import (
	"net/http"
	"time"

	"github.com/Mannan-Ali/Backend/calc"
	"github.com/Mannan-Ali/Backend/db"
	"github.com/Mannan-Ali/Backend/db/internaldb"
	"github.com/Mannan-Ali/Backend/httputils"
	"github.com/go-chi/chi"
)

func ReturnRoutes(apiCgf *APIConfig) chi.Router {
	rRoute := chi.NewRouter()

	rRoute.Get("/cumulative_return", apiCgf.handlerCalcReturns)
	rRoute.Get("/cagr", apiCgf.handlerCalcCagr)
	return rRoute
}
func (apiCgf *APIConfig) fetchPricehelper(r *http.Request) ([]float64, time.Time, time.Time, error) {
	//parse dates from r (user input)
	stockIDStr := r.URL.Query().Get("stock_id")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	// handling date errors
	stockID, startDate, endDate, err := HandlerQueryParams(stockIDStr, startDateStr, endDateStr)
	if err != nil {
		return nil, time.Time{}, time.Time{}, err
	}
	//fix the returned format for the adjusted_col to float array/slice so that we can call the periodiccalcreturn func with it
	dbPrices, err := apiCgf.DB.GetDailyPricesByStockInDateRange(r.Context(),
		internaldb.GetDailyPricesByStockInDateRangeParams{
			StockID:       stockID,
			FromTradeDate: startDate,
			ToTradeDate:   endDate,
		})
	if err != nil {
		return nil, time.Time{}, time.Time{}, err
	}
	// do return json using the model.go file.
	cvDbPrices := db.DatabaseNullStringsToFloats(dbPrices)
	return cvDbPrices, startDate, endDate, nil
}
func (apiCgf *APIConfig) handlerCalcReturns(w http.ResponseWriter, r *http.Request) {
	cvDbPrices, startDate, endDate, err := apiCgf.fetchPricehelper(r)
	if err != nil {
		httputils.ResponseWithError(w, 500, err.Error())
		return
	}
	calculatedReturns, err := calc.CalculateCumulativeReturns(cvDbPrices, startDate, endDate)
	if err != nil {
		httputils.ResponseWithError(w, 400, "Error while calculating returns")
		return
	}
	httputils.ResponseWithJSON(w, 201, calculatedReturns)
}

func (apiCgf *APIConfig) handlerCalcCagr(w http.ResponseWriter, r *http.Request) {
	cvDbPrices, startDate, endDate, err := apiCgf.fetchPricehelper(r)
	if err != nil {
		httputils.ResponseWithError(w, 500, err.Error())
		return
	}

	calculateCagr, err := calc.CalculateCAGR(cvDbPrices, startDate, endDate)
	if err != nil {
		httputils.ResponseWithError(w, 400, err.Error())
		return
	}
	httputils.ResponseWithJSON(w, 201, calculateCagr)
}
