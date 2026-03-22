package users

import (
	"context"
	"fmt"
	"time"

	repo "github.com/deegha/moneyBadgerApi/internal/adapters/postgresql/sqlc"
	"github.com/deegha/moneyBadgerApi/internal/env"
	"github.com/deegha/moneyBadgerApi/internal/hash"
	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	login(ctx context.Context, email, password string) (LoginResponse, error)
	register(ctx context.Context, name, email, password string) (RegisterResponse, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) login(ctx context.Context, email, password string) (LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("user not found")
	}

	isValid, err := hash.VerifyPassword(password, user.PasswordHash)
	if !isValid || err != nil {
		return LoginResponse{}, fmt.Errorf("invalid credentials")
	}

	// 1. Create the JWT Claims
	claims := jwt.MapClaims{
		"sub":  user.ID.String(),
		"exp":  time.Now().Add(time.Hour * 72).Unix(), // 3 days expiry
		"iat":  time.Now().Unix(),
		"tier": user.Tier,
	}

	// 2. Sign the token with your Secret Key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(env.GetString("JWT_SECRET", "supersecretkey")))
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{User: user, Token: t}, nil
}

func (s *svc) register(ctx context.Context, name, email, password string) (RegisterResponse, error) {

	hashedPassword, err := hash.HashPassword(password)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("failed to process security: %w", err)
	}

	user, err := s.repo.CreateUser(ctx, repo.CreateUserParams{
		FullName:     name,
		Email:        email,
		PasswordHash: hashedPassword,
	})

	if err != nil {
		return RegisterResponse{}, fmt.Errorf("failed to create account: %w", err)
	}

	// 3. Generate JWT (Automatic Login)
	claims := jwt.MapClaims{
		"sub":  user.ID.String(),
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
		"iat":  time.Now().Unix(),
		"tier": user.Tier,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(env.GetString("JWT_SECRET", "supersecretkey")))
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("account created, but login failed: %w", err)
	}

	return RegisterResponse{User: user, Token: t}, nil
}
