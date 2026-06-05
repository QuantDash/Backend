
-- Symbols Table
CREATE TABLE IF NOT EXISTS symbols (
    stock_id INT AUTO_INCREMENT PRIMARY KEY,
    symbol VARCHAR(255) UNIQUE NOT NULL,
    company_name VARCHAR(255),
    status VARCHAR(10) DEFAULT 'Active'
);

-- Daily Prices Table
CREATE TABLE IF NOT EXISTS daily_prices (
    stock_id INT NOT NULL,
    trade_date DATE NOT NULL,
    open_price DECIMAL(10, 2),
    high_price DECIMAL(10, 2),
    low_price DECIMAL(10, 2),
    close_price DECIMAL(10, 2),
    adjusted_open DECIMAL(10, 2),
    adjusted_high DECIMAL(10, 2),
    adjusted_low DECIMAL(10, 2),
    adjusted_close DECIMAL(10, 2),
    volume BIGINT,
    adjusted_volume BIGINT,
    PRIMARY KEY (stock_id, trade_date),
    FOREIGN KEY (stock_id) REFERENCES symbols(stock_id)
);

-- Corporate_actions 
CREATE TABLE IF NOT EXISTS corporate_actions (
    stock_id INT NOT NULL,
    ex_date DATE NOT NULL,
    action_type VARCHAR(20) NOT NULL,
    dividend_amount DECIMAL(10, 2) DEFAULT 0.00,
    split_ratio DECIMAL(10, 4) DEFAULT 1.0000,
    PRIMARY KEY (stock_id, ex_date, action_type),
    FOREIGN KEY (stock_id) REFERENCES symbols(stock_id)
);

