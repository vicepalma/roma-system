package security

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var ErrForbidden = errors.New("forbidden")

func RoleOf(db *gorm.DB, userID string) (string, error) {
	var role string
	err := db.Table("users").Select("role").Where("id = ?", userID).Scan(&role).Error
	if err != nil {
		return "", err
	}
	if role == "" {
		return "", gorm.ErrRecordNotFound
	}
	return role, nil
}

func RequireRole(db *gorm.DB, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := UserID(c)
		if uid == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		got, err := RoleOf(db.WithContext(c.Request.Context()), uid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
			return
		}
		if got != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}

func IsCoachOf(db *gorm.DB, coachID, discipleID string) (bool, error) {
	if coachID == "" || discipleID == "" {
		return false, nil
	}
	var count int64
	err := db.Table("coach_links").
		Where("coach_id = ? AND disciple_id = ? AND status = 'accepted'", coachID, discipleID).
		Count(&count).Error
	return count > 0, err
}

func CanAccessDisciple(db *gorm.DB, actorID, discipleID string) (bool, error) {
	if actorID == "" || discipleID == "" {
		return false, nil
	}
	if actorID == discipleID {
		return true, nil
	}
	return IsCoachOf(db, actorID, discipleID)
}

func RequireSelfOrCoachOf(db *gorm.DB, param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		discipleID := c.Param(param)
		if discipleID == "" {
			discipleID = c.Query(param)
		}
		ok, err := CanAccessDisciple(db.WithContext(c.Request.Context()), UserID(c), discipleID)
		abortAccess(c, ok, err)
	}
}

func IsProgramOwner(db *gorm.DB, actorID, programID string) (bool, error) {
	var count int64
	err := db.Table("programs").Where("id = ? AND owner_id = ?", programID, actorID).Count(&count).Error
	return count > 0, err
}

func ProgramKind(db *gorm.DB, programID string) (string, error) {
	var kind string
	err := db.Table("programs").Select("kind").Where("id = ?", programID).Scan(&kind).Error
	if err != nil {
		return "", err
	}
	if kind == "" {
		return "", gorm.ErrRecordNotFound
	}
	return kind, nil
}

func IsProgramMutable(db *gorm.DB, actorID, programID string) (bool, error) {
	role, err := RoleOf(db, actorID)
	if err != nil {
		return false, err
	}
	var count int64
	switch role {
	case "coach":
		err = db.Table("programs").
			Where("id = ? AND owner_id = ? AND kind = 'coach_program'", programID, actorID).
			Count(&count).Error
	case "disciple":
		err = db.Table("programs").
			Where("id = ? AND owner_id = ? AND kind = 'self_training'", programID, actorID).
			Count(&count).Error
	default:
		return false, nil
	}
	return count > 0, err
}

func IsProgramReadable(db *gorm.DB, actorID, programID string) (bool, error) {
	ok, err := IsProgramOwner(db, actorID, programID)
	if err != nil || ok {
		return ok, err
	}
	kind, err := ProgramKind(db, programID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	if kind == "self_training" {
		return false, nil
	}
	var count int64
	err = db.Table("assignments AS a").
		Joins("LEFT JOIN coach_links cl ON cl.disciple_id = a.disciple_id AND cl.coach_id = ? AND cl.status = 'accepted'", actorID).
		Where("a.program_id = ? AND (a.disciple_id = ? OR cl.id IS NOT NULL)", programID, actorID).
		Count(&count).Error
	return count > 0, err
}

func IsProgramOwnerByWeek(db *gorm.DB, actorID, weekID string) (bool, error) {
	var count int64
	err := db.Table("program_weeks AS w").
		Joins("JOIN programs p ON p.id = w.program_id").
		Where("w.id = ? AND p.owner_id = ?", weekID, actorID).
		Count(&count).Error
	return count > 0, err
}

func IsProgramMutableByWeek(db *gorm.DB, actorID, weekID string) (bool, error) {
	var programID string
	err := db.Table("program_weeks").Select("program_id").Where("id = ?", weekID).Scan(&programID).Error
	if err != nil {
		return false, err
	}
	if programID == "" {
		return false, nil
	}
	return IsProgramMutable(db, actorID, programID)
}

func IsProgramOwnerByDay(db *gorm.DB, actorID, dayID string) (bool, error) {
	var count int64
	err := db.Table("program_days AS d").
		Joins("JOIN program_weeks w ON w.id = d.week_id").
		Joins("JOIN programs p ON p.id = w.program_id").
		Where("d.id = ? AND p.owner_id = ?", dayID, actorID).
		Count(&count).Error
	return count > 0, err
}

func IsProgramMutableByDay(db *gorm.DB, actorID, dayID string) (bool, error) {
	var programID string
	err := db.Table("program_days AS d").
		Select("w.program_id").
		Joins("JOIN program_weeks w ON w.id = d.week_id").
		Where("d.id = ?", dayID).
		Scan(&programID).Error
	if err != nil {
		return false, err
	}
	if programID == "" {
		return false, nil
	}
	return IsProgramMutable(db, actorID, programID)
}

func IsProgramOwnerByPrescription(db *gorm.DB, actorID, prescriptionID string) (bool, error) {
	var count int64
	err := db.Table("prescriptions AS pr").
		Joins("JOIN program_days d ON d.id = pr.day_id").
		Joins("JOIN program_weeks w ON w.id = d.week_id").
		Joins("JOIN programs p ON p.id = w.program_id").
		Where("pr.id = ? AND p.owner_id = ?", prescriptionID, actorID).
		Count(&count).Error
	return count > 0, err
}

func IsProgramMutableByPrescription(db *gorm.DB, actorID, prescriptionID string) (bool, error) {
	var programID string
	err := db.Table("prescriptions AS pr").
		Select("w.program_id").
		Joins("JOIN program_days d ON d.id = pr.day_id").
		Joins("JOIN program_weeks w ON w.id = d.week_id").
		Where("pr.id = ?", prescriptionID).
		Scan(&programID).Error
	if err != nil {
		return false, err
	}
	if programID == "" {
		return false, nil
	}
	return IsProgramMutable(db, actorID, programID)
}

func IsDayReadable(db *gorm.DB, actorID, dayID string) (bool, error) {
	ok, err := IsProgramOwnerByDay(db, actorID, dayID)
	if err != nil || ok {
		return ok, err
	}
	var kind string
	err = db.Table("program_days AS d").
		Select("p.kind").
		Joins("JOIN program_weeks w ON w.id = d.week_id").
		Joins("JOIN programs p ON p.id = w.program_id").
		Where("d.id = ?", dayID).
		Scan(&kind).Error
	if err != nil {
		return false, err
	}
	if kind == "self_training" {
		return false, nil
	}
	var count int64
	err = db.Table("program_days AS d").
		Joins("JOIN program_weeks w ON w.id = d.week_id").
		Joins("JOIN assignments a ON a.program_id = w.program_id").
		Joins("LEFT JOIN coach_links cl ON cl.disciple_id = a.disciple_id AND cl.coach_id = ? AND cl.status = 'accepted'", actorID).
		Where("d.id = ? AND (a.disciple_id = ? OR cl.id IS NOT NULL)", dayID, actorID).
		Count(&count).Error
	return count > 0, err
}

func RequireProgramOwner(db *gorm.DB, param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := IsProgramOwner(db.WithContext(c.Request.Context()), UserID(c), c.Param(param))
		abortAccess(c, ok, err)
	}
}

func RequireProgramMutable(db *gorm.DB, param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := IsProgramMutable(db.WithContext(c.Request.Context()), UserID(c), c.Param(param))
		abortAccess(c, ok, err)
	}
}

func RequireProgramOwnerByWeek(db *gorm.DB, param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := IsProgramOwnerByWeek(db.WithContext(c.Request.Context()), UserID(c), c.Param(param))
		abortAccess(c, ok, err)
	}
}

func RequireProgramMutableByWeek(db *gorm.DB, param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := IsProgramMutableByWeek(db.WithContext(c.Request.Context()), UserID(c), c.Param(param))
		abortAccess(c, ok, err)
	}
}

func RequireProgramOwnerByDay(db *gorm.DB, param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := IsProgramOwnerByDay(db.WithContext(c.Request.Context()), UserID(c), c.Param(param))
		abortAccess(c, ok, err)
	}
}

func RequireProgramMutableByDay(db *gorm.DB, param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := IsProgramMutableByDay(db.WithContext(c.Request.Context()), UserID(c), c.Param(param))
		abortAccess(c, ok, err)
	}
}

func RequireProgramOwnerByPrescription(db *gorm.DB, param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := IsProgramOwnerByPrescription(db.WithContext(c.Request.Context()), UserID(c), c.Param(param))
		abortAccess(c, ok, err)
	}
}

func RequireProgramMutableByPrescription(db *gorm.DB, param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := IsProgramMutableByPrescription(db.WithContext(c.Request.Context()), UserID(c), c.Param(param))
		abortAccess(c, ok, err)
	}
}

func RequireDayReadable(db *gorm.DB, param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := IsDayReadable(db.WithContext(c.Request.Context()), UserID(c), c.Param(param))
		abortAccess(c, ok, err)
	}
}

func CanAccessAssignment(db *gorm.DB, actorID, assignmentID string) (bool, error) {
	var row struct{ DiscipleID string }
	if err := db.Table("assignments").Select("disciple_id").Where("id = ?", assignmentID).Scan(&row).Error; err != nil {
		return false, err
	}
	if row.DiscipleID == "" {
		return false, nil
	}
	return CanAccessDisciple(db, actorID, row.DiscipleID)
}

func IsAssignmentOwnedByDisciple(db *gorm.DB, actorID, assignmentID string) (bool, error) {
	var count int64
	err := db.Table("assignments").Where("id = ? AND disciple_id = ?", assignmentID, actorID).Count(&count).Error
	return count > 0, err
}

func IsAssignmentActive(db *gorm.DB, assignmentID string) (bool, error) {
	var count int64
	err := db.Table("assignments").Where("id = ? AND is_active = true", assignmentID).Count(&count).Error
	return count > 0, err
}

func CanAccessSession(db *gorm.DB, actorID, sessionID string) (bool, error) {
	var row struct{ DiscipleID string }
	if err := db.Table("session_logs").Select("disciple_id").Where("id = ?", sessionID).Scan(&row).Error; err != nil {
		return false, err
	}
	if row.DiscipleID == "" {
		return false, nil
	}
	return CanAccessDisciple(db, actorID, row.DiscipleID)
}

func IsSessionOwnedByDisciple(db *gorm.DB, actorID, sessionID string) (bool, error) {
	var count int64
	err := db.Table("session_logs").Where("id = ? AND disciple_id = ?", sessionID, actorID).Count(&count).Error
	return count > 0, err
}

func CanAccessSet(db *gorm.DB, actorID, setID string) (bool, error) {
	var row struct{ DiscipleID string }
	err := db.Table("set_logs AS st").
		Select("s.disciple_id").
		Joins("JOIN session_logs s ON s.id = st.session_id").
		Where("st.id = ?", setID).
		Scan(&row).Error
	if err != nil {
		return false, err
	}
	if row.DiscipleID == "" {
		return false, nil
	}
	return CanAccessDisciple(db, actorID, row.DiscipleID)
}

func IsSetOwnedByDisciple(db *gorm.DB, actorID, setID string) (bool, error) {
	var count int64
	err := db.Table("set_logs AS st").
		Joins("JOIN session_logs s ON s.id = st.session_id").
		Where("st.id = ? AND s.disciple_id = ?", setID, actorID).
		Count(&count).Error
	return count > 0, err
}

func IsAssignmentDay(db *gorm.DB, assignmentID, dayID string) (bool, error) {
	var count int64
	err := db.Table("assignments AS a").
		Joins("JOIN program_weeks w ON w.program_id = a.program_id").
		Joins("JOIN program_days d ON d.week_id = w.id").
		Where("a.id = ? AND d.id = ?", assignmentID, dayID).
		Count(&count).Error
	return count > 0, err
}

func IsPrescriptionInSessionDay(db *gorm.DB, sessionID, prescriptionID string) (bool, error) {
	var count int64
	err := db.Table("session_logs AS s").
		Joins("JOIN prescriptions p ON p.day_id = s.day_id").
		Where("s.id = ? AND p.id = ?", sessionID, prescriptionID).
		Count(&count).Error
	return count > 0, err
}

func abortAccess(c *gin.Context, ok bool, err error) {
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	c.Next()
}
