package categories

import (
	"context"
	"fmt"

	repo "github.com/deegha/moneyBadgerApi/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	CreateCategories(ctx context.Context, args CreateCategoryRequest) (repo.Category, error)
	ListCategories(ctx context.Context, userID pgtype.UUID) ([]repo.GetUserCategoriesWithBudgetsRow, error)
}

type svc struct {
	repo repo.Queries
	db   *pgx.Conn
}

func NewService(repo repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) CreateCategories(ctx context.Context, args CreateCategoryRequest) (repo.Category, error) {

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

	_, err = qtx.CreateOrUpdateBudget(ctx, repo.CreateOrUpdateBudgetParams{
		UserID:      args.UserID,
		CategoryID:  pgtype.UUID{Bytes: category.ID, Valid: true},
		LimitAmount: args.LimitAmount,
		Month:       args.Month,
		Year:        args.Year,
	})

	if err != nil {
		return repo.Category{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return repo.Category{}, err
	}

	return category, nil
}

func (s *svc) ListCategories(ctx context.Context, userID pgtype.UUID) ([]repo.GetUserCategoriesWithBudgetsRow, error) {

	return s.repo.GetUserCategoriesWithBudgets(ctx, repo.GetUserCategoriesWithBudgetsParams{
		UserID: userID,
	})
}
