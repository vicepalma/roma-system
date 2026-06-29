package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/domain"
	"github.com/vicepalma/roma-system/backend/internal/security"
)

type fakeUserRepo struct {
	byID    map[string]*domain.User
	byEmail map[string]*domain.User
}

func (r *fakeUserRepo) Create(context.Context, *domain.User) error { return nil }
func (r *fakeUserRepo) FindByID(_ context.Context, id string) (*domain.User, error) {
	return r.byID[id], nil
}
func (r *fakeUserRepo) FindByEmail(_ context.Context, email string) (*domain.User, error) {
	return r.byEmail[email], nil
}

func TestAuthLoginReturnsPersistedRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("JWT_SECRET", "test-secret")
	hash, err := security.HashPassword("secret123")
	if err != nil {
		t.Fatal(err)
	}
	repo := &fakeUserRepo{byEmail: map[string]*domain.User{
		"coach@example.test":    {ID: "coach-1", Email: "coach@example.test", Name: "Coach", PasswordHash: hash, Role: "coach"},
		"disciple@example.test": {ID: "disciple-1", Email: "disciple@example.test", Name: "Disciple", PasswordHash: hash, Role: "disciple"},
	}}
	r := gin.New()
	NewAuthHandler(repo, nil).Register(r.Group("/"))

	for _, tc := range []struct {
		email string
		role  string
	}{
		{"coach@example.test", "coach"},
		{"disciple@example.test", "disciple"},
	} {
		body := []byte(`{"email":"` + tc.email + `","password":"secret123"}`)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("login %s status=%d body=%s", tc.email, w.Code, w.Body.String())
		}
		var out struct {
			User struct {
				Role string `json:"role"`
			} `json:"user"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
			t.Fatal(err)
		}
		if out.User.Role != tc.role {
			t.Fatalf("role=%q want %q", out.User.Role, tc.role)
		}
	}
}

func TestMeRequiresTokenAndReturnsRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("JWT_SECRET", "test-secret")
	repo := &fakeUserRepo{byID: map[string]*domain.User{
		"coach-1": {ID: "coach-1", Email: "coach@example.test", Name: "Coach", Role: "coach"},
	}}
	r := gin.New()
	NewAuthHandler(repo, nil).Register(r.Group("/"))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/me", nil))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated /me status=%d", w.Code)
	}

	tokens, err := security.GenerateTokens("coach-1")
	if err != nil {
		t.Fatal(err)
	}
	w = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("authenticated /me status=%d body=%s", w.Code, w.Body.String())
	}
	var out struct {
		Role string `json:"role"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out.Role != "coach" {
		t.Fatalf("role=%q want coach", out.Role)
	}
}
