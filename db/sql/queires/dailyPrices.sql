-- name: GetDailyPricesPaginated :many
SELECT * FROM daily_prices
ORDER BY stock_id, trade_date
LIMIT ?
OFFSET ?;

-- name: GetDailyPricesByStockInDateRange :many
SELECT adjusted_close,trade_date FROM daily_prices
WHERE stock_id = ? AND trade_date BETWEEN ? AND ?
ORDER BY trade_date;
