package db

import (
	"database/sql"
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
