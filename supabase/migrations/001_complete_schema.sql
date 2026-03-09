-- =====================================================
-- Share Market Simulator — Complete Backend Redesign
-- Fresh Supabase PostgreSQL Migration
-- All financial columns use NUMERIC(20,6)
-- =====================================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- =====================================================
-- DROP ALL EXISTING TABLES
-- =====================================================
DROP TABLE IF EXISTS market_index_history CASCADE;
DROP TABLE IF EXISTS price_triggers CASCADE;
DROP TABLE IF EXISTS trades CASCADE;
DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS portfolios CASCADE;
DROP TABLE IF EXISTS wallet_transfers CASCADE;
DROP TABLE IF EXISTS trading_wallets CASCADE;
DROP TABLE IF EXISTS main_wallets CASCADE;
DROP TABLE IF EXISTS ipo_applications CASCADE;
DROP TABLE IF EXISTS ipos CASCADE;
DROP TABLE IF EXISTS stock_prices CASCADE;
DROP TABLE IF EXISTS companies CASCADE;
DROP TABLE IF EXISTS stock_transactions CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS user_portfolios CASCADE;
DROP TABLE IF EXISTS virtual_wallets CASCADE;
DROP TABLE IF EXISTS market_events CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;
-- =====================================================
-- USERS
-- =====================================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    kyc_status VARCHAR(50) DEFAULT 'pending',
    role VARCHAR(50) DEFAULT 'user',
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    profile_image_id VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_users_email ON users(email);
-- =====================================================
-- COMPANIES
-- =====================================================
CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    symbol VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    sector VARCHAR(50) NOT NULL DEFAULT 'General',
    total_supply BIGINT NOT NULL,
    current_price NUMERIC(20, 6) NOT NULL DEFAULT 0,
    market_cap NUMERIC(20, 6) NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_companies_symbol ON companies(symbol);
CREATE INDEX idx_companies_sector ON companies(sector);
CREATE INDEX idx_companies_is_active ON companies(is_active);
-- =====================================================
-- IPOS
-- =====================================================
CREATE TABLE ipos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    price_per_share NUMERIC(20, 6) NOT NULL,
    total_shares BIGINT NOT NULL,
    allocated_shares BIGINT NOT NULL DEFAULT 0,
    max_per_applicant BIGINT NOT NULL DEFAULT 100,
    open_at TIMESTAMPTZ NOT NULL,
    close_at TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    -- pending, open, closed, allocated
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_ipos_company_id ON ipos(company_id);
CREATE INDEX idx_ipos_status ON ipos(status);
-- =====================================================
-- IPO_APPLICATIONS
-- =====================================================
CREATE TABLE ipo_applications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ipo_id UUID NOT NULL REFERENCES ipos(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    shares_requested BIGINT NOT NULL,
    shares_allocated BIGINT NOT NULL DEFAULT 0,
    amount_paid NUMERIC(20, 6) NOT NULL DEFAULT 0,
    amount_refunded NUMERIC(20, 6) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    -- pending, allocated, not_allocated, refunded
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(ipo_id, user_id)
);
CREATE INDEX idx_ipo_applications_ipo_id ON ipo_applications(ipo_id);
CREATE INDEX idx_ipo_applications_user_id ON ipo_applications(user_id);
-- =====================================================
-- MAIN_WALLETS (deposits, withdrawals, real-money)
-- =====================================================
CREATE TABLE main_wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    balance NUMERIC(20, 6) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_main_wallets_user_id ON main_wallets(user_id);
-- =====================================================
-- TRADING_WALLETS (buying shares, receiving sale proceeds)
-- =====================================================
CREATE TABLE trading_wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    balance NUMERIC(20, 6) NOT NULL DEFAULT 0,
    locked_balance NUMERIC(20, 6) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_trading_wallets_user_id ON trading_wallets(user_id);
-- =====================================================
-- WALLET_TRANSFERS (audit log for main ↔ trading)
-- =====================================================
CREATE TABLE wallet_transfers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount NUMERIC(20, 6) NOT NULL,
    direction VARCHAR(20) NOT NULL,
    -- main_to_trading, trading_to_main
    status VARCHAR(20) NOT NULL DEFAULT 'completed',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_wallet_transfers_user_id ON wallet_transfers(user_id);
-- =====================================================
-- PORTFOLIOS (user share holdings)
-- =====================================================
CREATE TABLE portfolios (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    quantity BIGINT NOT NULL DEFAULT 0,
    avg_buy_price NUMERIC(20, 6) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, company_id)
);
CREATE INDEX idx_portfolios_user_id ON portfolios(user_id);
CREATE INDEX idx_portfolios_user_company ON portfolios(user_id, company_id);
-- =====================================================
-- ORDERS (order book: buy/sell limit/market orders)
-- =====================================================
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    side VARCHAR(10) NOT NULL,
    -- buy, sell
    order_type VARCHAR(10) NOT NULL DEFAULT 'limit',
    -- limit, market
    price NUMERIC(20, 6) NOT NULL DEFAULT 0,
    quantity BIGINT NOT NULL,
    filled_qty BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'open',
    -- open, partially_filled, filled, cancelled
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_orders_company_side ON orders(company_id, side);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_company_side_price ON orders(company_id, side, price, created_at);
-- =====================================================
-- TRADES (matched order executions)
-- =====================================================
CREATE TABLE trades (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    buy_order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    sell_order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    buyer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    seller_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    price NUMERIC(20, 6) NOT NULL,
    quantity BIGINT NOT NULL,
    total_amount NUMERIC(20, 6) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_trades_company_id ON trades(company_id);
CREATE INDEX idx_trades_buyer_id ON trades(buyer_id);
CREATE INDEX idx_trades_seller_id ON trades(seller_id);
CREATE INDEX idx_trades_created_at ON trades(created_at DESC);
-- =====================================================
-- STOCK_PRICES (OHLCV candle data)
-- =====================================================
CREATE TABLE stock_prices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    open_price NUMERIC(20, 6) NOT NULL,
    high_price NUMERIC(20, 6) NOT NULL,
    low_price NUMERIC(20, 6) NOT NULL,
    close_price NUMERIC(20, 6) NOT NULL,
    volume BIGINT NOT NULL DEFAULT 0,
    timestamp TIMESTAMPTZ NOT NULL,
    timeframe VARCHAR(10) NOT NULL DEFAULT '1D',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_stock_prices_company_time ON stock_prices(company_id, timestamp DESC);
CREATE INDEX idx_stock_prices_company_timeframe ON stock_prices(company_id, timeframe, timestamp DESC);
-- =====================================================
-- PRICE_TRIGGERS (auto-sell when price hits target)
-- =====================================================
CREATE TABLE price_triggers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    trigger_price NUMERIC(20, 6) NOT NULL,
    shares_qty BIGINT NOT NULL,
    direction VARCHAR(10) NOT NULL DEFAULT 'above',
    -- above, below
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    -- active, triggered, cancelled
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_price_triggers_user_id ON price_triggers(user_id);
CREATE INDEX idx_price_triggers_company_status ON price_triggers(company_id, status);
-- =====================================================
-- MARKET_INDEX_HISTORY
-- =====================================================
CREATE TABLE market_index_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    index_value NUMERIC(20, 6) NOT NULL,
    change NUMERIC(20, 6) NOT NULL DEFAULT 0,
    change_percent NUMERIC(20, 6) NOT NULL DEFAULT 0,
    total_turnover NUMERIC(20, 6) NOT NULL DEFAULT 0,
    total_volume BIGINT NOT NULL DEFAULT 0,
    total_market_cap NUMERIC(20, 6) NOT NULL DEFAULT 0,
    advances INTEGER NOT NULL DEFAULT 0,
    declines INTEGER NOT NULL DEFAULT 0,
    unchanged INTEGER NOT NULL DEFAULT 0,
    total_companies INTEGER NOT NULL DEFAULT 0,
    previous_close NUMERIC(20, 6) NOT NULL DEFAULT 0,
    timestamp TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_market_index_history_timestamp ON market_index_history(timestamp DESC);
-- =====================================================
-- AUTO-UPDATE TRIGGER
-- =====================================================
CREATE OR REPLACE FUNCTION update_updated_at_column() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
CREATE TRIGGER update_users_updated_at BEFORE
UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_companies_updated_at BEFORE
UPDATE ON companies FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_ipos_updated_at BEFORE
UPDATE ON ipos FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_ipo_applications_updated_at BEFORE
UPDATE ON ipo_applications FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_main_wallets_updated_at BEFORE
UPDATE ON main_wallets FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_trading_wallets_updated_at BEFORE
UPDATE ON trading_wallets FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_portfolios_updated_at BEFORE
UPDATE ON portfolios FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_orders_updated_at BEFORE
UPDATE ON orders FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_stock_prices_updated_at BEFORE
UPDATE ON stock_prices FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_price_triggers_updated_at BEFORE
UPDATE ON price_triggers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- =====================================================
-- ROW LEVEL SECURITY
-- =====================================================
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE companies ENABLE ROW LEVEL SECURITY;
ALTER TABLE ipos ENABLE ROW LEVEL SECURITY;
ALTER TABLE ipo_applications ENABLE ROW LEVEL SECURITY;
ALTER TABLE main_wallets ENABLE ROW LEVEL SECURITY;
ALTER TABLE trading_wallets ENABLE ROW LEVEL SECURITY;
ALTER TABLE wallet_transfers ENABLE ROW LEVEL SECURITY;
ALTER TABLE portfolios ENABLE ROW LEVEL SECURITY;
ALTER TABLE orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE trades ENABLE ROW LEVEL SECURITY;
ALTER TABLE stock_prices ENABLE ROW LEVEL SECURITY;
ALTER TABLE price_triggers ENABLE ROW LEVEL SECURITY;
ALTER TABLE market_index_history ENABLE ROW LEVEL SECURITY;
-- Public read on market data
CREATE POLICY "public_read_companies" ON companies FOR
SELECT USING (true);
CREATE POLICY "public_read_stock_prices" ON stock_prices FOR
SELECT USING (true);
CREATE POLICY "public_read_ipos" ON ipos FOR
SELECT USING (true);
CREATE POLICY "public_read_orders" ON orders FOR
SELECT USING (true);
CREATE POLICY "public_read_trades" ON trades FOR
SELECT USING (true);
CREATE POLICY "public_read_market_index" ON market_index_history FOR
SELECT USING (true);
-- Service role full access (backend uses service role key)
CREATE POLICY "all_users" ON users FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "all_companies" ON companies FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "all_ipos" ON ipos FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "all_ipo_applications" ON ipo_applications FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "all_main_wallets" ON main_wallets FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "all_trading_wallets" ON trading_wallets FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "all_wallet_transfers" ON wallet_transfers FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "all_portfolios" ON portfolios FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "all_orders" ON orders FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "all_trades" ON trades FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "all_stock_prices" ON stock_prices FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "all_price_triggers" ON price_triggers FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "all_market_index" ON market_index_history FOR ALL USING (true) WITH CHECK (true);