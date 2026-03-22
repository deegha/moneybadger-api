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
    limit_amount, 
    month, 
    year
) VALUES (
    $1, $2, $3, $4, $5
)
ON CONFLICT (user_id, category_id, month, year) 
DO UPDATE SET 
    limit_amount = EXCLUDED.limit_amount
RETURNING *;

-- name: GetBudgetByCategory :one
SELECT * FROM budgets
WHERE user_id = $1 
  AND category_id = $2 
  AND month = $3 
  AND year = $4;

-- name: ListBudgetsByMonth :many
SELECT * FROM budgets
WHERE user_id = $1 
  AND month = $2 
  AND year = $3;


-- name: GetUserCategoriesWithBudgets :many
SELECT 
    c.*, 
    b.limit_amount as budget_limit,
    b.month as budget_month,
    b.year as budget_year
FROM categories c
LEFT JOIN budgets b ON 
    c.id = b.category_id AND 
    b.user_id = $1 AND 
    b.month = $2 AND 
    b.year = $3
WHERE c.user_id = $1
ORDER BY c.name ASC;