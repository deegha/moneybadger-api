package transactions

import (
	"context"
	"fmt"

	repo "github.com/deegha/moneyBadgerApi/internal/adapters/postgresql/sqlc"
)

type Service interface {
	ListTransactions(ctx context.Context) ([]repo.GetTransactionsFilteredRow, error)
	CreateTransaction(ctx context.Context, arg CreateTransactionRequest) (repo.Transaction, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) ListTransactions(ctx context.Context) ([]repo.GetTransactionsFilteredRow, error) {

	return s.repo.GetTransactionsFiltered(ctx, repo.GetTransactionsFilteredParams{})
}

func (s *svc) CreateTransaction(ctx context.Context, arg CreateTransactionRequest) (repo.Transaction, error) {

	if arg.Type != "income" && arg.Type != "expense" {
		return repo.Transaction{}, fmt.Errorf("type must be either 'income' or 'expense'")
	}
	return s.repo.CreateTransaction(ctx, repo.CreateTransactionParams{
		UserID:       arg.UserID,
		Amount:       arg.Amount,
		Description:  arg.Description,
		Date:         arg.Date,
		CategoryID:   arg.CategoryID,
		Type:         repo.TransactionType(arg.Type),
		MerchantName: arg.MerchantName,
		IsRecurring:  arg.IsRecurring,
	})
}
