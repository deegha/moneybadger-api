package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/deegha/moneyBadgerApi/internal/env"
)

func main() {
	ctx := context.Background()
	cfg := config{
		addr: ":3000",
		db: dbConfig{
			dsn: env.GetString(
				"DATABASE_URL",
				"host=localhost port=5432 user=atelier_admin password=atelier_password_2026 dbname=moneybadger sslmode=disable",
			),
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// conn, err := pgx.Connect(ctx, cfg.db.dsn)
	conn, err := pgxpool.New(ctx, cfg.db.dsn)
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	logger.Info("Successfully connected to database")

	api := application{
		config: cfg,
		db:     conn,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("Server has fail to start", "error", err)
		os.Exit(1)
	}
}
