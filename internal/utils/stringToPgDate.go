package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func StringToPgDate(dateStr string) (pgtype.Date, error) {
	// Parse the string into a standard Go time object
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return pgtype.Date{}, err
	}

	// Assign to pgtype.Date
	return pgtype.Date{
		Time:  t,
		Valid: true, // Crucial: tells pgx this isn't a NULL value
	}, nil
}
