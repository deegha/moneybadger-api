package utils

import "github.com/jackc/pgx/v5/pgtype"

func ParseUUID(s string) (pgtype.UUID, error) {
	var u pgtype.UUID
	// Scan handles the string parsing and sets the Valid flag automatically
	err := u.Scan(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return u, nil
}
