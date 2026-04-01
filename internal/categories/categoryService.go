package categories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	repo "github.com/deegha/moneyBadgerApi/internal/adapters/postgresql/sqlc"
)

type CategoryService interface {
	CreateCategories(ctx context.Context, args CreateCategoryRequest) (repo.Category, error)
	ListCategories(
		ctx context.Context,
		arge GetCategories,
	) ([]repo.GetUserCategoriesWithBudgetsRow, error)
}

type svc struct {
	repo repo.Queries
	db   *pgxpool.Pool
}

func NewService(repo repo.Queries, db *pgxpool.Pool) CategoryService {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) CreateCategories(
	ctx context.Context,
	args CreateCategoryRequest,
) (repo.Category, error) {
	if args.Name == "" {
		return repo.Category{}, fmt.Errorf("Name of the category is required")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Category{}, fmt.Errorf("Something wrong when creating transaction %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	category, err := qtx.CreateCategory(ctx, repo.CreateCategoryParams{
		UserID:   args.UserID,
		Name:     args.Name,
		Icon:     args.Icon,
		ColorHex: args.ColorHex,
	})
	if err != nil {
		return repo.Category{}, err
	}

	now := time.Now()

	_, err = qtx.CreateOrUpdateBudget(ctx, repo.CreateOrUpdateBudgetParams{
		UserID:      args.UserID,
		CategoryID:  pgtype.UUID{Bytes: category.ID, Valid: true},
		LimitAmount: args.LimitAmount,
		Month:       int32(now.Month()),
		Year:        int32(now.Year()),
	})
	if err != nil {
		return repo.Category{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return repo.Category{}, err
	}

	return category, nil
}

func (s *svc) ListCategories(
	ctx context.Context,
	args GetCategories,
) ([]repo.GetUserCategoriesWithBudgetsRow, error) {

	return s.repo.GetUserCategoriesWithBudgets(ctx, repo.GetUserCategoriesWithBudgetsParams{
		UserID: args.UserID,
		Month:  args.Month,
		Year:   args.Year,
	})
}
