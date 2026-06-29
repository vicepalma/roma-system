package security

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestRoleAndCoachDiscipleAccess(t *testing.T) {
	db, mock, cleanup := mockGorm(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT role FROM "users"`)).
		WithArgs("coach-1").
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("coach"))
	role, err := RoleOf(db, "coach-1")
	if err != nil || role != "coach" {
		t.Fatalf("RoleOf role=%q err=%v", role, err)
	}

	ok, err := CanAccessDisciple(db, "disciple-1", "disciple-1")
	if err != nil || !ok {
		t.Fatalf("self access ok=%v err=%v", ok, err)
	}

	mock.ExpectQuery(`SELECT count\(\*\) FROM "coach_links"`).
		WithArgs("coach-1", "disciple-1").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	ok, err = CanAccessDisciple(db, "coach-1", "disciple-1")
	if err != nil || !ok {
		t.Fatalf("linked coach access ok=%v err=%v", ok, err)
	}

	mock.ExpectQuery(`SELECT count\(\*\) FROM "coach_links"`).
		WithArgs("coach-2", "disciple-1").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	ok, err = CanAccessDisciple(db, "coach-2", "disciple-1")
	if err != nil || ok {
		t.Fatalf("foreign coach access ok=%v err=%v", ok, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestOwnershipGuards(t *testing.T) {
	db, mock, cleanup := mockGorm(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT count\(\*\) FROM "programs"`).
		WithArgs("program-1", "coach-1").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	ok, err := IsProgramOwner(db, "coach-1", "program-1")
	if err != nil || !ok {
		t.Fatalf("program owner ok=%v err=%v", ok, err)
	}

	mock.ExpectQuery(`SELECT count\(\*\) FROM "assignments"`).
		WithArgs("assignment-1", "disciple-1").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	ok, err = IsAssignmentOwnedByDisciple(db, "disciple-1", "assignment-1")
	if err != nil || !ok {
		t.Fatalf("assignment owner ok=%v err=%v", ok, err)
	}

	mock.ExpectQuery(`SELECT count\(\*\) FROM "session_logs"`).
		WithArgs("session-1", "disciple-1").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	ok, err = IsSessionOwnedByDisciple(db, "disciple-1", "session-1")
	if err != nil || !ok {
		t.Fatalf("session owner ok=%v err=%v", ok, err)
	}

	mock.ExpectQuery(`SELECT count\(\*\) FROM set_logs AS st`).
		WithArgs("set-1", "disciple-1").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	ok, err = IsSetOwnedByDisciple(db, "disciple-1", "set-1")
	if err != nil || !ok {
		t.Fatalf("set owner ok=%v err=%v", ok, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestSelfTrainingProgramMutability(t *testing.T) {
	db, mock, cleanup := mockGorm(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT role FROM "users"`)).
		WithArgs("disciple-1").
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("disciple"))
	mock.ExpectQuery(`SELECT count\(\*\) FROM "programs"`).
		WithArgs("program-1", "disciple-1").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	ok, err := IsProgramMutable(db, "disciple-1", "program-1")
	if err != nil || !ok {
		t.Fatalf("disciple self-training mutable ok=%v err=%v", ok, err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT role FROM "users"`)).
		WithArgs("coach-1").
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("coach"))
	mock.ExpectQuery(`SELECT count\(\*\) FROM "programs"`).
		WithArgs("program-1", "coach-1").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	ok, err = IsProgramMutable(db, "coach-1", "program-1")
	if err != nil || ok {
		t.Fatalf("foreign coach self-training mutable ok=%v err=%v", ok, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestResourceAccessThroughDiscipleLink(t *testing.T) {
	db, mock, cleanup := mockGorm(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT disciple_id FROM "assignments"`)).
		WithArgs("assignment-1").
		WillReturnRows(sqlmock.NewRows([]string{"disciple_id"}).AddRow("disciple-1"))
	mock.ExpectQuery(`SELECT count\(\*\) FROM "coach_links"`).
		WithArgs("coach-1", "disciple-1").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	ok, err := CanAccessAssignment(db, "coach-1", "assignment-1")
	if err != nil || !ok {
		t.Fatalf("linked coach assignment access ok=%v err=%v", ok, err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT disciple_id FROM "session_logs"`)).
		WithArgs("session-1").
		WillReturnRows(sqlmock.NewRows([]string{"disciple_id"}).AddRow("disciple-1"))
	mock.ExpectQuery(`SELECT count\(\*\) FROM "coach_links"`).
		WithArgs("coach-2", "disciple-1").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	ok, err = CanAccessSession(db, "coach-2", "session-1")
	if err != nil || ok {
		t.Fatalf("foreign coach session access ok=%v err=%v", ok, err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT s.disciple_id FROM set_logs AS st`)).
		WithArgs("set-1").
		WillReturnRows(sqlmock.NewRows([]string{"disciple_id"}).AddRow("disciple-1"))
	ok, err = CanAccessSet(db, "disciple-1", "set-1")
	if err != nil || !ok {
		t.Fatalf("disciple set access ok=%v err=%v", ok, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestSessionConsistencyGuards(t *testing.T) {
	db, mock, cleanup := mockGorm(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT count\(\*\) FROM assignments AS a`).
		WithArgs("assignment-1", "day-1").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	ok, err := IsAssignmentDay(db, "assignment-1", "day-1")
	if err != nil || !ok {
		t.Fatalf("assignment day ok=%v err=%v", ok, err)
	}

	mock.ExpectQuery(`SELECT count\(\*\) FROM session_logs AS s`).
		WithArgs("session-1", "prescription-foreign").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	ok, err = IsPrescriptionInSessionDay(db, "session-1", "prescription-foreign")
	if err != nil || ok {
		t.Fatalf("foreign prescription ok=%v err=%v", ok, err)
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
