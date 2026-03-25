package transactions

import (
	repo "github.com/deegha/moneyBadgerApi/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type TransactionType string

type CreateTransactionRequest struct {
	UserID       pgtype.UUID     `json:"user_id"`
	CategoryID   pgtype.UUID     `json:"category_id"`
	Amount       pgtype.Numeric  `json:"amount"`
	Type         TransactionType `json:"type"`
	Description  pgtype.Text     `json:"description"`
	MerchantName pgtype.Text     `json:"merchant_name"`
	Date         pgtype.Date     `json:"date"`
	IsRecurring  pgtype.Bool     `json:"is_recurring"`
}

type ChartData struct {
	Monthly []repo.GetMonthlySpendingOverviewRow
	Weekly  []repo.GetWeeklySpendingOverviewRow
	Daily   []repo.GetSpendingOverviewRow
}

type OverViewParams struct {
	UserID pgtype.UUID `json:"user_id"`
	Year   int32       `json:"year"`
	Month  int32       `json:"month"`
}

type OverViewRequest struct {
	Year  int32 `json:"year"`
	Month int32 `json:"month"`
}
