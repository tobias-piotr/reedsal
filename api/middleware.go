package api

import (
	"context"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type ClaimKey string

const (
	SubKey = ClaimKey("sub")
)

// TODO: Two separate places for this?
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			RespondWithError(w, NewAPIError(http.StatusUnauthorized, "Missing credentials", nil))
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			RespondWithError(w, NewAPIError(http.StatusUnauthorized, "Invalid credentials", nil))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx := context.WithValue(r.Context(), SubKey, claims["sub"])
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			RespondWithError(w, NewAPIError(http.StatusUnauthorized, "Invalid credentials (missing sub claim)", nil))
		}
	})
}
