package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Mannan-Ali/Backend/db"
	"github.com/Mannan-Ali/Backend/db/internaldb"
	"github.com/Mannan-Ali/Backend/httputils"
)

// We created this functions as we are using all this in every routing file.

// Read about  in about.txt point 15
type CalcFunc[T any] func(prices []float64, dates []time.Time) (T, error)

func HandleCalcRequest[T any](apiCgf *APIConfig, w http.ResponseWriter, r *http.Request, fn CalcFunc[T]) {

	prices, dates, err := apiCgf.HandlerFetchPrice(r)
	if err != nil {
		httputils.ResponseWithError(w, 400, err.Error())
		return
	}
	if len(prices) == 0 || len(dates) == 0 {
		httputils.ResponseWithError(w, 400, "no price data")
		return
	}

	result, err := fn(prices, dates)
	if err != nil {
		httputils.ResponseWithError(w, 400, err.Error())
		return
	}

	httputils.ResponseWithJSON(w, 200, result)
}

func (apiCgf *APIConfig) HandlerFetchPrice(r *http.Request) ([]float64, []time.Time, error) {
	//parse dates from r (user input)
	stockIDStr := r.URL.Query().Get("stock_id")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	// handling date errors
	stockID, startDate, endDate, err := HandlerQueryParams(stockIDStr, startDateStr, endDateStr)
	if err != nil {
		return nil, nil, err
	}
	//fix the returned format for the adjusted_col to float array/slice so that we can call the periodiccalcreturn func with it
	dbPrices, err := apiCgf.DB.GetDailyPricesByStockInDateRange(r.Context(),
		internaldb.GetDailyPricesByStockInDateRangeParams{
			StockID:       stockID,
			FromTradeDate: startDate,
			ToTradeDate:   endDate,
		})
	if err != nil {
		return nil, nil, err
	}
	// do return json using the model.go file.
	cvDbPrices, dbDates := db.DatabasePricesToFloatsAndDates(dbPrices)
	return cvDbPrices, dbDates, nil
}

func HandlerQueryParams(stockIDStr, startDateStr, endDateStr string) (int32, time.Time, time.Time, error) {

	//1. empty strings
	if stockIDStr == "" || startDateStr == "" || endDateStr == "" {
		return 0, time.Time{}, time.Time{}, errors.New("missing stock ID, start date, or end date")
	}

	//2. Right dates passed

	//ParseInt returns an int64
	stockID, err := strconv.ParseInt(stockIDStr, 10, 32)
	if err != nil {
		return 0, time.Time{}, time.Time{}, errors.New("invalid stock ID")
	}
	//what is layout read in about.txt 13
	reference_layout := "2006-01-02" // Standard YYYY-MM-DD format

	startDate, err := time.Parse(reference_layout, startDateStr)
	if err != nil {
		return 0, time.Time{}, time.Time{}, errors.New("start_date Input Invalid")

	}
	endDate, err := time.Parse(reference_layout, endDateStr)
	if err != nil {
		return 0, time.Time{}, time.Time{}, errors.New("end_date Input Invalid")
	}
	//if start_date is after end_date
	if startDate.After(endDate) {
		return 0, time.Time{}, time.Time{}, errors.New("start_date cannot be higher then end_date")
	}
	return int32(stockID), startDate, endDate, nil
}
