package calc

import (
	"errors"
	"fmt"
	"time"
)

type SharpeRatio struct {
	SharpeRatioValue float64 `json:"sharpe_ratio_value"`
}

// read about.txt point 16 why are we using daily return mean annualized instead of using cagr as they sound same
// and also as we are not using cagr now so sharp ratio can be done for sub years too
func CalculateSharpeRatio(prices []float64, dates []time.Time, riskFreeRate float64) (SharpeRatio, error) {
	// Anything less than 3 months of data contains too much statistical noise
	if len(dates) < 64 || len(prices) < 64 {
		return SharpeRatio{}, errors.New("cannot calculate Sharpe Ratio for less then 3 months")
	}

	fmt.Printf("Calculating Sharpe Ratio from: %v to %v\n", dates[0].Format("2006-01-02"), dates[len(dates)-1].Format("2006-01-02"))

	dailyReturns, err := CalculatePeriodicReturns(prices, dates)
	if err != nil {
		return SharpeRatio{}, err
	}
	stdDeviationVal, err := CalculateStandardDeviation(dailyReturns.PeriodicReturnVal, dates)
	if err != nil {
		return SharpeRatio{}, err
	}
	// Safeguard against Division by Zero
	if stdDeviationVal.AnnualizedVolatility == 0 {
		return SharpeRatio{}, errors.New("standard deviation is exactly zero; Sharpe ratio is undefined")
	}
	annualizedMeanReturn := stdDeviationVal.Mean * 252.0
	sharpeRatio := (annualizedMeanReturn - riskFreeRate) / stdDeviationVal.AnnualizedVolatility

	return SharpeRatio{
		SharpeRatioValue: sharpeRatio,
	}, nil
}
