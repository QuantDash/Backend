# QuantDash Backend

A REST API backend built in **Go** that powers the QuantDash quantitative analysis platform. It connects to a MySQL database of historical NSE (National Stock Exchange of India) price data and exposes endpoints for calculating institutional-grade financial metrics on any stock or date range.

---

## What It Does

The API accepts a stock symbol, a date range, and optional parameters (e.g., risk-free rate), then fetches the adjusted closing prices from the database and runs the requested calculation — returning clean, structured JSON.

All calculations are performed on **adjusted prices** (corrected for splits and dividends) to ensure historically accurate results.

---

## Tech Stack

| Layer | Technology |
|:---|:---|
| Language | Go |
| Router | [chi](https://github.com/go-chi/chi) |
| Database | MySQL |
| Query Layer | [sqlc](https://sqlc.dev/) — type-safe, generated SQL queries |
| DB Driver | `go-sql-driver/mysql` |
| Env Config | `godotenv` |
| Linear Algebra | `gonum` (used in Multi-Kelly covariance matrix) |

---

## API Endpoints

All routes are mounted under `/v1`.

| Method | Endpoint | Description |
|:---|:---|:---|
| `GET` | `/healthcheck` | Server liveness check |
| `GET` | `/v1/returns` | Periodic (daily) returns and cumulative return |
| `GET` | `/v1/dd` | Maximum Drawdown and Maximum Drawdown Duration |
| `GET` | `/v1/std` | Standard Deviation and Annualized Volatility |
| `GET` | `/v1/sharpe_ratio_value` | Sharpe Ratio |
| `GET` | `/v1/sortino_ratio_value` | Sortino Ratio |
| `GET` | `/v1/cagr_to_drawdown_ratio` | MAR Ratio (CAGR ÷ Max Drawdown) |
| `GET` | `/v1/4_sigma_threshold` | 4-Sigma outlier detection (Black Swan events) |
| `GET` | `/v1/kelly_formula` | Single-asset and Multi-asset Kelly Criterion |
| `GET` | `/v1/get_stock_details` | Paginated stock listing from the database |

## Calculations Implemented

### Returns
- **Periodic (Daily) Returns** — `(P[i] - P[i-1]) / P[i-1]`
- **Cumulative Return** — Total return over the full period
- **CAGR** — Compound Annual Growth Rate

### Risk Metrics
- **Standard Deviation** — Daily and annualized (`× √252`) volatility
- **Maximum Drawdown** — Peak-to-trough decline with exact peak and trough dates
- **Maximum Drawdown Duration** — Longest period spent below a prior peak (in calendar days)

### Risk-Adjusted Performance
- **Sharpe Ratio** — Annualized excess return per unit of annualized volatility
- **Sortino Ratio** — Like Sharpe, but only penalizes downside volatility
- **MAR Ratio** — CAGR divided by Maximum Drawdown

### Tail Risk
- **4-Sigma Threshold** — Identifies extreme daily return events (> 4 standard deviations from the mean). Returns the count, bounds, and the exact date and return value of each outlier. Useful for Black Swan analysis.

### Position Sizing
- **Single-Asset Kelly Criterion** — `f* = (μ - r) / σ²` — the mathematically optimal fraction of capital to allocate to one stock
- **Multi-Asset Kelly Criterion** — Solves `F = C⁻¹ × M` (inverse covariance matrix × mean excess returns vector) for optimal portfolio-level allocations across multiple assets simultaneously using `gonum`

---

## Getting Started

### Prerequisites
- Go 1.21+
- MySQL 8.0+
- A populated NSE price database

### 1. Clone the repository

```bash
git clone https://github.com/QuantDash/Backend.git
cd Backend
```

### 3. Install dependencies

```bash
go mod tidy
```

### 4. Run the server

```bash
go run .
```

The server starts on the port defined in your `.env`. Test it:

```bash
curl http://localhost:8080/healthcheck
```


## Architecture Notes

- **`calc/` is completely decoupled from HTTP and the database.** Every function in `calc/` takes plain Go slices (`[]float64`, `[]time.Time`) and returns a typed result struct. This means calculations can be tested in isolation without a running server or database.
- **`api/` is the glue layer.** Route handlers resolve the symbol → `stock_id`, fetch prices from the DB, then call the relevant `calc/` function and return the result as JSON.
- **`sqlc` generates all database code.** The raw SQL lives in `db/sql/queires/`. Running `sqlc generate` produces the type-safe Go functions in `db/internaldb/`. You never write raw SQL string operations in application code.
- **CORS is open** (`https://*`, `http://*`) — intended for local development. Restrict `AllowedOrigins` before any public deployment.

---
