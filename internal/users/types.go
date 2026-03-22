package users

import repo "github.com/deegha/moneyBadgerApi/internal/adapters/postgresql/sqlc"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  repo.User
	Token string
}

type RegisterResponse struct {
	User  repo.User
	Token string
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
