package transactions

import (
	"log"
	"net/http"

	repo "github.com/deegha/moneyBadgerApi/internal/adapters/postgresql/sqlc"
	"github.com/deegha/moneyBadgerApi/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(s Service) *handler {
	return &handler{service: s}
}

func (h *handler) ListTransactions(w http.ResponseWriter, r *http.Request) {

	transactions, err := h.service.ListTransactions(r.Context())

	if err != nil {
		log.Printf("error listing transactions: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(transactions) == 0 {
		json.Writer(w, http.StatusOK, []repo.GetTransactionsFilteredRow{}, "No transactions found")
		return
	}

	json.Writer(w, http.StatusOK, transactions, "Successfully retrieved transactions")
}

func (h *handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest

	if err := json.Reader(r, &req); err != nil {
		log.Printf("error parsing request body: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	transaction, err := h.service.CreateTransaction(r.Context(), req)
	if err != nil {
		log.Printf("error creating transaction: %v", err)
		json.Writer(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	json.Writer(w, http.StatusCreated, transaction, "Successfully created transaction")
}
