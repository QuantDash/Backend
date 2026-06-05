-- name: GetDailyPrices :many
SELECT * FROM daily_prices;

-- name: GetPricesOnDate :many
SELECT adjusted_close FROM daily_prices WHERE trade_date 
BETWEEN ? AND ?; 
