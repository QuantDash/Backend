package calc

import (
	"errors"
	"fmt"
	"math"
	"time"
)

// read docs to understand what MAR ratio is :
type CAGRDR struct {
	CagrToDrawdownratioValue float64 `json:"cagr_to_drawdown_ratio_value"`
	DurationDays             int64   `json:"duration_days"`
}

func CalculateCagrToDrawdownRatio(prices []float64, dates []time.Time) (CAGRDR, error) {

	fmt.Printf("Calculating CAGR to Drawdown Ratio from: %v to %v\n", dates[0].Format("2006-01-02"), dates[len(dates)-1].Format("2006-01-02"))

	maxDrawdown, err := CalculateMaximumDrawDown(prices, dates)
	if err != nil {
		return CAGRDR{}, err
	}
	if maxDrawdown.MaxDrawDown == 0 {
		return CAGRDR{}, errors.New("cannot calculate MAR: maximum drawdown is zero")
	}
	cagrValue, err := CalculateCAGR(prices, dates)
	if err != nil {
		return CAGRDR{}, err
	}

	return CAGRDR{
		CagrToDrawdownratioValue: cagrValue.CAGRValue / math.Abs(maxDrawdown.MaxDrawDown),
		DurationDays:             int64(dates[len(dates)-1].Sub(dates[0]).Hours() / 24.0),
	}, nil
}
