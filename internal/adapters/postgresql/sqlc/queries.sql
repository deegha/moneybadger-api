-- name: CreateUser :one
INSERT INTO users (
    full_name, email, password_hash, tier
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: CreateCategory :one
INSERT INTO categories (
    user_id, name, icon, color_hex, is_default
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUserCategories :many
SELECT * FROM categories
WHERE user_id = $1
ORDER BY name ASC;

-- name: CreateTransaction :one
INSERT INTO transactions (
    user_id, category_id, amount, type, description, merchant_name, date, is_recurring
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetRecentTransactions :many
SELECT 
    t.*, 
    c.name as category_name, 
    c.icon as category_icon, 
    c.color_hex as category_color
FROM transactions t
LEFT JOIN categories c ON t.category_id = c.id
WHERE t.user_id = $1
ORDER BY t.date DESC, t.created_at DESC
LIMIT $2;

-- name: GetMonthlySummary :one
-- Calculates the "Total Income" and "Total Expenses" for the current month
SELECT 
    SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END)::DECIMAL(15,2) as total_income,
    SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END)::DECIMAL(15,2) as total_expense
FROM transactions
WHERE user_id = $1 
  AND date >= date_trunc('month', CURRENT_DATE)
  AND date <= (date_trunc('month', CURRENT_DATE) + interval '1 month - 1 day');

-- name: GetUpcomingBills :many
SELECT * FROM recurring_bills
WHERE user_id = $1 
  AND is_active = true 
  AND next_due_date <= (CURRENT_DATE + interval '30 days')
ORDER BY next_due_date ASC;

-- name: UpdateRecurringBillDate :exec
UPDATE recurring_bills
SET next_due_date = $2, updated_at = now()
WHERE id = $1;


-- name: GetTransactionsFiltered :many
-- Fetch transactions with pagination and date range filtering
-- Supports category-specific filtering if category_id is provided
SELECT 
    t.id,
    t.amount,
    t.type,
    t.description,
    t.merchant_name,
    t.date,
    t.is_recurring,
    c.name as category_name,
    c.icon as category_icon,
    c.color_hex as category_color
FROM transactions t
LEFT JOIN categories c ON t.category_id = c.id
WHERE t.user_id = $1
  AND (t.date >= sqlc.narg('start_date') OR sqlc.narg('start_date') IS NULL)
  AND (t.date <= sqlc.narg('end_date') OR sqlc.narg('end_date') IS NULL)
  AND (t.category_id = sqlc.narg('category_id') OR sqlc.narg('category_id') IS NULL)
ORDER BY t.date DESC, t.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetTransactionsCount :one
-- Used to calculate total pages in the UI
SELECT COUNT(*) FROM transactions
WHERE user_id = $1
  AND (date >= sqlc.narg('start_date') OR sqlc.narg('start_date') IS NULL)
  AND (date <= sqlc.narg('end_date') OR sqlc.narg('end_date') IS NULL);

-- name: CreateOrUpdateBudget :one
INSERT INTO budgets (
    user_id,
    category_id,
    limit_amount
) VALUES (
    $1, $2, $3
)
ON CONFLICT (user_id, category_id)
DO UPDATE SET
    limit_amount = EXCLUDED.limit_amount
RETURNING *;

-- name: GetBudgetByCategory :one
SELECT * FROM budgets
WHERE user_id = $1
  AND category_id = $2;

-- name: ListBudgets :many
SELECT * FROM budgets
WHERE user_id = $1;

-- name: GetUserCategoriesWithBudgets :many
SELECT
    c.id,
    c.user_id,
    c.name,
    c.icon,
    c.color_hex,
    c.created_at,
    COALESCE(b.limit_amount, 0.00)::DECIMAL(15,2) as budget_limit,
    COALESCE(s.total_spent, 0.00)::DECIMAL(15,2) as total_spent,
    -- Calculate percentage for the progress bar, capped at 100 or allowed to exceed
    CASE
        WHEN COALESCE(b.limit_amount, 0) = 0 THEN 0
        ELSE ROUND((COALESCE(s.total_spent, 0) / b.limit_amount) * 100)
    END as spent_percentage
FROM categories c
LEFT JOIN budgets b ON c.id = b.category_id
LEFT JOIN (
    /* Subquery to sum transactions for the current month */
    SELECT
        category_id,
        SUM(amount) as total_spent
    FROM transactions
    WHERE
        EXTRACT(MONTH FROM date) = EXTRACT(MONTH FROM CURRENT_DATE) AND
        EXTRACT(YEAR FROM date) = EXTRACT(YEAR FROM CURRENT_DATE) AND
        type = 'expense'
    GROUP BY category_id
) s ON c.id = s.category_id
WHERE c.user_id = $1
ORDER BY c.name ASC;


-- name: GetSpendingOverview :many
WITH RECURSIVE days AS (
    -- Start at the first day of the given month/year
    SELECT 
        make_date(sqlc.arg('year')::int, sqlc.arg('month')::int, 1)::date AS day
    UNION ALL
    -- Increment by 1 day until we hit the last day of the month
    SELECT 
        (day + interval '1 day')::date
    FROM days
    WHERE day < (make_date(sqlc.arg('year')::int, sqlc.arg('month')::int, 1) + interval '1 month - 1 day')::date
)
SELECT 
    d.day,
    -- Aggregate expenses, defaulting to 0 for days with no activity
    COALESCE(SUM(t.amount), 0)::numeric AS total_amount
FROM days d
LEFT JOIN transactions t ON d.day = t.date 
    AND t.user_id = $1 
    AND t.type = 'expense'
GROUP BY d.day
ORDER BY d.day ASC;


-- name: GetMonthlySpendingOverview :many
WITH RECURSIVE days AS (
    -- Start at the first day of the requested month
    SELECT 
        make_date(sqlc.arg('year')::int, sqlc.arg('month')::int, 1)::date AS day
    UNION ALL
    -- Increment by 1 day until the end of the month
    SELECT 
        (day + interval '1 day')::date
    FROM days
    WHERE day < (make_date(sqlc.arg('year')::int, sqlc.arg('month')::int, 1) + interval '1 month - 1 day')::date
)
SELECT 
    d.day,
    COALESCE(SUM(t.amount), 0)::numeric AS total_amount
FROM days d
LEFT JOIN transactions t ON d.day = t.date 
    AND t.user_id = $1 
    AND t.type = 'expense'
GROUP BY d.day
ORDER BY d.day ASC;


-- name: GetWeeklySpendingOverview :many
SELECT 
    date_trunc('week', date)::date AS week_start,
    SUM(amount)::numeric AS total_amount
FROM transactions
WHERE user_id = $1 
    AND type = 'expense'
    AND date >= CURRENT_DATE - INTERVAL '12 weeks' -- Shows last 3 months of habits
GROUP BY week_start
ORDER BY week_start ASC;