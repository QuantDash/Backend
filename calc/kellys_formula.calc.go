package calc

import (
	"errors"
	"fmt"
	"math"
	"time"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

type KellyResult struct {
	KellyPercentage               float64 `json:"kelly_percentage"`
	HalfKellyPercentage           float64 `json:"half_kelly_percentage"`
	OptimalGrowthRateForFullKelly float64 `json:"optimal_growth_rate_for_full_kelly"`
	OptimalGrowthRateForHalfKelly float64 `json:"optimal_growth_rate_for_half_kelly"`
	AnnualizedMeanReturn          float64 `json:"annualized_mean_return"`
	AnnualizedVolatility          float64 `json:"annualized_volatility"`
	ExcessReturn                  float64 `json:"excess_return"`
	SharpeRatio                   float64 `json:"sharpe_ratio"`
}

type MultiKellyResult struct {
	OptimalAllocations []float64 `json:"optimal_allocations"`
}

func CalculateSingleKellyPercentage(prices []float64, dates []time.Time, riskFreeRate float64) (KellyResult, error) {
	// Use a similar validation as Sharpe/Sortino for data robustness.
	if len(prices) < 64 {
		return KellyResult{}, errors.New("cannot calculate Kelly Criterion for less than 3 months of data")
	}

	fmt.Printf("Calculating Kelly Criterion from: %v to %v\n", dates[0].Format("2006-01-02"), dates[len(dates)-1].Format("2006-01-02"))

	// 1. Get daily returns.
	dailyReturns, err := CalculatePeriodicReturns(prices, dates)
	if err != nil {
		return KellyResult{}, err
	}
	stdDev, err := CalculateStandardDeviation(dailyReturns.PeriodicReturnVal, dates)
	if err != nil {
		return KellyResult{}, err
	}

	if stdDev.AnnualizedVolatility == 0 {
		return KellyResult{}, errors.New("annualized volatility is zero; Kelly Criterion is undefined")
	}

	//Calculate the components for the formulas.
	annualizedMeanReturn := stdDev.Mean * 252.0
	// one thing to note here is we cannot predict the right return we get from this stock,
	// so what we do is use the past historical return as a guess for the future expected return.
	excessReturn := annualizedMeanReturn - riskFreeRate            // This is 'm'
	annualizedVariance := math.Pow(stdDev.AnnualizedVolatility, 2) // This is 's²'

	// Calculate the Kelly Percentage (f*) = expected return from that stock / that stocks volatility square
	kellyPercentage := excessReturn / annualizedVariance

	// Calculate the Optimal Growth Rate (g) for Full Kelly.
	// The formula g = r + (S^2 / 2) is the growth rate at the optimal fraction f* (Full Kelly).
	sharpeRatio := (annualizedMeanReturn - riskFreeRate) / stdDev.AnnualizedVolatility // This is 'S'
	//read point 20. in about.txt on why did we add riskfree rate in growtrate formula
	optimalGrowthRateForFullKelly := riskFreeRate + (math.Pow(sharpeRatio, 2) / 2.0)

	// Important thing to note here is we cannot do optimalGrowthRateForFullKelly*.75 because
	//25% reduction in growth applies only to the excess return you get for taking risk, not to the entire growth rate as it also inovles adding riskfreerate.
	// Calculate the Growth Rate for Half Kelly.
	// The excess growth for Half Kelly is 75% of the excess growth for Full Kelly.
	// Excess Growth (Full) = g_full - r
	// Excess Growth (Half) = (g_full - r) * 0.75
	// g_half = r + (g_full - r) * 0.75
	optimalGrowthRateForHalfKelly := riskFreeRate + (optimalGrowthRateForFullKelly-riskFreeRate)*0.75

	return KellyResult{
		KellyPercentage:               kellyPercentage,
		HalfKellyPercentage:           kellyPercentage / 2.0,
		OptimalGrowthRateForFullKelly: optimalGrowthRateForFullKelly,
		OptimalGrowthRateForHalfKelly: optimalGrowthRateForHalfKelly,
		AnnualizedMeanReturn:          annualizedMeanReturn,
		AnnualizedVolatility:          stdDev.AnnualizedVolatility,
		ExcessReturn:                  excessReturn,
		SharpeRatio:                   sharpeRatio,
	}, nil
}

// CalculateMultiKellyPercentage computes the optimal allocation for a portfolio of assets.
// read about.txt point 23 and 24
func CalculateMultiKellyPercentage(priceSeries [][]float64, riskFreeRate float64) (MultiKellyResult, error) {
	if len(priceSeries) == 0 {
		return MultiKellyResult{}, errors.New("priceSeries cannot be empty")
	}

	numAssets := len(priceSeries)
	if numAssets == 0 {
		return MultiKellyResult{}, errors.New("no assets provided")
	}

	// Validate that all assets have the same number of price points and enough data.
	numObservations := len(priceSeries[0])
	if numObservations < 64 {
		return MultiKellyResult{}, errors.New("cannot calculate Kelly Criterion for less than 3 months of data")
	}
	for i, series := range priceSeries {
		if len(series) != numObservations {
			return MultiKellyResult{}, fmt.Errorf("asset %d has a different number of observations", i)
		}
	}

	// 1. Calculate daily returns for each asset.
	// The returns matrix will have (numObservations - 1) rows and numAssets columns.
	numReturns := numObservations - 1
	returnsData := make([]float64, numReturns*numAssets)
	for i, assetPrices := range priceSeries {
		// We don't need the dates for this internal calculation.
		returns, err := CalculatePeriodicReturns(assetPrices, nil)
		if err != nil {
			return MultiKellyResult{}, fmt.Errorf("error calculating returns for asset %d: %w", i, err)
		}
		// Place the returns for this asset into the correct column.
		for j, ret := range returns.PeriodicReturnVal {
			returnsData[j*numAssets+i] = ret
		}
	}
	returnsMatrix := mat.NewDense(numReturns, numAssets, returnsData)

	// 2. Calculate the vector of mean excess returns (M).
	meanExcessReturns := mat.NewVecDense(numAssets, nil)
	dailyRiskFreeRate := riskFreeRate / 252.0
	for i := 0; i < numAssets; i++ {
		col := mat.Col(nil, i, returnsMatrix)
		meanReturn := stat.Mean(col, nil)
		meanExcessReturn := meanReturn - dailyRiskFreeRate
		meanExcessReturns.SetVec(i, meanExcessReturn)
	}
	// Annualize the mean excess returns vector.
	meanExcessReturns.ScaleVec(252.0, meanExcessReturns)

	// 3. Calculate the covariance matrix of daily returns (C).
	var cov mat.SymDense
	stat.CovarianceMatrix(&cov, returnsMatrix, nil)

	// 4. Annualize the covariance matrix.
	// C_annual = C_daily * 252
	cov.ScaleSym(252.0, &cov)

	// 5. Invert the covariance matrix (C⁻¹).
	var invCov mat.Dense
	err := invCov.Inverse(&cov)
	if err != nil {
		// This often happens if assets are perfectly correlated or there isn't enough data.
		return MultiKellyResult{}, fmt.Errorf("failed to invert covariance matrix: %w. Check for highly correlated assets or insufficient data", err)
	}

	// 6. Calculate the optimal allocations vector (F = C⁻¹ * M).
	var allocationsVec mat.VecDense
	allocationsVec.MulVec(&invCov, meanExcessReturns)

	// 7. Convert the result to a simple slice.
	optimalAllocations := make([]float64, numAssets)
	for i := 0; i < numAssets; i++ {
		optimalAllocations[i] = allocationsVec.AtVec(i)
	}

	return MultiKellyResult{
		OptimalAllocations: optimalAllocations,
	}, nil
}
