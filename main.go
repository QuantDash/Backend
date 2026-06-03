package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Print("Check Mannan")
	godotenv.Load() //to load the env file we use this function

	//to use env go has the follwing functoin
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the env")
	}
	mainRouter := chi.NewRouter()
	mainRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	mainRouter.Get("/", HealthCheck)
	server := &http.Server{
		Addr:    ":" + portString,
		Handler: mainRouter,
	}
	fmt.Println("Starting server on PORT:", portString)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Health Check OK!"))
}

/*
Matrix to implement:
1. Sharp Ratio -It calculates your return divided by your volatility (risk).
1.1 annualize Sharpe Ratio
2. Information Ratio: Similar to sharp ratio
3. Drawdown: How long u endure no profit
4. Transsaction Cost : Cost or tax you have to pay to brokers site (zerodha)
5. Basis Points :
6.  The 4-Sigma Threshold : Standard Deviation of stocks
7. MAR ratio
8. CAGR: we will calculate this but this involves returns that is considerd not good metrix
so we will add a flag about it that why cagr is bad
9 kellys formula
10 something related to bellcurve and gaussian formula

From Gemini:
1. Base Returns

Daily Returns: The raw day-to-day percentage changes.

Cumulative Returns: Total growth over the selected time period.

CAGR (with flag): As you smartly noted, include the flag explaining that CAGR hides intra-year volatility.

2. Risk & Distribution (The Math Heavy Section)

Standard Deviation (Annualized): The raw volatility of the stock.

Gaussian/Bell Curve Distribution: To show the user if the stock's returns are normally distributed.

The 4-Sigma Threshold: To highlight "Black Swan" events (days the stock moved so violently it broke normal statistical probability).

3. Drawdown (Pain Tolerance)

Max Drawdown: The absolute worst peak-to-trough drop in the selected timeframe.

Drawdown Duration: How many days the stock took to recover to its previous high.

4. Risk-Adjusted Scoring (The "Quality" Scores)

Sharpe Ratio (Annualized): Return divided by total volatility.

Sortino Ratio (Annualized): Return divided by downside volatility.

Information Ratio: How well the stock performed compared to a benchmark (e.g., Nifty 50).

MAR Ratio: CAGR divided by Max Drawdown.

5. Real-World Adjustments

Transaction Costs (Zerodha / STT): Simulating what happens to the return if the user bought/sold.

Kelly's Formula: Telling the user what percentage of their total portfolio they should safely allocate to this specific stock based on its historical math.



*/
