-- name: GetExplorerStocksPaginated :many
SELECT 
    s.symbol, 
    s.status, 
    dp.trade_date, 
    dp.close_price, 
    dp.adjusted_close, 
    dp.volume, 
    dp.adjusted_volume
FROM symbols s
JOIN daily_prices dp ON s.stock_id = dp.stock_id
INNER JOIN (
    SELECT stock_id, MAX(trade_date) as max_date
    FROM daily_prices
    GROUP BY stock_id
) latest ON dp.stock_id = latest.stock_id AND dp.trade_date = latest.max_date
ORDER BY s.symbol
LIMIT ? OFFSET ?;


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