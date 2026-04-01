package categories

import "github.com/jackc/pgx/v5/pgtype"

type CreateCategoryRequest struct {
	UserID      pgtype.UUID    `json:"user_id"`
	Name        string         `json:"name"`
	Icon        pgtype.Text    `json:"icon"`
	ColorHex    pgtype.Text    `json:"color_hex"`
	IsDefault   pgtype.Bool    `json:"is_default"`
	LimitAmount pgtype.Numeric `json:"limit_amount"`
}

type GetCategories struct {
	UserID pgtype.UUID
}
