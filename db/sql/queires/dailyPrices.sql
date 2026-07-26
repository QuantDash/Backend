-- name: GetDailyPricesPaginated :many
SELECT * FROM daily_prices
ORDER BY stock_id, trade_date
LIMIT ?
OFFSET ?;

-- name: GetStockIdsBySymbols :many
SELECT stock_id FROM symbols WHERE symbol IN (sqlc.slice(symbols));

-- name: GetDailyPricesByStockInDateRange :many
SELECT adjusted_close,trade_date FROM daily_prices
WHERE stock_id = ? AND trade_date BETWEEN ? AND ?
ORDER BY trade_date;


-- name: GetMultiStockPricesInDateRange :many
SELECT stock_id, adjusted_close, trade_date
FROM daily_prices
WHERE stock_id IN (sqlc.slice(stock_ids)) AND trade_date BETWEEN ? AND ?
ORDER BY stock_id, trade_date;