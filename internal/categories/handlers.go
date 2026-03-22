package categories

import (
	"log"
	"net/http"

	"github.com/deegha/moneyBadgerApi/internal/json"
	auth "github.com/deegha/moneyBadgerApi/internal/middleware"
	"github.com/jackc/pgx/v5/pgtype"
)

type handler struct {
	service Service
}

func NewHandler(s Service) *handler {
	return &handler{service: s}
}

func (h *handler) CreateCategories(w http.ResponseWriter, r *http.Request) {
	var request CreateCategoryRequest

	if err := json.Reader(r, &request); err != nil {
		log.Printf("error parsing request body: %v", err)
		json.Writer(w, http.StatusBadRequest, nil, err.Error())

		return
	}

	category, err := h.service.CreateCategories(r.Context(), request)

	if err != nil {
		log.Printf("Error creating category: %v", err)
		json.Writer(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	json.Writer(w, http.StatusCreated, category, "Successfully created the category")
}

func (h *handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(auth.UserIDKey).(pgtype.UUID)

	if !ok {
		json.Writer(w, http.StatusUnauthorized, nil, "Authentication failed")
	}

	categories, err := h.service.ListCategories(r.Context(), userIDStr)

	if err != nil {
		json.Writer(w, http.StatusInternalServerError, nil, "Failed to fetch the categories")
	}

	json.Writer(w, http.StatusOK, categories, "Successfully fetched categories")
}
