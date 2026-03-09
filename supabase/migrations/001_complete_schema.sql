-- ============================================================
-- Share Market Simulator — Complete Database Schema (V3)
-- Fresh start with: company events, extended company fields,
-- candlestick data, profile image URL, and more.
-- ============================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ────────────────────── USERS ──────────────────────
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    phone TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    kyc_status TEXT NOT NULL DEFAULT 'pending' CHECK (kyc_status IN ('pending','verified','rejected','under_review')),
    role TEXT NOT NULL DEFAULT 'user' CHECK (role IN ('user','admin')),
    is_admin BOOLEAN DEFAULT FALSE,
    profile_image_url TEXT DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- ────────────────────── COMPANIES ──────────────────────
CREATE TABLE IF NOT EXISTS companies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    symbol TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    sector TEXT NOT NULL DEFAULT 'General',
    description TEXT DEFAULT '',
    total_supply BIGINT NOT NULL DEFAULT 0,
    shares_outstanding NUMERIC(20,2) DEFAULT 0,
    current_price NUMERIC(20,6) DEFAULT 0,
    market_cap NUMERIC(20,6) DEFAULT 0,
    eps NUMERIC(20,6) DEFAULT 0,
    pe_ratio NUMERIC(20,6) DEFAULT 0,
    book_value NUMERIC(20,6) DEFAULT 0,
    pbv NUMERIC(20,6) DEFAULT 0,
    week_52_high NUMERIC(20,6) DEFAULT 0,
    week_52_low NUMERIC(20,6) DEFAULT 0,
    avg_120_day NUMERIC(20,6) DEFAULT 0,
    yield_1_year NUMERIC(20,6) DEFAULT 0,
    listed_date DATE DEFAULT CURRENT_DATE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- ────────────────────── COMPANY EVENTS ──────────────────────
CREATE TABLE IF NOT EXISTS company_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL CHECK (event_type IN (
        'agm','dividend','bonus_share','right_share','quarterly_report',
        'board_meeting','financial_results','stock_split','merger_acquisition','ipo_announcement'
    )),
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    event_date TIMESTAMPTZ NOT NULL,
    fiscal_year TEXT DEFAULT '',
    status TEXT DEFAULT 'upcoming' CHECK (status IN ('upcoming','ongoing','completed','cancelled')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- ────────────────────── IPOs ──────────────────────
CREATE TABLE IF NOT EXISTS ipos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    price_per_share NUMERIC(20,6) NOT NULL,
    total_shares BIGINT NOT NULL,
    allocated_shares BIGINT DEFAULT 0,
    max_per_applicant BIGINT NOT NULL,
    open_at TIMESTAMPTZ NOT NULL,
    close_at TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','open','closed','allocated')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- ────────────────────── IPO APPLICATIONS ──────────────────────
CREATE TABLE IF NOT EXISTS ipo_applications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ipo_id UUID NOT NULL REFERENCES ipos(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    shares_requested BIGINT NOT NULL,
    shares_allocated BIGINT DEFAULT 0,
    amount_paid NUMERIC(20,6) NOT NULL,
    amount_refunded NUMERIC(20,6) DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','allocated','not_allocated','refunded')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(ipo_id, user_id)
);

-- ────────────────────── WALLETS ──────────────────────
CREATE TABLE IF NOT EXISTS main_wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    balance NUMERIC(20,6) DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS trading_wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    balance NUMERIC(20,6) DEFAULT 0,
    locked_balance NUMERIC(20,6) DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS wallet_transfers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount NUMERIC(20,6) NOT NULL,
    direction TEXT NOT NULL CHECK (direction IN ('main_to_trading','trading_to_main')),
    status TEXT NOT NULL DEFAULT 'completed',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ────────────────────── PORTFOLIOS ──────────────────────
CREATE TABLE IF NOT EXISTS portfolios (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    quantity BIGINT NOT NULL DEFAULT 0,
    avg_buy_price NUMERIC(20,6) DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, company_id)
);

-- ────────────────────── ORDERS ──────────────────────
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    side TEXT NOT NULL CHECK (side IN ('buy','sell')),
    order_type TEXT NOT NULL DEFAULT 'limit' CHECK (order_type IN ('limit','market')),
    price NUMERIC(20,6) NOT NULL,
    quantity BIGINT NOT NULL,
    filled_qty BIGINT DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open','partially_filled','filled','cancelled')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- ────────────────────── TRADES ──────────────────────
CREATE TABLE IF NOT EXISTS trades (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    buy_order_id UUID NOT NULL REFERENCES orders(id),
    sell_order_id UUID NOT NULL REFERENCES orders(id),
    buyer_id UUID NOT NULL REFERENCES users(id),
    seller_id UUID NOT NULL REFERENCES users(id),
    price NUMERIC(20,6) NOT NULL,
    quantity BIGINT NOT NULL,
    total_amount NUMERIC(20,6) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ────────────────────── STOCK PRICES (OHLCV Candles) ──────────────────────
CREATE TABLE IF NOT EXISTS stock_prices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    open_price NUMERIC(20,6) NOT NULL,
    high_price NUMERIC(20,6) NOT NULL,
    low_price NUMERIC(20,6) NOT NULL,
    close_price NUMERIC(20,6) NOT NULL,
    volume BIGINT DEFAULT 0,
    turnover NUMERIC(20,6) DEFAULT 0,
    change_percent NUMERIC(10,4) DEFAULT 0,
    timestamp TIMESTAMPTZ NOT NULL,
    timeframe TEXT NOT NULL DEFAULT '1D' CHECK (timeframe IN ('1m','5m','15m','1H','1D','1W','1M')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- ────────────────────── PRICE TRIGGERS ──────────────────────
CREATE TABLE IF NOT EXISTS price_triggers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    trigger_price NUMERIC(20,6) NOT NULL,
    shares_qty BIGINT NOT NULL,
    direction TEXT NOT NULL CHECK (direction IN ('above','below')),
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','triggered','cancelled')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- ────────────────────── MARKET INDEX HISTORY ──────────────────────
CREATE TABLE IF NOT EXISTS market_index_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    index_value NUMERIC(20,6) NOT NULL,
    change NUMERIC(20,6) DEFAULT 0,
    change_percent NUMERIC(10,4) DEFAULT 0,
    total_turnover NUMERIC(20,6) DEFAULT 0,
    total_volume BIGINT DEFAULT 0,
    total_market_cap NUMERIC(20,6) DEFAULT 0,
    advances INT DEFAULT 0,
    declines INT DEFAULT 0,
    unchanged INT DEFAULT 0,
    total_companies INT DEFAULT 0,
    timestamp TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ────────────────────── INDEXES ──────────────────────
CREATE INDEX IF NOT EXISTS idx_companies_sector ON companies(sector);
CREATE INDEX IF NOT EXISTS idx_companies_listed_date ON companies(listed_date);
CREATE INDEX IF NOT EXISTS idx_companies_active ON companies(is_active);
CREATE INDEX IF NOT EXISTS idx_stock_prices_company_time ON stock_prices(company_id, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_stock_prices_timeframe ON stock_prices(timeframe);
CREATE INDEX IF NOT EXISTS idx_orders_company_side ON orders(company_id, side, status);
CREATE INDEX IF NOT EXISTS idx_orders_user ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_trades_company ON trades(company_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_trades_buyer ON trades(buyer_id);
CREATE INDEX IF NOT EXISTS idx_trades_seller ON trades(seller_id);
CREATE INDEX IF NOT EXISTS idx_portfolios_user ON portfolios(user_id);
CREATE INDEX IF NOT EXISTS idx_company_events_company ON company_events(company_id);
CREATE INDEX IF NOT EXISTS idx_company_events_date ON company_events(event_date);
CREATE INDEX IF NOT EXISTS idx_company_events_type ON company_events(event_type);
CREATE INDEX IF NOT EXISTS idx_price_triggers_company ON price_triggers(company_id, status);

-- ────────────────────── AUTO-UPDATE TRIGGERS ──────────────────────
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
DECLARE
    tbl TEXT;
BEGIN
    FOR tbl IN
        SELECT unnest(ARRAY[
            'users','companies','company_events','ipos','ipo_applications',
            'main_wallets','trading_wallets','portfolios','orders','stock_prices','price_triggers'
        ])
    LOOP
        EXECUTE format('
            DROP TRIGGER IF EXISTS trigger_update_%s ON %I;
            CREATE TRIGGER trigger_update_%s
                BEFORE UPDATE ON %I
                FOR EACH ROW EXECUTE FUNCTION update_updated_at();
        ', tbl, tbl, tbl, tbl);
    END LOOP;
END $$;

-- ────────────────────── ROW LEVEL SECURITY ──────────────────────
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE companies ENABLE ROW LEVEL SECURITY;
ALTER TABLE company_events ENABLE ROW LEVEL SECURITY;
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

-- Public read for market data
DO $$
DECLARE
    tbl TEXT;
BEGIN
    FOR tbl IN
        SELECT unnest(ARRAY[
            'companies','company_events','stock_prices','market_index_history','ipos','trades'
        ])
    LOOP
        EXECUTE format('
            DROP POLICY IF EXISTS public_read_%s ON %I;
            CREATE POLICY public_read_%s ON %I FOR SELECT USING (true);
        ', tbl, tbl, tbl, tbl);
    END LOOP;
END $$;

-- Service role full access
DO $$
DECLARE
    tbl TEXT;
BEGIN
    FOR tbl IN
        SELECT unnest(ARRAY[
            'users','companies','company_events','ipos','ipo_applications',
            'main_wallets','trading_wallets','wallet_transfers','portfolios',
            'orders','trades','stock_prices','price_triggers','market_index_history'
        ])
    LOOP
        EXECUTE format('
            DROP POLICY IF EXISTS service_role_all_%s ON %I;
            CREATE POLICY service_role_all_%s ON %I FOR ALL USING (true) WITH CHECK (true);
        ', tbl, tbl, tbl, tbl);
    END LOOP;
END $$;
