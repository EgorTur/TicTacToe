package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"tic-tac-toe/internal/application/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware_ValidToken(t *testing.T) {
	provider := auth.NewJwtProvider("access_secret", "refresh_secret", 15*time.Minute, 24*time.Hour)
	userID := uuid.New()
	token, _ := provider.GenerateAccessToken(userID)

	authenticator := NewJwtAuthenticator(provider)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid, ok := r.Context().Value("userID").(uuid.UUID)
		assert.True(t, ok)
		assert.Equal(t, userID, uid)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	authenticator.Middleware(nextHandler).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMiddleware_MissingHeader(t *testing.T) {
	provider := auth.NewJwtProvider("access_secret", "refresh_secret", time.Hour, time.Hour)
	authenticator := NewJwtAuthenticator(provider)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	authenticator.Middleware(nextHandler).ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMiddleware_MalformedHeader(t *testing.T) {
	provider := auth.NewJwtProvider("access_secret", "refresh_secret", time.Hour, time.Hour)
	authenticator := NewJwtAuthenticator(provider)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	})

	testCases := []string{
		"Basic something",   // не Bearer
		"Bearertoken",       // нет пробела
		"",                  // пусто
	}
	for _, header := range testCases {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		if header != "" {
			req.Header.Set("Authorization", header)
		}
		w := httptest.NewRecorder()
		authenticator.Middleware(nextHandler).ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code, "for header: %s", header)
	}
}

func TestMiddleware_ExpiredToken(t *testing.T) {
	// токен просрочен сразу
	provider := auth.NewJwtProvider("access_secret", "refresh_secret", -1*time.Hour, 24*time.Hour)
	userID := uuid.New()
	token, _ := provider.GenerateAccessToken(userID)

	authenticator := NewJwtAuthenticator(provider)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	authenticator.Middleware(nextHandler).ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMiddleware_InvalidSignature(t *testing.T) {
	provider := auth.NewJwtProvider("access_secret", "refresh_secret", time.Hour, time.Hour)
	userID := uuid.New()
	token, _ := provider.GenerateAccessToken(userID)

	// неправильный секрет
	anotherProvider := auth.NewJwtProvider("wrong_secret", "refresh_secret", time.Hour, time.Hour)
	authenticator := NewJwtAuthenticator(anotherProvider)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	authenticator.Middleware(nextHandler).ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}