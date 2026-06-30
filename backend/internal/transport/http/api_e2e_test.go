package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/repository"
	"github.com/vicepalma/roma-system/backend/internal/security"
	"github.com/vicepalma/roma-system/backend/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	e2ePassword = "secret123"

	e2eCoach1    = "10000000-0000-0000-0000-000000000001"
	e2eCoach2    = "10000000-0000-0000-0000-000000000002"
	e2eDisciple1 = "10000000-0000-0000-0000-000000000011"
	e2eDisciple2 = "10000000-0000-0000-0000-000000000012"
	e2eDisciple3 = "10000000-0000-0000-0000-000000000013"
)

func TestE2EAPIPermissionsWithCleanDB(t *testing.T) {
	dsn := os.Getenv("ROMA_E2E_DB_URL")
	if dsn == "" {
		t.Skip("set ROMA_E2E_DB_URL to run E2E API tests against a migrated test Postgres database")
	}
	requireSafeE2EDSN(t, dsn)
	t.Setenv("JWT_SECRET", "e2e-test-secret")
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open e2e db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()

	cleanAndSeedE2EDB(t, db)
	r := e2eRouter(db)

	coach1Token, coach1ID := e2eLogin(t, r, "coach1.e2e@example.test")
	coach2Token, _ := e2eLogin(t, r, "coach2.e2e@example.test")
	disciple1Token, disciple1ID := e2eLogin(t, r, "disciple1.e2e@example.test")
	disciple2Token, disciple2ID := e2eLogin(t, r, "disciple2.e2e@example.test")
	disciple3Token, disciple3ID := e2eLogin(t, r, "disciple3.e2e@example.test")

	e2eAssertMe(t, r, coach1Token, "coach")
	e2eAssertMe(t, r, disciple1Token, "disciple")
	e2eRequest(t, r, http.MethodGet, "/api/exercises", "", nil, http.StatusUnauthorized)

	exerciseID := e2eCreateExercise(t, r, coach1Token, "E2E Bench Press")
	e2eRequest(t, r, http.MethodPost, "/api/exercises", disciple1Token, gin.H{"name": "E2E Disciple Forbidden", "primary_muscle": "chest"}, http.StatusForbidden)
	e2eRequest(t, r, http.MethodGet, "/api/exercises", disciple1Token, nil, http.StatusOK)

	programID := e2eCreateProgram(t, r, coach1Token, "E2E Coach Program")
	foreignProgramID := e2eCreateProgram(t, r, coach2Token, "E2E Foreign Program")
	e2eRequest(t, r, http.MethodPost, "/api/programs", disciple1Token, gin.H{"title": "Disciple Program"}, http.StatusForbidden)
	e2eRequest(t, r, http.MethodPut, "/api/programs/"+programID, coach2Token, gin.H{"title": "Foreign Mutation"}, http.StatusForbidden)

	weekID := e2ePostID(t, r, http.MethodPost, "/api/programs/"+programID+"/weeks", coach1Token, gin.H{"week_index": 1}, http.StatusCreated)
	dayID := e2ePostID(t, r, http.MethodPost, "/api/programs/"+programID+"/weeks/"+weekID+"/days", coach1Token, gin.H{"day_index": 1}, http.StatusCreated)
	prescriptionID := e2ePostID(t, r, http.MethodPost, "/api/programs/days/"+dayID+"/prescriptions", coach1Token, gin.H{
		"exercise_id": exerciseID,
		"series":      3,
		"reps":        "8-10",
		"position":    1,
	}, http.StatusCreated)

	foreignWeekID := e2ePostID(t, r, http.MethodPost, "/api/programs/"+foreignProgramID+"/weeks", coach2Token, gin.H{"week_index": 1}, http.StatusCreated)
	foreignDayID := e2ePostID(t, r, http.MethodPost, "/api/programs/"+foreignProgramID+"/weeks/"+foreignWeekID+"/days", coach2Token, gin.H{"day_index": 1}, http.StatusCreated)
	foreignPrescriptionID := e2ePostID(t, r, http.MethodPost, "/api/programs/days/"+foreignDayID+"/prescriptions", coach2Token, gin.H{
		"exercise_id": exerciseID,
		"series":      2,
		"reps":        "12",
		"position":    1,
	}, http.StatusCreated)

	assignmentID := e2ePostID(t, r, http.MethodPost, "/api/coach/assignments", coach1Token, gin.H{
		"disciple_id": disciple1ID,
		"program_id":  programID,
		"start_date":  "2026-06-29",
	}, http.StatusCreated)
	e2eRequest(t, r, http.MethodPost, "/api/coach/assignments", coach1Token, gin.H{
		"disciple_id": e2eDisciple2,
		"program_id":  programID,
		"start_date":  "2026-06-29",
	}, http.StatusForbidden)
	e2eRequest(t, r, http.MethodPost, "/api/coach/assignments", coach1Token, gin.H{
		"disciple_id": disciple1ID,
		"program_id":  foreignProgramID,
		"start_date":  "2026-06-29",
	}, http.StatusForbidden)
	e2eRequest(t, r, http.MethodPost, "/api/coach/assignments", disciple1Token, gin.H{
		"disciple_id": disciple1ID,
		"program_id":  programID,
		"start_date":  "2026-06-29",
	}, http.StatusForbidden)
	e2eRequest(t, r, http.MethodPost, "/api/coach/assignments/"+assignmentID+"/activate?disciple_id="+disciple1ID, coach1Token, nil, http.StatusNoContent)

	sessionID := e2ePostID(t, r, http.MethodPost, "/api/sessions", disciple1Token, gin.H{
		"assignment_id": assignmentID,
		"day_id":        dayID,
	}, http.StatusCreated)
	e2eRequest(t, r, http.MethodPost, "/api/sessions", disciple2Token, gin.H{
		"assignment_id": assignmentID,
		"day_id":        dayID,
	}, http.StatusForbidden)
	e2eRequest(t, r, http.MethodPost, "/api/sessions", disciple1Token, gin.H{
		"assignment_id": assignmentID,
		"day_id":        foreignDayID,
	}, http.StatusBadRequest)
	e2eRequest(t, r, http.MethodGet, "/api/sessions/"+sessionID, coach1Token, nil, http.StatusOK)
	e2eRequest(t, r, http.MethodGet, "/api/sessions/"+sessionID, coach2Token, nil, http.StatusForbidden)

	setID := e2ePostID(t, r, http.MethodPost, "/api/sessions/"+sessionID+"/sets", disciple1Token, gin.H{
		"prescription_id": prescriptionID,
		"set_index":       1,
		"weight":          100,
		"reps":            8,
	}, http.StatusCreated)
	e2eRequest(t, r, http.MethodPost, "/api/sessions/"+sessionID+"/sets", disciple1Token, gin.H{
		"prescription_id": foreignPrescriptionID,
		"set_index":       2,
		"reps":            10,
	}, http.StatusBadRequest)
	e2eRequest(t, r, http.MethodDelete, "/api/sessions/"+sessionID+"/sets/"+setID, disciple2Token, nil, http.StatusForbidden)

	e2eSetAssignmentActive(t, db, assignmentID, false)
	e2eRequest(t, r, http.MethodPost, "/api/sessions", disciple1Token, gin.H{
		"assignment_id": assignmentID,
		"day_id":        dayID,
	}, http.StatusConflict)
	e2eRequest(t, r, http.MethodGet, "/api/sessions/"+sessionID, disciple1Token, nil, http.StatusOK)
	e2eRequest(t, r, http.MethodGet, "/api/sessions/"+sessionID, coach1Token, nil, http.StatusOK)

	e2eRequest(t, r, http.MethodGet, "/api/history", disciple1Token, nil, http.StatusOK)
	e2eRequest(t, r, http.MethodGet, "/api/history?disciple_id="+disciple1ID, coach1Token, nil, http.StatusOK)
	e2eRequest(t, r, http.MethodGet, "/api/history?disciple_id="+disciple1ID, coach2Token, nil, http.StatusForbidden)
	e2eRequest(t, r, http.MethodGet, "/api/history?disciple_id="+disciple1ID, disciple2Token, nil, http.StatusForbidden)

	selfProgramID := e2eCreateSelfTrainingProgram(t, r, disciple3Token, "E2E Self Routine")
	selfProgramBID := e2eCreateSelfTrainingProgram(t, r, disciple3Token, "E2E Self Routine B")
	e2eRequest(t, r, http.MethodGet, "/api/programs/"+selfProgramID, disciple3Token, nil, http.StatusOK)
	e2eRequest(t, r, http.MethodGet, "/api/programs/"+selfProgramID, disciple2Token, nil, http.StatusForbidden)
	e2eRequest(t, r, http.MethodGet, "/api/programs/"+selfProgramID, coach1Token, nil, http.StatusForbidden)
	e2eRequest(t, r, http.MethodPost, "/api/programs/"+programID+"/self-assignment", disciple3Token, gin.H{"start_date": "2026-06-29"}, http.StatusForbidden)
	e2eRequest(t, r, http.MethodPost, "/api/programs/"+selfProgramID+"/self-assignment", coach1Token, gin.H{"start_date": "2026-06-29"}, http.StatusForbidden)
	selfWeekID := e2ePostID(t, r, http.MethodPost, "/api/programs/"+selfProgramID+"/weeks", disciple3Token, gin.H{"week_index": 1}, http.StatusCreated)
	selfDayID := e2ePostID(t, r, http.MethodPost, "/api/programs/"+selfProgramID+"/weeks/"+selfWeekID+"/days", disciple3Token, gin.H{"day_index": 1}, http.StatusCreated)
	e2ePostID(t, r, http.MethodPost, "/api/programs/days/"+selfDayID+"/prescriptions", disciple3Token, gin.H{
		"exercise_id": exerciseID,
		"series":      3,
		"reps":        "10",
		"position":    1,
	}, http.StatusCreated)
	selfWeekBID := e2ePostID(t, r, http.MethodPost, "/api/programs/"+selfProgramBID+"/weeks", disciple3Token, gin.H{"week_index": 1}, http.StatusCreated)
	selfDayBID := e2ePostID(t, r, http.MethodPost, "/api/programs/"+selfProgramBID+"/weeks/"+selfWeekBID+"/days", disciple3Token, gin.H{"day_index": 1}, http.StatusCreated)
	selfPrescriptionBID := e2ePostID(t, r, http.MethodPost, "/api/programs/days/"+selfDayBID+"/prescriptions", disciple3Token, gin.H{
		"exercise_id": exerciseID,
		"series":      4,
		"reps":        "8",
		"position":    1,
	}, http.StatusCreated)
	e2eRequest(t, r, http.MethodPut, "/api/programs/"+selfProgramID, disciple2Token, gin.H{"title": "Steal Routine"}, http.StatusForbidden)
	e2eRequest(t, r, http.MethodPost, "/api/coach/assignments", disciple3Token, gin.H{
		"disciple_id": disciple2ID,
		"program_id":  selfProgramID,
		"start_date":  "2026-06-29",
	}, http.StatusForbidden)

	e2eInsertActiveCoachAssignment(t, db, programID, disciple3ID, e2eCoach1)
	selfAssignmentID := e2ePostID(t, r, http.MethodPost, "/api/programs/"+selfProgramID+"/self-assignment", disciple3Token, gin.H{"start_date": "2026-06-29"}, http.StatusCreated)
	e2eAssertSelfAssignmentState(t, db, disciple3ID, selfProgramID, true)
	e2eAssertActiveSelfAssignmentCount(t, db, disciple3ID, 1)
	e2eAssertActiveCoachAssignmentCount(t, db, disciple3ID, 1)
	e2eAssertActiveAssignment(t, r, disciple3Token, selfAssignmentID)
	e2eAssertAssignmentDays(t, r, disciple3Token, selfAssignmentID, selfDayID)
	selfAssignmentBID := e2ePostID(t, r, http.MethodPost, "/api/programs/"+selfProgramBID+"/self-assignment", disciple3Token, gin.H{"start_date": "2026-06-30"}, http.StatusCreated)
	if selfAssignmentBID == selfAssignmentID {
		t.Fatalf("second self-training activation reused first assignment id")
	}
	e2eAssertSelfAssignmentState(t, db, disciple3ID, selfProgramID, false)
	e2eAssertSelfAssignmentState(t, db, disciple3ID, selfProgramBID, true)
	e2eAssertActiveSelfAssignmentCount(t, db, disciple3ID, 1)
	e2eAssertActiveCoachAssignmentCount(t, db, disciple3ID, 1)
	e2eAssertAssignmentDays(t, r, disciple3Token, selfAssignmentBID, selfDayBID)

	e2eRequest(t, r, http.MethodPost, "/api/sessions", disciple3Token, gin.H{
		"assignment_id": selfAssignmentID,
		"day_id":        selfDayID,
	}, http.StatusConflict)
	e2eRequest(t, r, http.MethodPost, "/api/sessions", disciple3Token, gin.H{
		"assignment_id": selfAssignmentBID,
		"day_id":        selfDayID,
	}, http.StatusBadRequest)
	selfSessionID := e2ePostID(t, r, http.MethodPost, "/api/sessions", disciple3Token, gin.H{
		"assignment_id": selfAssignmentBID,
		"day_id":        selfDayBID,
	}, http.StatusCreated)
	e2ePostID(t, r, http.MethodPost, "/api/sessions/"+selfSessionID+"/sets", disciple3Token, gin.H{
		"prescription_id": selfPrescriptionBID,
		"set_index":       1,
		"reps":            10,
	}, http.StatusCreated)
	e2eRequest(t, r, http.MethodGet, "/api/sessions/"+selfSessionID, disciple2Token, nil, http.StatusForbidden)
	e2eRequest(t, r, http.MethodGet, "/api/sessions/"+selfSessionID, coach1Token, nil, http.StatusForbidden)
	e2eRequest(t, r, http.MethodGet, "/api/history?disciple_id="+disciple3ID, disciple3Token, nil, http.StatusOK)
	e2eRequest(t, r, http.MethodGet, "/api/history?disciple_id="+disciple3ID, coach1Token, nil, http.StatusForbidden)

	if coach1ID != e2eCoach1 {
		t.Fatalf("coach id=%s want %s", coach1ID, e2eCoach1)
	}
}

func e2eRouter(db *gorm.DB) *gin.Engine {
	userRepo := repository.NewUserRepository(db)
	exRepo := repository.NewExerciseRepository(db)
	progRepo := repository.NewProgramRepository(db)
	assignRepo := repository.NewAssignmentRepository(db)
	histRepo := repository.NewHistoryRepository(db)
	coachRepo := repository.NewCoachRepository(db)
	sessRepo := repository.NewSessionRepository(db)
	inviteRepo := repository.NewInviteRepository(db)
	adRepo := repository.NewAssignmentDaysRepository(db)

	histSvc := service.NewHistoryService(histRepo)
	coachSvc := service.NewCoachService(coachRepo, histSvc, db, assignRepo)
	sessSvc := service.NewSessionService(sessRepo, coachSvc)

	r := gin.New()
	NewAuthHandler(userRepo, db).Register(r.Group("/"))
	api := r.Group("/api", security.AuthRequired())
	NewExerciseHandler(service.NewExerciseService(exRepo), db).Register(api)
	NewProgramHandler(service.NewProgramService(progRepo), db).Register(api)
	NewSessionHandler(sessSvc, db).Register(api)
	NewHistoryHandler(histSvc, "UTC", db).Register(api)
	NewCoachHandler(coachSvc, histSvc, userRepo, db).Register(api)
	NewInviteHandler(service.NewInviteService(inviteRepo, coachSvc, "")).Register(api)
	NewAssignmentDaysHandler(service.NewAssignmentDaysService(adRepo, coachSvc)).Register(api)
	NewMeHandler(histSvc, coachSvc, sessSvc).Register(api)
	return r
}

func cleanAndSeedE2EDB(t *testing.T, db *gorm.DB) {
	t.Helper()
	for _, table := range []string{
		"set_logs", "cardio_segments", "session_logs", "assignments", "prescriptions",
		"program_days", "program_weeks", "program_versions", "programs", "exercises",
		"coach_links", "master_disciple", "invite_codes", "invitations", "checkins",
		"user_flags", "methods", "users",
	} {
		if err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE").Error; err != nil {
			t.Fatalf("truncate %s: %v", table, err)
		}
	}

	hash, err := security.HashPassword(e2ePassword)
	if err != nil {
		t.Fatal(err)
	}
	users := []struct {
		id, email, name, role string
	}{
		{e2eCoach1, "coach1.e2e@example.test", "E2E Coach 1", "coach"},
		{e2eCoach2, "coach2.e2e@example.test", "E2E Coach 2", "coach"},
		{e2eDisciple1, "disciple1.e2e@example.test", "E2E Disciple 1", "disciple"},
		{e2eDisciple2, "disciple2.e2e@example.test", "E2E Disciple 2", "disciple"},
		{e2eDisciple3, "disciple3.e2e@example.test", "E2E Disciple 3", "disciple"},
	}
	for _, u := range users {
		if err := db.Exec(
			`INSERT INTO users (id, email, password_hash, name, role) VALUES (?, ?, ?, ?, ?)`,
			u.id, u.email, hash, u.name, u.role,
		).Error; err != nil {
			t.Fatalf("insert user %s: %v", u.email, err)
		}
	}
	for _, link := range []struct{ coachID, discipleID string }{
		{e2eCoach1, e2eDisciple1},
		{e2eCoach2, e2eDisciple2},
	} {
		if err := db.Exec(
			`INSERT INTO coach_links (coach_id, disciple_id, status) VALUES (?, ?, 'accepted')`,
			link.coachID, link.discipleID,
		).Error; err != nil {
			t.Fatalf("insert coach link: %v", err)
		}
	}
}

func e2eLogin(t *testing.T, r http.Handler, email string) (string, string) {
	t.Helper()
	resp := e2eRequest(t, r, http.MethodPost, "/auth/login", "", gin.H{"email": email, "password": e2ePassword}, http.StatusOK)
	var out struct {
		User struct {
			ID   string `json:"id"`
			Role string `json:"role"`
		} `json:"user"`
		Tokens struct {
			Access string `json:"access"`
		} `json:"tokens"`
	}
	e2eDecode(t, resp, &out)
	if out.Tokens.Access == "" {
		t.Fatalf("empty access token for %s", email)
	}
	return out.Tokens.Access, out.User.ID
}

func e2eAssertMe(t *testing.T, r http.Handler, token string, role string) {
	t.Helper()
	resp := e2eRequest(t, r, http.MethodGet, "/me", token, nil, http.StatusOK)
	var out struct {
		Role string `json:"role"`
	}
	e2eDecode(t, resp, &out)
	if out.Role != role {
		t.Fatalf("/me role=%q want %q", out.Role, role)
	}
}

func e2eCreateExercise(t *testing.T, r http.Handler, token, name string) string {
	t.Helper()
	return e2ePostID(t, r, http.MethodPost, "/api/exercises", token, gin.H{"name": name, "primary_muscle": "chest"}, http.StatusCreated)
}

func e2eCreateProgram(t *testing.T, r http.Handler, token, title string) string {
	t.Helper()
	return e2ePostID(t, r, http.MethodPost, "/api/programs", token, gin.H{"title": title}, http.StatusCreated)
}

func e2eCreateSelfTrainingProgram(t *testing.T, r http.Handler, token, title string) string {
	t.Helper()
	return e2ePostID(t, r, http.MethodPost, "/api/programs", token, gin.H{"title": title, "kind": "self_training"}, http.StatusCreated)
}

func e2eAssertActiveAssignment(t *testing.T, r http.Handler, token, assignmentID string) {
	t.Helper()
	resp := e2eRequest(t, r, http.MethodGet, "/api/me/assignment/active", token, nil, http.StatusOK)
	var out struct {
		ID string `json:"id"`
	}
	e2eDecode(t, resp, &out)
	if out.ID != assignmentID {
		t.Fatalf("active assignment id=%s want %s", out.ID, assignmentID)
	}
}

func e2eAssertAssignmentDays(t *testing.T, r http.Handler, token, assignmentID, dayID string) {
	t.Helper()
	resp := e2eRequest(t, r, http.MethodGet, "/api/assignments/"+assignmentID+"/days", token, nil, http.StatusOK)
	var out struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	}
	e2eDecode(t, resp, &out)
	for _, item := range out.Items {
		if item.ID == dayID {
			return
		}
	}
	t.Fatalf("assignment %s days did not include day %s; got %#v", assignmentID, dayID, out.Items)
}

func e2eSetAssignmentActive(t *testing.T, db *gorm.DB, assignmentID string, active bool) {
	t.Helper()
	if err := db.Exec(`UPDATE assignments SET is_active = ? WHERE id = ?`, active, assignmentID).Error; err != nil {
		t.Fatalf("set assignment active=%v: %v", active, err)
	}
}

func e2eInsertActiveCoachAssignment(t *testing.T, db *gorm.DB, programID, discipleID, coachID string) {
	t.Helper()
	if err := db.Exec(`
		INSERT INTO assignments (program_id, program_version, disciple_id, assigned_by, start_date, is_active)
		VALUES (?, 1, ?, ?, '2026-06-29', true)
	`, programID, discipleID, coachID).Error; err != nil {
		t.Fatalf("insert active coach assignment: %v", err)
	}
}

func e2eAssertSelfAssignmentState(t *testing.T, db *gorm.DB, discipleID, programID string, active bool) {
	t.Helper()
	var got bool
	if err := db.Raw(`
		SELECT is_active
		FROM assignments
		WHERE disciple_id = ? AND assigned_by = ? AND program_id = ?
	`, discipleID, discipleID, programID).Scan(&got).Error; err != nil {
		t.Fatalf("read self assignment state: %v", err)
	}
	if got != active {
		t.Fatalf("self assignment active=%v want %v for program %s", got, active, programID)
	}
}

func e2eAssertActiveSelfAssignmentCount(t *testing.T, db *gorm.DB, discipleID string, want int64) {
	t.Helper()
	var got int64
	if err := db.Raw(`
		SELECT COUNT(*)
		FROM assignments a
		JOIN programs p ON p.id = a.program_id
		WHERE a.disciple_id = ?
		  AND a.assigned_by = ?
		  AND a.is_active = true
		  AND p.kind = 'self_training'
	`, discipleID, discipleID).Scan(&got).Error; err != nil {
		t.Fatalf("count active self assignments: %v", err)
	}
	if got != want {
		t.Fatalf("active self assignments=%d want %d", got, want)
	}
}

func e2eAssertActiveCoachAssignmentCount(t *testing.T, db *gorm.DB, discipleID string, want int64) {
	t.Helper()
	var got int64
	if err := db.Raw(`
		SELECT COUNT(*)
		FROM assignments
		WHERE disciple_id = ?
		  AND assigned_by <> disciple_id
		  AND is_active = true
	`, discipleID).Scan(&got).Error; err != nil {
		t.Fatalf("count active coach assignments: %v", err)
	}
	if got != want {
		t.Fatalf("active coach assignments=%d want %d", got, want)
	}
}

func e2ePostID(t *testing.T, r http.Handler, method, path, token string, body any, want int) string {
	t.Helper()
	resp := e2eRequest(t, r, method, path, token, body, want)
	var out struct {
		ID string `json:"id"`
	}
	e2eDecode(t, resp, &out)
	if out.ID == "" {
		t.Fatalf("%s %s returned empty id: %s", method, path, resp)
	}
	return out.ID
}

func e2eRequest(t *testing.T, r http.Handler, method, path, token string, body any, want int) []byte {
	t.Helper()
	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		reader = bytes.NewReader(raw)
	}
	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if strings.TrimSpace(token) != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != want {
		t.Fatalf("%s %s status=%d want=%d body=%s", method, path, w.Code, want, w.Body.String())
	}
	return w.Body.Bytes()
}

func e2eDecode(t *testing.T, raw []byte, out any) {
	t.Helper()
	if err := json.Unmarshal(raw, out); err != nil {
		t.Fatalf("decode %s: %v", string(raw), err)
	}
}

func requireSafeE2EDSN(t *testing.T, dsn string) {
	t.Helper()
	if strings.Contains(dsn, "sslmode=require") {
		t.Skip("external managed databases are not allowed for destructive E2E tests")
	}
	if !strings.Contains(dsn, "localhost") && !strings.Contains(dsn, "127.0.0.1") {
		t.Fatalf("ROMA_E2E_DB_URL must point to a local test database, got %q", redactDSN(dsn))
	}
}

func redactDSN(dsn string) string {
	if at := strings.LastIndex(dsn, "@"); at >= 0 {
		return fmt.Sprintf("postgres://<redacted>%s", dsn[at:])
	}
	return dsn
}
