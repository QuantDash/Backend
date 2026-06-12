package calc

import (
	"errors"
	"fmt"
	"math"
	"time"
)

type CumulativeReturn struct {
	CumulativeReturnVal float64 `json:"cumulative_return_value"`
}

type PeriodicReturn struct {
	PeriodicReturnVal []float64 `json:"periodic_return_value"`
}

type CAGR struct {
	CAGRValue float64 `json:"cagr_value"`
}

func CalculateCumulativeReturns(prices []float64, dates []time.Time) (CumulativeReturn, error) {
	if len(prices) < 2 {
		return CumulativeReturn{}, errors.New("cannot calculate returns: requires at least 2 prices")
	}
	fmt.Printf("Calculating Cumulative Return from: %v to %v\n", dates[0].Format("2006-01-02"), dates[len(dates)-1].Format("2006-01-02"))
	cumuReturn := (prices[len(prices)-1] - prices[0]) / prices[0]
	return CumulativeReturn{
		CumulativeReturnVal: cumuReturn,
	}, nil
}

func CalculatePeriodicReturns(prices []float64, dates []time.Time) (PeriodicReturn, error) {
	if len(prices) < 2 {
		return PeriodicReturn{}, errors.New("cannot calculate periodic returns: requires at least 2 prices")
	}
	fmt.Printf("Calculating periodic Return from: %v to %v\n", dates[0].Format("2006-01-02"), dates[len(dates)-1].Format("2006-01-02"))
	returns := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		if prices[i-1] == 0.0 {
			return PeriodicReturn{}, fmt.Errorf("invalid price data: price at index %d is 0.0, causing division by zero", i-1)
		}

		returns[i-1] = (prices[i] - prices[i-1]) / prices[i-1]
	}
	//values are in decimal not yet x by 100, as needed for math in future
	return PeriodicReturn{
		PeriodicReturnVal: returns,
	}, nil
}

// read docs understand what CAGR is.
func CalculateCAGR(prices []float64, dates []time.Time) (CAGR, error) {
	if len(prices) < 2 {
		return CAGR{}, errors.New("cannot calculate returns: requires at least 2 prices")
	}
	fmt.Printf("Calculating CAGR from: %v to %v\n", dates[0].Format("2006-01-02"), dates[len(dates)-1].Format("2006-01-02"))
	//.Sub() returns time.Duration.a time.Duration is just a single, massive integer representing the exact
	// number of nanoseconds between enddate and startdate.
	totalDuration := dates[len(dates)-1].Sub(dates[0])

	totalDays := totalDuration.Hours() / 24.0
	//here we use 365 and not 365.25 as 365 is the minima if greater then that we can execute it anyways
	if totalDays < 365 {
		return CAGR{}, errors.New("CAGR cannot be calculated for less then 1 Year")
	}

	//If the time duration is exactly 365.0 days, 365.0 / 365.2425 equals 0.99933 years. When the formula calculates the exponent 1 / 0.99933, it becomes 1.0006 instead of 1.
	// It is mathematically accepted in quantitative finance because we are strictly normalizing to the astronomical solar year rather than a flat calendar year.
	totalYears := totalDays / 365.2425
	//That (1/totalYears) exponent is the Annualizer similar to std deviation we have root 252 for annualizer
	cagr := math.Pow(prices[len(prices)-1]/prices[0], (1/totalYears)) - 1
	return CAGR{
		CAGRValue: cagr,
	}, nil
}
