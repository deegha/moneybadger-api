package transactions

import (
	"github.com/jackc/pgx/v5/pgtype"

	repo "github.com/deegha/moneyBadgerApi/internal/adapters/postgresql/sqlc"
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

type ListTransacitonsRequest struct {
	UserID     pgtype.UUID `json:"user_id"`
	Limit      int32       `json:"limit"`
	Offset     int32       `json:"offset"`
	StartDate  pgtype.Date `json:"start_date"`
	EndDate    pgtype.Date `json:"end_date"`
	CategoryID pgtype.UUID `json:"category_id"`
}

type ListTransacitonsResponse struct {
	Transactions []repo.GetTransactionsFilteredRow `json:"transactions"`
	TotalCount   int                               `json:"total_count"`
}
