package api

import (
	"net/http"

	"github.com/Mannan-Ali/Backend/calc"
	"github.com/Mannan-Ali/Backend/db"
	"github.com/Mannan-Ali/Backend/db/internaldb"
	"github.com/Mannan-Ali/Backend/httputils"
	"github.com/go-chi/chi"
)

func ReturnRoutes(apiCgf *APIConfig) chi.Router {
	rRoute := chi.NewRouter()

	rRoute.Get("/", apiCgf.handlerCalcReturns)

	return rRoute
}

func (apiCgf *APIConfig) handlerCalcReturns(w http.ResponseWriter, r *http.Request) {

	//parse dates from r (user input)
	stockIDStr := r.URL.Query().Get("stock_id")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	// handling date errors
	stockID, startDate, endDate, err := HandlerQueryParams(stockIDStr, startDateStr, endDateStr)
	if err != nil {
		httputils.ResponseWithError(w, 400, err.Error())
		return
	}
	//fix the returned format for the adjusted_col to float array/slice so that we can call the periodiccalcreturn func with it
	dbPrices, err := apiCgf.DB.GetDailyPricesByStockInDateRange(r.Context(),
		internaldb.GetDailyPricesByStockInDateRangeParams{
			StockID:       stockID,
			FromTradeDate: startDate,
			ToTradeDate:   endDate,
		})
	if err != nil {
		httputils.ResponseWithError(w, 500, "Failed to get Prices")
		return
	}
	// do return json using the model.go file.
	cvDbPrices := db.DatabaseNullStringsToFloats(dbPrices)
	calculatedReturns, err := calc.CalculateCumulativeReturns(cvDbPrices, startDateStr, endDateStr)
	if err != nil {
		httputils.ResponseWithError(w, 400, "Error while calculating returns")
	}
	httputils.ResponseWithJSON(w, 201, calculatedReturns)
}
