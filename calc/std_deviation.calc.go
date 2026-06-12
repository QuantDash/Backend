package calc

import (
	"errors"
	"fmt"
	"math"
	"time"
)

// read docs understand what std_deviation is.
// basic defination :
// The standard deviation (σ) is a measure that is used to quantify the amount of variation or dispersion of data from its mean
// basically how much a stock variate on day to day basis so you dont have to worry to come and check it everytime

// Standard Deviation and Daily Volatility are the exact same thing
type StdDeviation struct {
	StandardDeviationValue float64 `json:"std_deviation_value"`   // This represents Daily Volatility
	AnnualizedVolatility   float64 `json:"annualized_volatility"` // it is StandardDeviationValue/daily volatility * root 252
	Mean                   float64 `json:"mean"`
	TotalTradingDays       int64   `json:"total_trading_days"`
	TotalCalenderDays      int64   `json:"total_calender_days"`
}

func CalculateStandardDeviation(dailyReturns []float64, dates []time.Time) (StdDeviation, error) {
	if len(dailyReturns) < 2 {
		return StdDeviation{}, errors.New("cannot calculate Standard Deviation: requires at least 2 returns")
	}
	fmt.Printf("Calculating Standard Deviation from: %v to %v\n", dates[0].Format("2006-01-02"), dates[len(dates)-1].Format("2006-01-02"))
	n := float64(len(dailyReturns))

	var totalSum float64
	for _, val := range dailyReturns {
		totalSum += val
	}
	mean := totalSum / n

	var variance float64
	for _, val := range dailyReturns {
		dif := val - mean
		variance += (dif * dif)
	}
	//(n-1 Bessel's correction)
	variance = variance / (n - 1)
	stdDev := math.Sqrt(variance)
	return StdDeviation{
		StandardDeviationValue: stdDev,
		AnnualizedVolatility:   (stdDev * math.Sqrt(252)), //annualized val
		Mean:                   mean,
		TotalTradingDays:       int64(len(dates)),
		TotalCalenderDays:      int64(dates[len(dates)-1].Sub(dates[0]).Hours() / 24.0),
	}, nil
}
