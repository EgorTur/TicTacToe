package web

import (
	"context"
	"net/http"
	"strings"
	"tic-tac-toe/internal/application/auth"
)

type JwtAuthenticator struct {
	jwtProvider *auth.JwtProvider
}

type Authenticator interface {
	Middleware(next http.Handler) http.Handler
}

func NewJwtAuthenticator(jwtProwider *auth.JwtProvider) *JwtAuthenticator {
	return &JwtAuthenticator{
		jwtProvider: jwtProwider,
	}
}

func (ua *JwtAuthenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeError(w, http.StatusUnauthorized, "authorization header required")
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			writeError(w, http.StatusUnauthorized, "invalid authorization header format")
			return
		}
		tokenString := parts[1]
		userID, err := ua.jwtProvider.ValidateAccessToken(tokenString)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
