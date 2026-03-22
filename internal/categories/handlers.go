package categories

import (
	"log"
	"net/http"

	repo "github.com/deegha/moneyBadgerApi/internal/adapters/postgresql/sqlc"
	"github.com/deegha/moneyBadgerApi/internal/json"
	auth "github.com/deegha/moneyBadgerApi/internal/middleware"
)

type handler struct {
	service Service
}

func NewHandler(s Service) *handler {
	return &handler{service: s}
}

func (h *handler) CreateCategories(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())

	if err != nil {
		json.Writer(w, http.StatusUnauthorized, nil, "Unauthorized")
		return
	}

	var request CreateCategoryRequest

	request.UserID = userID

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

	userID, err := auth.GetUserID(r.Context())
	if err != nil {
		json.Writer(w, http.StatusUnauthorized, nil, "Unauthorized")
		return
	}

	categories, err := h.service.ListCategories(r.Context(), userID)
	if err != nil {
		json.Writer(w, http.StatusInternalServerError, nil, "Failed to fetch the categories")
		return
	}

	if len(categories) == 0 {
		json.Writer(w, http.StatusOK, []repo.GetUserCategoriesWithBudgetsRow{}, "Successfully fetched categories")
		return
	}

	json.Writer(w, http.StatusOK, categories, "Successfully fetched categories")
}
