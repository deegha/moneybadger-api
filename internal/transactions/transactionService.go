package transactions

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/sync/errgroup"

	repo "github.com/deegha/moneyBadgerApi/internal/adapters/postgresql/sqlc"
)

type TransactionService interface {
	ListTransactions(
		ctx context.Context,
		args ListTransacitonsRequest,
	) (ListTransacitonsResponse, error)
	CreateTransaction(ctx context.Context, arg CreateTransactionRequest) (repo.Transaction, error)
	GetSummaryMonth(ctx context.Context, UserID pgtype.UUID) (repo.GetMonthlySummaryRow, error)
	GetOverView(ctx context.Context, args OverViewParams) (ChartData, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) TransactionService {
	return &svc{
		repo: repo,
	}
}

func (s *svc) ListTransactions(
	ctx context.Context,
	args ListTransacitonsRequest,
) (ListTransacitonsResponse, error) {
	trnasactions, err := s.repo.GetTransactionsFiltered(
		ctx,
		repo.GetTransactionsFilteredParams(args),
	)
	if err != nil {
		return ListTransacitonsResponse{}, err
	}

	totalCount, err := s.repo.GetTransactionsCount(ctx, repo.GetTransactionsCountParams{
		UserID:    args.UserID,
		StartDate: args.StartDate,
		EndDate:   args.EndDate,
	})
	if err != nil {
		return ListTransacitonsResponse{}, err
	}

	return ListTransacitonsResponse{
		Transactions: trnasactions,
		TotalCount:   int(totalCount),
	}, nil
}

func (s *svc) CreateTransaction(
	ctx context.Context,
	arg CreateTransactionRequest,
) (repo.Transaction, error) {
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

func (s *svc) GetSummaryMonth(
	ctx context.Context,
	UserID pgtype.UUID,
) (repo.GetMonthlySummaryRow, error) {
	return s.repo.GetMonthlySummary(ctx, UserID)
}

func (s *svc) GetOverView(ctx context.Context, args OverViewParams) (ChartData, error) {
	var chartData ChartData

	dataWeekly := make(chan []repo.GetWeeklySpendingOverviewRow)
	dataMonthly := make(chan []repo.GetMonthlySpendingOverviewRow)
	dataDaily := make(chan []repo.GetSpendingOverviewRow)

	g, ctx := errgroup.WithContext(ctx)

	// wg := sync.WaitGroup{}

	// 1. Fetch Weekly Data
	g.Go(func() error {
		weekly, err := s.repo.GetWeeklySpendingOverview(ctx, args.UserID)
		if err != nil {
			return err
		}
		dataWeekly <- weekly
		close(dataWeekly)
		return nil
	})

	// 2. Fetch Monthly Data
	g.Go(func() error {
		monthly, err := s.repo.GetMonthlySpendingOverview(
			ctx,
			repo.GetMonthlySpendingOverviewParams{
				UserID: args.UserID,
				Month:  args.Month,
				Year:   args.Year,
			},
		)
		if err != nil {
			return err
		}
		dataMonthly <- monthly
		close(dataMonthly)
		return nil
	})

	// 3. Fetch Daily Data
	g.Go(func() error {
		daily, err := s.repo.GetSpendingOverview(ctx, repo.GetSpendingOverviewParams{
			UserID: args.UserID,
			Month:  args.Month,
			Year:   args.Year,
		})
		if err != nil {
			return err
		}
		dataDaily <- daily
		close(dataDaily)
		return nil
	})

	chartData = ChartData{
		Weekly:  <-dataWeekly,
		Daily:   <-dataDaily,
		Monthly: <-dataMonthly,
	}

	// Wait for all goroutines to finish
	if err := g.Wait(); err != nil {
		return ChartData{}, err
	}

	return chartData, nil
}
