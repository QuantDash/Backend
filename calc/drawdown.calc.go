package calc

import (
	"errors"
	"fmt"
	"time"
)

/*
What to use here raw prices or dailyreturns:
When not to use raw prices and use dailyreturns
- Raw prices have splits/dividends/unadjusted data → a cumulative return index is cleaner
- You're combining multiple assets and need time-weighted returns
When raw prices are fine:
- You already have adjusted prices (most backtesting/API data)
- Less code, fewer edge cases, same result
- The function signature stays simple
For our codebase, raw prices is the right choice. CalculatePeriodicReturns exists for things that genuinely need period returns (like volatility calculation, Sharpe ratio, etc.), not for drawdown.
*/

type MaxDrawDownResult struct {
	MaxDrawDown float64    `json:"max_drawdown"`
	TroughDate  *time.Time `json:"trough_date"`
}

type MaxDrawDownDurationResult struct {
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	DurationDays int64      `json:"duration_days"`
}

func CalculateMaximumDrawDown(prices []float64, dates []time.Time) (MaxDrawDownResult, error) {
	//check
	if len(prices) < 2 {
		return MaxDrawDownResult{}, errors.New("cannot calculate Maximum Drawdown: requires at least 2 prices")
	}
	if len(prices) != len(dates) {
		return MaxDrawDownResult{}, errors.New("the length of dates and prices does not match")
	}

	fmt.Printf("Calculating MaximumDrawDown from %v to %v\n", dates[0].Format("2006-01-02"), dates[len(dates)-1].Format("2006-01-02"))
	//trackers
	peakPrice := prices[0]
	//pointer as if we dont get drawdown then we need null in json
	var ddDate *time.Time
	maxDrawDown := 0.0

	for i := 1; i < len(prices); i++ {
		if prices[i] > peakPrice {
			peakPrice = prices[i]
		} else {
			drawdown := (prices[i] - peakPrice) / peakPrice
			if drawdown < maxDrawDown {
				maxDrawDown = drawdown
				ddDate = &dates[i]
			}
		}
	}
	return MaxDrawDownResult{
		MaxDrawDown: maxDrawDown,
		TroughDate:  ddDate,
	}, nil
}
func CalculateMaximumDrawDownDuration(prices []float64, dates []time.Time) (MaxDrawDownDurationResult, error) {
	//check
	if len(prices) < 2 {
		return MaxDrawDownDurationResult{}, errors.New("cannot calculate MaxDrawdownDuration: requires at least 2 prices")
	}
	if len(prices) != len(dates) {
		return MaxDrawDownDurationResult{}, errors.New("The length of dates and prices does not match")
	}
	fmt.Printf("Calculating MaximumDrawDownDuration from %v to %v\n", dates[0].Format("2006-01-02"), dates[len(dates)-1].Format("2006-01-02"))

	//trackers
	peakPrice := prices[0]
	peakIndex := 0
	maxDuration := time.Duration(0)

	// read point 14 in about.txt for why we need wasUnderWater
	wasUnderWater := false

	var drawdownStartDate *time.Time // start of the *current* underwater period
	var lastStartDate *time.Time
	var lastEndDate *time.Time
	for i := 1; i < len(prices); i++ {
		//returns time.Duration in nanoseconds
		drawdownDuration := dates[i].Sub(dates[peakIndex])

		if prices[i] >= peakPrice {
			if wasUnderWater {
				if drawdownDuration > maxDuration {
					lastStartDate = drawdownStartDate
					lastEndDate = &dates[i]
					maxDuration = drawdownDuration
				}
				wasUnderWater = false
			}
			peakPrice = prices[i]
			peakIndex = i
		} else {
			if !wasUnderWater {
				// just went underwater — record the start
				wasUnderWater = true
				drawdownStartDate = &dates[peakIndex]
			}
			if drawdownDuration > maxDuration {
				lastStartDate = drawdownStartDate
				lastEndDate = &dates[i]
				maxDuration = drawdownDuration
			}
		}
	}
	return MaxDrawDownDurationResult{
		StartDate:    lastStartDate,
		EndDate:      lastEndDate,
		DurationDays: int64(maxDuration.Hours() / 24.0),
	}, nil
}
