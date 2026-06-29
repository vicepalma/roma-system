package http

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/repository"
	"github.com/vicepalma/roma-system/backend/internal/security"
	"github.com/vicepalma/roma-system/backend/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type fakeExerciseService struct{}

func (fakeExerciseService) List(context.Context, repository.ExerciseFilter) ([]repository.Exercise, int64, error) {
	return []repository.Exercise{}, 0, nil
}
func (fakeExerciseService) Create(context.Context, service.CreateExercise) (*repository.Exercise, error) {
	return &repository.Exercise{ID: "exercise-1", Name: "Bench", PrimaryMuscle: "chest"}, nil
}
func (fakeExerciseService) Get(context.Context, string) (*repository.Exercise, error) {
	return nil, nil
}
func (fakeExerciseService) Update(context.Context, string, service.UpdateExercise) (*repository.Exercise, error) {
	return nil, nil
}
func (fakeExerciseService) Delete(context.Context, string) error { return nil }

func TestExercisePermissions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, cleanup := mockGorm(t)
	defer cleanup()

	r := gin.New()
	api := r.Group("/api")
	api.Use(func(c *gin.Context) {
		switch c.GetHeader("X-Test-User") {
		case "coach":
			c.Set(security.CtxUserID, "coach-1")
		case "disciple":
			c.Set(security.CtxUserID, "disciple-1")
		}
		c.Next()
	})
	NewExerciseHandler(fakeExerciseService{}, db).Register(api)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT role FROM "users"`)).
		WithArgs("disciple-1").
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("disciple"))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/exercises", bytes.NewReader([]byte(`{"name":"Bench","primary_muscle":"chest"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Test-User", "disciple")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("disciple create status=%d want 403", w.Code)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/exercises", nil)
	req.Header.Set("X-Test-User", "disciple")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("disciple list status=%d want 200", w.Code)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT role FROM "users"`)).
		WithArgs("coach-1").
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("coach"))
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/exercises", bytes.NewReader([]byte(`{"name":"Bench","primary_muscle":"chest"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Test-User", "coach")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("coach create status=%d body=%s", w.Code, w.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func mockGorm(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db, mock, func() { _ = sqlDB.Close() }
}
