package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deegha/moneyBadgerApi/internal/env"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Define a custom type for context keys to avoid collisions
type contextKey string

const userIDKey contextKey = "user_id"

// GetUserID retrieves the typed pgtype.UUID from the context.
// Use this in your handlers to get a database-ready ID.
func GetUserID(ctx context.Context) (pgtype.UUID, error) {
	userID, ok := ctx.Value(userIDKey).(pgtype.UUID)
	if !ok {
		return pgtype.UUID{}, fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Get the cookie
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "Unauthorized: No session cookie", http.StatusUnauthorized)
			return
		}

		// 2. Parse the JWT
		tokenString := cookie.Value
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(env.GetString("JWT_SECRET", "supersecretkey")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// 3. Extract UserID (sub) from claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized: Invalid claims", http.StatusUnauthorized)
			return
		}

		userIDRaw, ok := claims["sub"].(string)
		if !ok {
			http.Error(w, "Unauthorized: Missing user ID", http.StatusUnauthorized)
			return
		}

		// 4. CONVERSION: Convert string to pgtype.UUID right here
		var userID pgtype.UUID
		if err := userID.Scan(userIDRaw); err != nil {
			// This handles cases where the JWT might have an invalid UUID string
			http.Error(w, "Unauthorized: Malformed user ID", http.StatusUnauthorized)
			return
		}

		// 5. Put the typed userID into context
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
