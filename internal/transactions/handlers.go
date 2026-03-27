package transactions

import (
	"log"
	"net/http"
	"strconv"

	"github.com/deegha/moneyBadgerApi/internal/json"
	auth "github.com/deegha/moneyBadgerApi/internal/middleware"
	"github.com/deegha/moneyBadgerApi/internal/utils"
)

type handler struct {
	service TransactionService
}

func NewHandler(s TransactionService) *handler {
	return &handler{service: s}
}

func (h *handler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())
	if err != nil {
		log.Printf("Error getting the user id %v", err)
		json.Writer(w, http.StatusInternalServerError, nil, err.Error())
	}

	query := r.URL.Query()

	Limit, _ := strconv.Atoi(query.Get("limit"))
	Offset, _ := strconv.Atoi(query.Get("offset"))
	StartDate, _ := utils.StringToPgDate(query.Get("start_date"))
	EndDate, _ := utils.StringToPgDate(query.Get("end_date"))
	CategoryID, _ := utils.ParseUUID(query.Get("category_id"))

	transactions, err := h.service.ListTransactions(r.Context(), ListTransacitonsRequest{
		UserID:     userID,
		Limit:      int32(Limit),
		Offset:     int32(Offset),
		StartDate:  StartDate,
		EndDate:    EndDate,
		CategoryID: CategoryID,
	})
	if err != nil {
		log.Printf("error listing transactions: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Writer(w, http.StatusOK, transactions, "Successfully retrieved transactions")
}

func (h *handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())
	if err != nil {
		log.Printf("Error getting the user id %v", err)
		json.Writer(w, http.StatusInternalServerError, nil, err.Error())
	}

	var req CreateTransactionRequest

	req.UserID = userID

	if err := json.Reader(r, &req); err != nil {
		log.Printf("error parsing request body: %v", err)
		json.Writer(w, http.StatusBadRequest, nil, err.Error())
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

func (h *handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())
	if err != nil {
		json.Writer(w, http.StatusUnauthorized, nil, "Unauthorized")
		return
	}

	summary, err := h.service.GetSummaryMonth(r.Context(), userID)
	if err != nil {
		log.Printf("Error while fetching transaction summary for user %v: %v", userID, err)
		json.Writer(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	json.Writer(w, http.StatusOK, summary, "Successfully fetched the summary")
}

func (h *handler) GetOverView(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())
	if err != nil {
		json.Writer(w, http.StatusUnauthorized, nil, "Unauthorized")
		return
	}

	query := r.URL.Query()

	month, _ := strconv.Atoi(query.Get("month"))
	year, _ := strconv.Atoi(query.Get("year"))

	overview, err := h.service.GetOverView(r.Context(), OverViewParams{
		UserID: userID,
		Month:  int32(month),
		Year:   int32(year),
	})
	if err != nil {
		log.Printf("Error while fetching transaction overview for user %v: %v", userID, err)
		json.Writer(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	json.Writer(w, http.StatusOK, overview, "Successfully fetched the overview")
}
