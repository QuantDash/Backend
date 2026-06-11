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
}

func CalculateStandardDeviation(prices []float64, startDate, endDate time.Time) (StdDeviation, error) {
	if len(prices) < 3 {
		return StdDeviation{}, errors.New("cannot calculate Standard Deviation: requires at least 3 prices")
	}
	dailyReturns, err := CalculatePeriodicReturns(prices, startDate, endDate)
	if err != nil {
		return StdDeviation{}, err
	}
	fmt.Printf("Calculating Standard Deviation from: %v to %v\n", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	n := float64(len(dailyReturns.PerodicReturnVal))

	var totalSum float64
	for _, val := range dailyReturns.PerodicReturnVal {
		totalSum += val
	}
	mean := totalSum / n

	var variance float64
	for _, val := range dailyReturns.PerodicReturnVal {
		dif := val - mean
		variance += (dif * dif)
	}
	//(n-1 Bessel's correction)
	variance = variance / (n - 1)
	stdDev := math.Sqrt(variance)
	return StdDeviation{
		StandardDeviationValue: stdDev,
		AnnualizedVolatility:   (stdDev * math.Sqrt(252)), //annualized val
	}, nil
}
