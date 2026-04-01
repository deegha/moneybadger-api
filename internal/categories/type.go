package categories

import "github.com/jackc/pgx/v5/pgtype"

type CreateCategoryRequest struct {
	UserID      pgtype.UUID    `json:"user_id"`
	Name        string         `json:"name"`
	Icon        pgtype.Text    `json:"icon"`
	ColorHex    pgtype.Text    `json:"color_hex"`
	IsDefault   pgtype.Bool    `json:"is_default"`
	Month       int32          `json:"month"`
	Year        int32          `json:"year"`
	LimitAmount pgtype.Numeric `json:"limit_amount"`
}

type GetCategories struct {
	UserID pgtype.UUID
	Month  int32
	Year   int32
}
