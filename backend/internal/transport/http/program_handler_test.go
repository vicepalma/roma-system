package http

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/domain"
	"github.com/vicepalma/roma-system/backend/internal/repository"
	"github.com/vicepalma/roma-system/backend/internal/security"
	"github.com/vicepalma/roma-system/backend/internal/service"
)

type fakeProgramService struct{}

func (fakeProgramService) CreateProgram(context.Context, string, string, *string) (*domain.Program, error) {
	return &domain.Program{ID: "program-1", OwnerID: "coach-1", Title: "Base"}, nil
}
func (fakeProgramService) CreateProgramWithKind(_ context.Context, ownerID, title string, notes *string, kind string) (*domain.Program, error) {
	return &domain.Program{ID: "program-1", OwnerID: ownerID, Title: title, Notes: notes, Kind: kind}, nil
}
func (fakeProgramService) ListMyPrograms(context.Context, string, int, int) ([]domain.Program, int64, error) {
	return nil, 0, nil
}
func (fakeProgramService) AddWeek(context.Context, string, int) (*domain.ProgramWeek, error) {
	return nil, nil
}
func (fakeProgramService) DeleteWeek(context.Context, string, string) error { return nil }
func (fakeProgramService) AddDay(context.Context, string, int, *string) (*domain.ProgramDay, error) {
	return nil, nil
}
func (fakeProgramService) AddPrescription(context.Context, *domain.Prescription) (*domain.Prescription, error) {
	return nil, nil
}
func (fakeProgramService) Assign(context.Context, string, string, string, time.Time, *time.Time) (*domain.Assignment, error) {
	return nil, nil
}
func (fakeProgramService) MyToday(context.Context, string, time.Time) (*domain.ProgramDay, []domain.Prescription, error) {
	return nil, nil, nil
}
func (fakeProgramService) List(context.Context, repository.ProgramFilter) ([]repository.Program, int64, error) {
	return nil, 0, nil
}
func (fakeProgramService) Create(context.Context, string, service.CreateProgram) (*repository.Program, error) {
	return nil, nil
}
func (fakeProgramService) Get(context.Context, string) (*repository.Program, error) { return nil, nil }
func (fakeProgramService) Update(context.Context, string, service.UpdateProgram) (*repository.Program, error) {
	return &repository.Program{ID: "program-1", OwnerID: "coach-1", Title: "Updated"}, nil
}
func (fakeProgramService) Delete(context.Context, string) error { return nil }
func (fakeProgramService) NewVersion(context.Context, string) (*repository.ProgramVersion, error) {
	return nil, nil
}
func (fakeProgramService) ListVersions(context.Context, string) ([]repository.ProgramVersion, error) {
	return nil, nil
}
func (fakeProgramService) ListWeeks(context.Context, string) ([]repository.ProgramWeek, error) {
	return nil, nil
}
func (fakeProgramService) ListDays(context.Context, string) ([]repository.ProgramDay, error) {
	return nil, nil
}
func (fakeProgramService) UpdateDay(context.Context, string, *string, *int) (*repository.ProgramDay, error) {
	return nil, nil
}
func (fakeProgramService) DeleteDay(context.Context, string) error { return nil }
func (fakeProgramService) ListPrescriptions(context.Context, string) ([]repository.PrescriptionRow, error) {
	return nil, nil
}
func (fakeProgramService) UpdatePrescription(context.Context, string, service.UpdatePrescription) (*repository.Prescription, error) {
	return nil, nil
}
func (fakeProgramService) DeletePrescription(context.Context, string) error { return nil }
func (fakeProgramService) ReorderPrescriptions(context.Context, string, []string) error {
	return nil
}
func (fakeProgramService) GetProgram(context.Context, string) (*repository.ProgramRow, error) {
	return nil, nil
}
func (fakeProgramService) UpdateProgram(context.Context, string, *string, *string, *string) (*repository.ProgramRow, error) {
	return nil, nil
}
func (fakeProgramService) DeleteProgram(context.Context, string) error { return nil }
func (fakeProgramService) CreateNextVersionClone(context.Context, string) (*repository.ProgramRow, error) {
	return nil, nil
}
func (fakeProgramService) CreateSelfAssignment(context.Context, string, string, time.Time) (*domain.Assignment, error) {
	return &domain.Assignment{ID: "assignment-1"}, nil
}

func TestProgramPermissions(t *testing.T) {
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
		case "foreign-coach":
			c.Set(security.CtxUserID, "coach-2")
		}
		c.Next()
	})
	NewProgramHandler(fakeProgramService{}, db).Register(api)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT role FROM "users"`)).
		WithArgs("disciple-1").
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("disciple"))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/programs", bytes.NewReader([]byte(`{"title":"Base"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Test-User", "disciple")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("disciple create program status=%d want 403", w.Code)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT role FROM "users"`)).
		WithArgs("coach-1").
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("coach"))
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/programs", bytes.NewReader([]byte(`{"title":"Base"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Test-User", "coach")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("coach create program status=%d body=%s", w.Code, w.Body.String())
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT role FROM "users"`)).
		WithArgs("coach-2").
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("coach"))
	mock.ExpectQuery(`SELECT count\(\*\) FROM "programs"`).
		WithArgs("program-1", "coach-2").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPut, "/api/programs/program-1", bytes.NewReader([]byte(`{"title":"Updated"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Test-User", "foreign-coach")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("foreign coach update program status=%d want 403", w.Code)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT role FROM "users"`)).
		WithArgs("disciple-1").
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("disciple"))
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/programs", bytes.NewReader([]byte(`{"title":"Own Routine","kind":"self_training"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Test-User", "disciple")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("disciple create self-training status=%d body=%s", w.Code, w.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
