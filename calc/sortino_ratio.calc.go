package calc

import (
	"errors"
	"fmt"
	"time"
)

type SortinoRatio struct {
	SortinoRatioValue float64 `json:"sortino_ratio_value"`
}

func CalculateSortinoRatio(prices []float64, dates []time.Time, target float64) (SortinoRatio, error) {
	// Anything less than 3 months of data contains too much statistical noise
	if len(dates) < 64 || len(prices) < 64 {
		return SortinoRatio{}, errors.New("cannot calculate Sortino Ratio for less then 3 months")
	}

	fmt.Printf("Calculating Sortino Ratio from: %v to %v\n", dates[0].Format("2006-01-02"), dates[len(dates)-1].Format("2006-01-02"))

	dailyReturns, err := CalculatePeriodicReturns(prices, dates)
	if err != nil {
		return SortinoRatio{}, err
	}

	// Convert the user's annualized target into a daily target to compare against daily returns
	//read point 17 in about.txt for why did we divide by 252 and why user passes it annualy and not dialy or monthly
	dailyTarget := target / 252.0
	stdDeviationVal, err := CalculateStandardDownwardDeviation(dailyReturns.PeriodicReturnVal, dates, dailyTarget)
	if err != nil {
		return SortinoRatio{}, err
	}
	// Safeguard against Division by Zero
	if stdDeviationVal.AnnualizedVolatility == 0 {
		return SortinoRatio{}, errors.New("downside deviation is exactly zero; Sortino ratio is undefined")
	}
	annualizedMeanReturn := stdDeviationVal.Mean * 252.0
	// in sortino instead of risk free rate we use target value at both places it can be the risk free rate
	//itself but it is good practise to call it target or mar
	sortinoRatio := (annualizedMeanReturn - target) / stdDeviationVal.AnnualizedVolatility

	return SortinoRatio{
		SortinoRatioValue: sortinoRatio,
	}, nil
}
