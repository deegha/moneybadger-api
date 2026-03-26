package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	repo "github.com/deegha/moneyBadgerApi/internal/adapters/postgresql/sqlc"
	"github.com/deegha/moneyBadgerApi/internal/categories"
	auth "github.com/deegha/moneyBadgerApi/internal/middleware"
	"github.com/deegha/moneyBadgerApi/internal/transactions"
	"github.com/deegha/moneyBadgerApi/internal/users"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		// Explicitly list your frontend URL, DO NOT use "*"
		AllowedOrigins:   []string{"http://localhost:3001"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true, // CRITICAL for cookies to work
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})

	// public route
	usersService := users.NewService(repo.New(app.db))
	usersHandler := users.NewHandler(usersService)

	r.Post("/v1/login", usersHandler.Login)
	r.Post("/v1/register", usersHandler.Register)

	transactionsService := transactions.NewService(repo.New(app.db))
	transactionsHandler := transactions.NewHandler(transactionsService)
	categoriesService := categories.NewService(*repo.New(app.db), app.db)
	categoriesHandler := categories.NewHandler(categoriesService)
	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)

		r.Get("/v1/transactions", transactionsHandler.ListTransactions)
		r.Get("/v1/transactions/summary", transactionsHandler.GetSummary)
		r.Get("/v1/transactions/overview", transactionsHandler.GetOverView)
		r.Post("/v1/transactions", transactionsHandler.CreateTransaction)

		r.Get("/v1/categories", categoriesHandler.ListCategories)
		r.Post("/v1/categories", categoriesHandler.CreateCategories)
	})

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Printf("starting server on %s", app.config.addr)

	return srv.ListenAndServe()
}

type application struct {
	config config
	db     *pgxpool.Pool
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}
