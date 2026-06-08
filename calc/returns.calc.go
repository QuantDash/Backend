package calc

import (
	"errors"
	"fmt"
)

func CalculateCumulativeReturns(prices []float64, fromDate, toDate string) (float64, error) {
	if len(prices) < 2 {
		return 0.0, errors.New("cannot calculate returns: requires at least 2 prices")
	}
	fmt.Printf("Calculating Cumulative Return from: %v to %v\n", fromDate, toDate)
	cumulativeReturn := (prices[len(prices)-1] - prices[0]) / prices[0]
	return cumulativeReturn, nil
}

func CalculatePeriodicReturns(prices []float64, fromDate, toDate string) ([]float64, error) {
	if len(prices) < 2 {
		return nil, errors.New("cannot calculate returns: requires at least 2 prices")
	}
	fmt.Printf("Calculating Perodic Return from: %v to %v\n", fromDate, toDate)
	returns := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		if prices[i-1] == 0.0 {
			return nil, fmt.Errorf("invalid price data: price at index %d is 0.0, causing division by zero", i-1)
		}

		returns[i-1] = (prices[i] - prices[i-1]) / prices[i-1]
	}
	//values are in decimal not yet x by 100, as needed for math in future
	return returns, nil

}
