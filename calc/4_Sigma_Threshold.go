package calc

import (
	"errors"
	"fmt"
	"time"
)

type Outlier struct {
	Date        time.Time `json:"date"`
	ReturnValue float64   `json:"return_value"`
}

// read about 4 sigma Threshold in docs
type SigmaThreshold struct {
	StandardDeviation StdDeviation `json:"standard_deviation"`
	OutlierCount      int64        `json:"4_sigma_threshold_count"`
	UpperBound        float64      `json:"upper_bound"`
	LowerBound        float64      `json:"lower_bound"`
	Outliers          []Outlier    `json:"outliers"`
}

func Calculate4SigmaThreshold(prices []float64, dates []time.Time) (SigmaThreshold, error) {
	fmt.Printf("Calculating 4 Sigma Threshold from: %v to %v\n", dates[0].Format("2006-01-02"), dates[len(dates)-1].Format("2006-01-02"))
	if len(dates) != len(prices) {
		return SigmaThreshold{}, errors.New("The length of dates and prices does not match")
	}
	if len(prices) < 3 {
		return SigmaThreshold{}, errors.New("cannot calculate 4-Sigma threshold: requires at least 3 prices")
	}
	dailyReturns, err := CalculatePeriodicReturns(prices, dates)
	if err != nil {
		return SigmaThreshold{}, err
	}
	stdDev, err := CalculateStandardDeviation(dailyReturns.PeriodicReturnVal, dates)
	if err != nil {
		return SigmaThreshold{}, err
	}
	//as i am evaluating a single day's behavior, the mathematical boundaries you build must also be based on a single day's risk.
	// exampple: std D = 2 percent
	//upperBound = 0.00 + (4 * 0.02) = +0.08 (+8%)
	// lowerBound = 0.00 - (4 * 0.02) = -0.08 (-8%)
	// The Rule: Any daily return between -8% and +8% is considered "normal enough." Anything beyond that is a 4-sigma outlier.
	upperBound := stdDev.Mean + (4 * stdDev.StandardDeviationValue)
	lowerBound := stdDev.Mean - (4 * stdDev.StandardDeviationValue)

	counter := 0
	outliers := []Outlier{}
	for i, val := range dailyReturns.PeriodicReturnVal {
		if val < lowerBound || val > upperBound {
			counter++
			outliers = append(outliers, Outlier{
				Date:        dates[i+1], // i+1 as daily returns are calculated from i-1 to i
				ReturnValue: val,
			})
		}
	}
	return SigmaThreshold{
		StandardDeviation: stdDev,
		OutlierCount:      int64(counter),
		UpperBound:        upperBound,
		LowerBound:        lowerBound,
		Outliers:          outliers,
	}, nil
}
