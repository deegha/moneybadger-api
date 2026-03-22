-- +goose Up
-- +goose StatementBegin

-- 1. Create Custom Types (Enums)
CREATE TYPE transaction_type AS ENUM ('income', 'expense', 'transfer');
CREATE TYPE frequency AS ENUM ('daily', 'weekly', 'monthly', 'yearly');
CREATE TYPE user_tier AS ENUM ('free', 'pro', 'atelier_pro');

-- 2. Users Table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    full_name TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    avatar_url TEXT,
    tier user_tier DEFAULT 'free',
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- 3. Categories Table
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    icon TEXT, 
    color_hex TEXT DEFAULT '#10B981',
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- 4. Transactions Table
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    amount DECIMAL(15, 2) NOT NULL,
    type transaction_type NOT NULL,
    description TEXT,
    merchant_name TEXT,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    is_recurring BOOLEAN DEFAULT false,
    recurring_bill_id UUID, -- Optional: link to the template
    created_at TIMESTAMPTZ DEFAULT now()
);

-- 5. Recurring Bills Table
CREATE TABLE recurring_bills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id),
    name TEXT NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    frequency frequency DEFAULT 'monthly',
    next_due_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- 6. Budgets Table
CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE CASCADE,
    limit_amount DECIMAL(15, 2) NOT NULL,
    month INTEGER NOT NULL CHECK (month BETWEEN 1 AND 12),
    year INTEGER NOT NULL,
    UNIQUE(user_id, category_id, month, year)
);

-- 7. Performance Indexes
CREATE INDEX idx_transactions_user_date ON transactions(user_id, date DESC);
CREATE INDEX idx_recurring_bills_due ON recurring_bills(next_due_date) WHERE is_active = true;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS budgets;
DROP TABLE IF EXISTS recurring_bills;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;

-- Drop custom types
DROP TYPE IF EXISTS user_tier;
DROP TYPE IF EXISTS frequency;
DROP TYPE IF EXISTS transaction_type;

-- +goose StatementEnd