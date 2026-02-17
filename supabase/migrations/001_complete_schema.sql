-- =====================================================
-- Supabase PostgreSQL Migration: Complete Stock Market Simulator Schema
-- This is a COMPLETE fresh migration - drops and recreates everything
-- Run this ONCE for a clean database setup
-- =====================================================

-- Enable UUID extension (usually enabled by default in Supabase)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =====================================================
-- DROP ALL TABLES (START FRESH)
-- Drop tables first (CASCADE automatically removes triggers, indexes, policies)
-- Then drop the function
-- =====================================================
DROP TABLE IF EXISTS stock_transactions CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS user_portfolios CASCADE;
DROP TABLE IF EXISTS virtual_wallets CASCADE;
DROP TABLE IF EXISTS market_events CASCADE;
DROP TABLE IF EXISTS stock_prices CASCADE;
DROP TABLE IF EXISTS companies CASCADE;
DROP TABLE IF EXISTS users CASCADE;

DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;

-- =====================================================
-- USERS TABLE
-- Purpose: Store user accounts with authentication
-- =====================================================
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name       VARCHAR(100)  NOT NULL,
    email           VARCHAR(255)  NOT NULL UNIQUE,
    phone           VARCHAR(20)   NOT NULL,
    password_hash   VARCHAR(255)  NOT NULL,
    kyc_status      VARCHAR(50)   DEFAULT 'pending',     -- pending, approved, rejected
    role            VARCHAR(50)   DEFAULT 'user',        -- user, admin
    profile_image_id VARCHAR(100),
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);

-- =====================================================
-- COMPANIES TABLE
-- Purpose: Store stock market company information
-- =====================================================
CREATE TABLE companies (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    symbol           VARCHAR(10)   NOT NULL UNIQUE,
    name             VARCHAR(100)  NOT NULL,
    sector           VARCHAR(50)   NOT NULL,
    market_cap       DOUBLE PRECISION NOT NULL,
    description      TEXT,
    founded_year     INTEGER,
    employees        INTEGER,
    total_shares     BIGINT        NOT NULL DEFAULT 10000000,  -- Total shares issued
    available_shares BIGINT        NOT NULL DEFAULT 10000000,  -- Shares available for trading
    is_active        BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_companies_symbol ON companies(symbol);
CREATE INDEX idx_companies_sector ON companies(sector);
CREATE INDEX idx_companies_is_active ON companies(is_active);
CREATE INDEX idx_companies_available_shares ON companies(available_shares);

-- =====================================================
-- STOCK_PRICES TABLE
-- Purpose: Store historical stock price data
-- =====================================================
CREATE TABLE stock_prices (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id  UUID          NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    open_price  DOUBLE PRECISION NOT NULL,
    high_price  DOUBLE PRECISION NOT NULL,
    low_price   DOUBLE PRECISION NOT NULL,
    close_price DOUBLE PRECISION NOT NULL,
    volume      BIGINT        NOT NULL DEFAULT 0,
    timestamp   TIMESTAMPTZ   NOT NULL,
    timeframe   VARCHAR(10)   NOT NULL DEFAULT '1D',      -- 1D, 1W, 1M, etc
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_stock_prices_company_time ON stock_prices(company_id, timestamp DESC);
CREATE INDEX idx_stock_prices_company_id ON stock_prices(company_id);

-- =====================================================
-- MARKET_EVENTS TABLE
-- Purpose: Store company events (earnings, dividends, news, etc)
-- =====================================================
CREATE TABLE market_events (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id        UUID          NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    event_type        VARCHAR(50)   NOT NULL,              -- earnings, news, dividend, merger, ipo, split
    title             VARCHAR(255)  NOT NULL,
    description       TEXT,
    impact_percentage DOUBLE PRECISION NOT NULL DEFAULT 0, -- Positive or negative market impact
    event_date        TIMESTAMPTZ   NOT NULL,
    created_at        TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_market_events_company_id ON market_events(company_id);

-- =====================================================
-- VIRTUAL_WALLETS TABLE
-- Purpose: Store user trading wallets (combines balance and fiat_balance)
-- Unified wallet for both trading balance and deposited funds
-- =====================================================
CREATE TABLE virtual_wallets (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id           UUID          NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    balance           DOUBLE PRECISION NOT NULL DEFAULT 0,          -- Available balance for trading
    total_invested    DOUBLE PRECISION NOT NULL DEFAULT 0,          -- Total amount invested in stocks
    total_profit_loss DOUBLE PRECISION NOT NULL DEFAULT 0,          -- Total profit/loss from trading
    fiat_balance      DOUBLE PRECISION DEFAULT 0,                   -- Original deposited fiat (for reference)
    locked            BOOLEAN       DEFAULT FALSE,
    version           INTEGER       DEFAULT 1,
    created_at        TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_virtual_wallets_user_id ON virtual_wallets(user_id);

-- =====================================================
-- USER_PORTFOLIOS TABLE
-- Purpose: Store user's stock holdings
-- =====================================================
CREATE TABLE user_portfolios (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id        UUID          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_id     UUID          NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    quantity       INTEGER       NOT NULL DEFAULT 0,
    average_price  DOUBLE PRECISION NOT NULL DEFAULT 0,
    total_invested DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    
    UNIQUE(user_id, company_id)  -- One portfolio entry per user per company
);

CREATE INDEX idx_user_portfolios_user_company ON user_portfolios(user_id, company_id);
CREATE INDEX idx_user_portfolios_user_id ON user_portfolios(user_id);

-- =====================================================
-- TRANSACTIONS TABLE
-- Purpose: Store wallet top-up and refund transactions
-- =====================================================
CREATE TABLE transactions (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      UUID          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type         VARCHAR(50)   NOT NULL,              -- topup, refund
    amount       DOUBLE PRECISION NOT NULL,
    status       VARCHAR(50)   NOT NULL DEFAULT 'pending',  -- pending, success, failed
    reference_id VARCHAR(100)  NOT NULL,
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_user_id ON transactions(user_id);

-- =====================================================
-- STOCK_TRANSACTIONS TABLE
-- Purpose: Store buy/sell trading history
-- =====================================================
CREATE TABLE stock_transactions (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_id      UUID          NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    type            VARCHAR(20)   NOT NULL,              -- buy, sell
    quantity        INTEGER       NOT NULL,
    price_per_share DOUBLE PRECISION NOT NULL,
    total_amount    DOUBLE PRECISION NOT NULL,
    status          VARCHAR(50)   NOT NULL DEFAULT 'pending',  -- pending, completed, failed, cancelled
    reference_id    VARCHAR(100),
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_stock_transactions_user_id ON stock_transactions(user_id);
CREATE INDEX idx_stock_transactions_company ON stock_transactions(company_id);
CREATE INDEX idx_stock_transactions_user_company ON stock_transactions(user_id, company_id);

-- =====================================================
-- AUTO-UPDATE TRIGGER FOR updated_at COLUMN
-- =====================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

-- Apply trigger to all tables
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_companies_updated_at 
    BEFORE UPDATE ON companies 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_stock_prices_updated_at 
    BEFORE UPDATE ON stock_prices 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_market_events_updated_at 
    BEFORE UPDATE ON market_events 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_virtual_wallets_updated_at 
    BEFORE UPDATE ON virtual_wallets 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_portfolios_updated_at 
    BEFORE UPDATE ON user_portfolios 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_transactions_updated_at 
    BEFORE UPDATE ON transactions 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_stock_transactions_updated_at 
    BEFORE UPDATE ON stock_transactions 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- ROW LEVEL SECURITY (RLS)
-- =====================================================
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE companies ENABLE ROW LEVEL SECURITY;
ALTER TABLE stock_prices ENABLE ROW LEVEL SECURITY;
ALTER TABLE market_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE virtual_wallets ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_portfolios ENABLE ROW LEVEL SECURITY;
ALTER TABLE transactions ENABLE ROW LEVEL SECURITY;
ALTER TABLE stock_transactions ENABLE ROW LEVEL SECURITY;

-- Allow public read on market data
CREATE POLICY "Allow public read on companies" ON companies FOR SELECT USING (true);
CREATE POLICY "Allow public read on stock_prices" ON stock_prices FOR SELECT USING (true);
CREATE POLICY "Allow public read on market_events" ON market_events FOR SELECT USING (true);

-- Allow all operations (service role bypasses RLS)
CREATE POLICY "Allow all on users" ON users FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "Allow all on virtual_wallets" ON virtual_wallets FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "Allow all on user_portfolios" ON user_portfolios FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "Allow all on transactions" ON transactions FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "Allow all on stock_transactions" ON stock_transactions FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "Allow all on companies_write" ON companies FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "Allow all on stock_prices_write" ON stock_prices FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "Allow all on market_events_write" ON market_events FOR ALL USING (true) WITH CHECK (true);

-- =====================================================
-- COMMENTS FOR DOCUMENTATION
-- =====================================================
COMMENT ON TABLE users IS 'User accounts with authentication';
COMMENT ON TABLE companies IS 'Stock market companies with market information';
COMMENT ON TABLE stock_prices IS 'Historical stock price data';
COMMENT ON TABLE market_events IS 'Company events (earnings, dividends, news, etc)';
COMMENT ON TABLE virtual_wallets IS 'User trading wallets with balance and investment tracking';
COMMENT ON TABLE user_portfolios IS 'User stock holdings';
COMMENT ON TABLE transactions IS 'Wallet top-up and refund transactions';
COMMENT ON TABLE stock_transactions IS 'Buy/sell trading history';

COMMENT ON COLUMN companies.total_shares IS 'Total shares issued by company';
COMMENT ON COLUMN companies.available_shares IS 'Shares available for trading (decreases when users buy)';
COMMENT ON COLUMN virtual_wallets.balance IS 'Available balance for trading (updates with buys/sells)';
COMMENT ON COLUMN virtual_wallets.total_invested IS 'Total amount invested in stocks';
COMMENT ON COLUMN virtual_wallets.total_profit_loss IS 'Total profit/loss from all trades';
COMMENT ON COLUMN virtual_wallets.fiat_balance IS 'Original deposited fiat amount (for reference/reconciliation)';
COMMENT ON COLUMN user_portfolios.average_price IS 'Weighted average buy price per share';
COMMENT ON COLUMN user_portfolios.total_invested IS 'Total amount invested in this position';

-- =====================================================
-- MIGRATION NOTES
-- =====================================================
-- 1. This migration creates a complete schema from scratch
-- 2. The virtual_wallets table combines balance and fiat_balance in ONE table
-- 3. All foreign keys use ON DELETE CASCADE for data integrity
-- 4. Unique constraints prevent duplicate records
-- 5. Indexes are optimized for common query patterns
-- 6. RLS policies allow the backend service role full access
-- 7. AUTO-UPDATE triggers keep updated_at column current
-- 8. Run this migration ONCE for fresh database setup
