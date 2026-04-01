-- Remove month and year from budgets; budget is now a simple per-category limit
ALTER TABLE budgets DROP CONSTRAINT budgets_user_id_category_id_month_year_key;
ALTER TABLE budgets DROP COLUMN month;
ALTER TABLE budgets DROP COLUMN year;
ALTER TABLE budgets ADD CONSTRAINT budgets_user_id_category_id_key UNIQUE (user_id, category_id);
