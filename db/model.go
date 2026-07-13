package db

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/Mannan-Ali/Backend/db/internaldb"
)

// As this file was in root directory it cannot be imported into other subpackages like api so for that we shifted it into db folder
// read point 8 in about.txt
// main reason for creation of this file is we cannot edit internalb/model.go and it reutns nullstring ({adjusted_col: {String: "150.5", Valid: true}})
// now we cannot use this for our maths and also cannot send it to user so we create this function to convert the model.go messy format to a better format

// read point 9 in about.txt for why do we use `json:"id"` here??

type CorporateAction struct {
	StockID        int32     `json:"stock_id"`
	ExDate         time.Time `json:"ex_date"`
	ActionType     string    `json:"action_type"`
	DividendAmount *float64  `json:"dividend_amount"`
	SplitRatio     *float64  `json:"split_ratio"`
}

func DatabaseCorpActionToCorpAction(dbCorpAc internaldb.CorporateAction) CorporateAction {
	// In Go, a normal float64 variable cannot be "empty/null". If you don't give it a value, Go automatically assigns it a default value of 0.0.
	// so if we dont use a pointer (*float64) and the dividend is null then it will be assigned 0.0 val to divident amout
	// but for *float it can be empty/nill as its a pointer hece if nill then we assing null to json corropspoding value for user to see
	var dividendPtr *float64
	var splitPtr *float64
	//can use the newly defined function here but leaving it this way for learning purpose in future
	if dbCorpAc.DividendAmount.Valid {
		// ParseFloat converts the string "1.50" into a float64 1.50
		if val, err := strconv.ParseFloat(dbCorpAc.DividendAmount.String, 64); err == nil {
			dividendPtr = &val // The & symbol gets the memory address (pointer)
		}
	}
	if dbCorpAc.SplitRatio.Valid {
		if val, err := strconv.ParseFloat(dbCorpAc.SplitRatio.String, 64); err == nil {
			splitPtr = &val
		}
	}

	return CorporateAction{
		StockID:        dbCorpAc.StockID,
		ExDate:         dbCorpAc.ExDate,
		ActionType:     dbCorpAc.ActionType,
		DividendAmount: dividendPtr,
		SplitRatio:     splitPtr,
	}
}

type DailyPrice struct {
	StockID        int32     `json:"stock_id"`
	TradeDate      time.Time `json:"trade_date"`
	OpenPrice      *float64  `json:"open_price"`
	HighPrice      *float64  `json:"high_price"`
	LowPrice       *float64  `json:"low_price"`
	ClosePrice     *float64  `json:"close_price"`
	AdjustedOpen   *float64  `json:"adjusted_open"`
	AdjustedHigh   *float64  `json:"adjusted_high"`
	AdjustedLow    *float64  `json:"adjusted_low"`
	AdjustedClose  *float64  `json:"adjusted_close"`
	Volume         *int64    `json:"volume"`
	AdjustedVolume *int64    `json:"adjusted_volume"`
}

func DatabaseDPriceToDPrice(dbDailyPrice internaldb.DailyPrice) DailyPrice {

	OpenPtr := TypeConversionInFloat64Helper(dbDailyPrice.OpenPrice)
	HighPtr := TypeConversionInFloat64Helper(dbDailyPrice.HighPrice)
	LowPtr := TypeConversionInFloat64Helper(dbDailyPrice.LowPrice)
	ClosePtr := TypeConversionInFloat64Helper(dbDailyPrice.ClosePrice)

	adjOpenPtr := TypeConversionInFloat64Helper(dbDailyPrice.AdjustedOpen)
	adjHighPtr := TypeConversionInFloat64Helper(dbDailyPrice.AdjustedHigh)
	adjLowPtr := TypeConversionInFloat64Helper(dbDailyPrice.AdjustedLow)
	adjClosePtr := TypeConversionInFloat64Helper(dbDailyPrice.AdjustedClose)

	volumePtr := TypeConversionInInt64Helper(dbDailyPrice.Volume)
	adjVolumePtr := TypeConversionInInt64Helper(dbDailyPrice.AdjustedVolume)

	return DailyPrice{
		StockID:        dbDailyPrice.StockID,
		TradeDate:      dbDailyPrice.TradeDate,
		OpenPrice:      OpenPtr,
		HighPrice:      HighPtr,
		LowPrice:       LowPtr,
		ClosePrice:     ClosePtr,
		AdjustedOpen:   adjOpenPtr,
		AdjustedHigh:   adjHighPtr,
		AdjustedLow:    adjLowPtr,
		AdjustedClose:  adjClosePtr,
		Volume:         volumePtr,
		AdjustedVolume: adjVolumePtr,
	}
}

// translate exactly one database row into one clean Go struct so when we call functions like GetDailyPrice we get mutiple values
// now to handle that we created this function so we can iterate/loop over the slice we get from interandb (and store it here in our slice) which itself makes a loop
// to store this, as the mysql will return multiple value ,so interdb models has to handle it so it does the same thing with nullstring values
// read point 10 in about.txt for better understanding
//
//	func databasePostsToPosts(dbPosts []database.Post) []Post {
//		posts := []Post{}
//		for _, dbPost := range dbPosts {
//			posts = append(posts, databasePostToPost(dbPost))
//		}
//		return posts
//	}

func DatabaseDPricesToDPrices(dbDailyPrices []internaldb.DailyPrice) []DailyPrice {
	prices := []DailyPrice{}
	for _, dbPrice := range dbDailyPrices {
		prices = append(prices, DatabaseDPriceToDPrice(dbPrice))
	}
	return prices
}

type Symbol struct {
	StockID     int32   `json:"stock_id"`
	Symbol      string  `json:"symbol"`
	CompanyName *string `json:"company_name"`
	Status      *string `json:"status"`
}

func DatabaseSymbolToSymbol(dbSymbol internaldb.Symbol) Symbol {
	var cmpNamePtr *string
	var status *string

	// We only need val, err := ... when we are doing a data type conversion that could potentially fail.
	if dbSymbol.CompanyName.Valid {
		cmpNamePtr = &dbSymbol.CompanyName.String
	}
	if dbSymbol.Status.Valid {
		status = &dbSymbol.Status.String
	}
	return Symbol{
		StockID:     dbSymbol.StockID,
		Symbol:      dbSymbol.Symbol,
		CompanyName: cmpNamePtr,
		Status:      status,
	}
}

func TypeConversionInFloat64Helper(ns sql.NullString) *float64 {
	if ns.Valid {
		// ParseFloat converts the string "1.50" into a float64 1.50
		if val, err := strconv.ParseFloat(ns.String, 64); err == nil {
			return &val // The & symbol gets the memory address (pointer)
		}
	}
	return nil
}

func TypeConversionInInt64Helper(ni sql.NullInt64) *int64 {
	//no need to convert as instad we already have it as NullInt64 not NullString so we can just return the .Int64 part
	if ni.Valid {
		val := ni.Int64
		return &val // The & symbol gets the memory address (pointer)
	}
	return nil
}

// This below comment is not much valid now as we changed the function
// we use this instead of databaseDPricesToDPricesOnDate like function as we can see in internaldb that its not a dailyprice struct for GetPriceOnDate
// it is a sql.Nullstring slice so we wrote this function also this can be used for any other col that returns []nullstring so its like a utility function
// read point about.txt 12 for more (read this for new func too)
func DatabasePricesToFloatsAndDates(dbRows []internaldb.GetDailyPricesByStockInDateRangeRow) ([]float64, []time.Time) {
	floats := []float64{}
	dates := []time.Time{}
	for _, row := range dbRows {
		if row.AdjustedClose.Valid {
			if val, err := strconv.ParseFloat(row.AdjustedClose.String, 64); err == nil {
				floats = append(floats, val)
				dates = append(dates, row.TradeDate)
			}
		}
	}
	return floats, dates
}

// AlignMultiStockPrices takes raw database rows for multiple stocks and aligns them by date.
// It returns a 2D slice of prices where rows are stocks and columns are dates,
// a slice of stock IDs corresponding to the rows, and an error if alignment fails.

// read point 25 and 25.1 in about.txt
func AlignMultiStockPrices(dbRows []internaldb.GetMultiStockPricesInDateRangeRow) ([][]float64, []int32, error) {
	if len(dbRows) == 0 {
		return nil, nil, errors.New("no price data returned from database")
	}

	// Step 1: Group all prices by stock ID and then by date for quick lookups.
	pricesByStock := make(map[int32]map[time.Time]float64)
	stockIDsMap := make(map[int32]struct{})

	for _, row := range dbRows {
		if !row.AdjustedClose.Valid {
			continue // Skip rows with null prices
		}

		price, err := strconv.ParseFloat(row.AdjustedClose.String, 64)
		if err != nil {
			continue // Skip rows with unparseable prices
		}

		// If we see a stock for the first time, initialize its inner map.
		if _, ok := pricesByStock[row.StockID]; !ok {
			pricesByStock[row.StockID] = make(map[time.Time]float64)
		}
		pricesByStock[row.StockID][row.TradeDate] = price
		stockIDsMap[row.StockID] = struct{}{}
	}

	if len(pricesByStock) == 0 {
		return nil, nil, errors.New("no valid price data found to align")
	}

	// Create a sorted slice of stock IDs to ensure a consistent output order.
	stockIDs := make([]int32, 0, len(stockIDsMap))
	for id := range stockIDsMap {
		stockIDs = append(stockIDs, id)
	}
	sort.Slice(stockIDs, func(i, j int) bool { return stockIDs[i] < stockIDs[j] })

	// Step 2: Find the intersection of dates where ALL stocks have price data.
	// Start with the dates of the first stock.
	commonDatesMap := make(map[time.Time]struct{})
	for date := range pricesByStock[stockIDs[0]] {
		commonDatesMap[date] = struct{}{}
	}

	// Filter this set by the dates of all other stocks.
	for _, id := range stockIDs[1:] {
		nextCommonDates := make(map[time.Time]struct{})
		for date := range commonDatesMap {
			if _, ok := pricesByStock[id][date]; ok {
				nextCommonDates[date] = struct{}{}
			}
		}
		commonDatesMap = nextCommonDates
	}

	if len(commonDatesMap) == 0 {
		return nil, nil, errors.New("no common trading dates found for the selected stocks in the given range")
	}

	// Edge Case: Ensure there are enough common dates for a statistically meaningful calculation.
	const minRequiredDataPoints = 64 // Consistent with calc functions (~3 months)
	if len(commonDatesMap) < minRequiredDataPoints {
		return nil, nil, fmt.Errorf("insufficient common trading dates for a reliable calculation: found %d, require at least %d", len(commonDatesMap), minRequiredDataPoints)
	}

	// Convert the map of common dates to a sorted slice.
	sortedDates := make([]time.Time, 0, len(commonDatesMap))
	for date := range commonDatesMap {
		sortedDates = append(sortedDates, date)
	}
	sort.Slice(sortedDates, func(i, j int) bool { return sortedDates[i].Before(sortedDates[j]) })

	// Step 3: Build the final, aligned price matrix.
	alignedPrices := make([][]float64, len(stockIDs))
	for i, id := range stockIDs {
		alignedPrices[i] = make([]float64, len(sortedDates))
		for j, date := range sortedDates {
			alignedPrices[i][j] = pricesByStock[id][date]
		}
	}

	return alignedPrices, stockIDs, nil
}
