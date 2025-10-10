package service

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/vicepalma/roma-system/backend/internal/domain"
	"github.com/vicepalma/roma-system/backend/internal/repository"
)

// ---- PIVOT (para charts) ----

type PivotResponse struct {
	Columns []string                         `json:"columns"`
	Rows    []map[string]interface{}         `json:"rows"`
	Catalog []repository.ExerciseCatalogItem `json:"catalog,omitempty"`
	Mode    string                           `json:"mode"`
	Days    int                              `json:"days"`
}

type orderedSet struct {
	m map[string]struct{}
	a []string
}

func (s *orderedSet) add(v string) {
	if v == "" {
		return
	}
	if s.m == nil {
		s.m = map[string]struct{}{}
	}
	if _, ok := s.m[v]; ok {
		return
	}
	s.m[v] = struct{}{}
	s.a = append(s.a, v)
}
func (s *orderedSet) values() []string { return s.a }

type HistoryService interface {
	GetHistory(ctx context.Context, discipleID string, days int) (*HistoryResponse, error)
	GetPRs(ctx context.Context, discipleID string) ([]repository.PRRow, error)

	GetDailyByExercise(ctx context.Context, discipleID string, days int, includeCatalog bool, tz string) ([]repository.DailyExerciseVolume, []repository.ExerciseCatalogItem, error)
	GetDailyByMuscle(ctx context.Context, discipleID string, days int, tz string) ([]repository.DailyMuscleVolume, error)

	GetPivotByExercise(ctx context.Context, discipleID string, days int, includeCatalog bool, metric string, tz string) (*PivotResponse, error)
	GetPivotByMuscle(ctx context.Context, discipleID string, days int, metric string, tz string) (*PivotResponse, error)
}

type HistoryResponse struct {
	Sessions []domain.SessionLog `json:"sessions"`
}

type historyService struct{ repo repository.HistoryRepository }

func NewHistoryService(r repository.HistoryRepository) HistoryService {
	return &historyService{repo: r}
}

func clampDays(days int) int {
	if days <= 0 || days > 180 {
		return 14
	}
	return days
}

func sinceFromDays(days int) time.Time {
	d := clampDays(days)
	return time.Now().UTC().AddDate(0, 0, -d)
}

func (s *historyService) GetHistory(ctx context.Context, discipleID string, days int) (*HistoryResponse, error) {
	since := sinceFromDays(days)
	sessions, err := s.repo.ListRecentSessions(ctx, discipleID, since)
	if err != nil {
		return nil, err
	}
	return &HistoryResponse{Sessions: sessions}, nil
}

func (s *historyService) GetPRs(ctx context.Context, discipleID string) ([]repository.PRRow, error) {
	return s.repo.BestSetsByExercise(ctx, discipleID)
}

// summaries crudos (para pivot y summary)
func (s *historyService) GetDailyByExercise(ctx context.Context, discipleID string, days int, includeCatalog bool, tz string) ([]repository.DailyExerciseVolume, []repository.ExerciseCatalogItem, error) {
	loc := normTZ(tz)
	sinceDate := sinceLocalDate(days, loc)

	rows, err := s.repo.DailyVolumeByExercise(ctx, discipleID, sinceDate, loc.String())
	if err != nil {
		return nil, nil, err
	}

	var catalog []repository.ExerciseCatalogItem
	if includeCatalog {
		catalog, err = s.repo.ListRelevantExercisesForUser(ctx, discipleID)
		if err != nil {
			return nil, nil, err
		}
	}

	out := fillExerciseGapsWithCatalog(rows, catalog, mustParseDate(sinceDate, loc), time.Now().In(loc))
	return out, catalog, nil
}

func (s *historyService) GetDailyByMuscle(ctx context.Context, discipleID string, days int, tz string) ([]repository.DailyMuscleVolume, error) {
	loc := normTZ(tz)
	sinceDate := sinceLocalDate(days, loc)

	rows, err := s.repo.DailyVolumeByMuscle(ctx, discipleID, sinceDate, loc.String())
	if err != nil {
		return nil, err
	}
	return fillMuscleGaps(rows, mustParseDate(sinceDate, loc), time.Now().In(loc)), nil
}

/* ----------------- gap fillers ----------------- */

func dateKey(t time.Time) string { return t.Format("2006-01-02") }

// gap filler que usa catálogo para “fijar” todas las claves (aunque no aparezcan en datos)
func fillExerciseGapsWithCatalog(in []repository.DailyExerciseVolume, catalog []repository.ExerciseCatalogItem, start, end time.Time) []repository.DailyExerciseVolume {
	// claves desde datos
	type key struct{ id, name string }
	keys := map[key]struct{}{}
	for _, r := range in {
		keys[key{r.ExerciseID, r.ExerciseName}] = struct{}{}
	}
	// + claves desde catálogo
	for _, c := range catalog {
		keys[key{c.ID, c.Name}] = struct{}{}
	}

	// index existente
	existing := make(map[string]map[key]repository.DailyExerciseVolume)
	for _, r := range in {
		if _, ok := existing[r.Date]; !ok {
			existing[r.Date] = make(map[key]repository.DailyExerciseVolume)
		}
		existing[r.Date][key{r.ExerciseID, r.ExerciseName}] = r
	}

	out := make([]repository.DailyExerciseVolume, 0, len(in))
	for d := start.Truncate(24 * time.Hour); !d.After(end); d = d.Add(24 * time.Hour) {
		day := dateKey(d)
		for k := range keys {
			if byDay, ok := existing[day]; ok {
				if row, ok2 := byDay[k]; ok2 {
					out = append(out, row)
					continue
				}
			}
			// gap → 0s
			out = append(out, repository.DailyExerciseVolume{
				Date:         day,
				ExerciseID:   k.id,
				ExerciseName: k.name,
				Volume:       0,
				Sets:         0,
				Reps:         0,
			})
		}
	}

	// orden: fecha asc, nombre asc
	sort.Slice(out, func(i, j int) bool {
		if out[i].Date == out[j].Date {
			return out[i].ExerciseName < out[j].ExerciseName
		}
		return out[i].Date < out[j].Date
	})
	return out
}

func fillMuscleGaps(in []repository.DailyMuscleVolume, start, end time.Time) []repository.DailyMuscleVolume {
	// claves presentes
	muscles := map[string]struct{}{}
	for _, r := range in {
		muscles[r.PrimaryMuscle] = struct{}{}
	}
	// index existente
	existing := make(map[string]map[string]repository.DailyMuscleVolume)
	for _, r := range in {
		if _, ok := existing[r.Date]; !ok {
			existing[r.Date] = make(map[string]repository.DailyMuscleVolume)
		}
		existing[r.Date][r.PrimaryMuscle] = r
	}

	out := make([]repository.DailyMuscleVolume, 0, len(in))
	for d := start.Truncate(24 * time.Hour); !d.After(end); d = d.Add(24 * time.Hour) {
		day := dateKey(d)
		for m := range muscles {
			if byDay, ok := existing[day]; ok {
				if row, ok2 := byDay[m]; ok2 {
					out = append(out, row)
					continue
				}
			}
			out = append(out, repository.DailyMuscleVolume{
				Date:          day,
				PrimaryMuscle: m,
				Volume:        0,
				Sets:          0,
				Reps:          0,
			})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Date == out[j].Date {
			return out[i].PrimaryMuscle < out[j].PrimaryMuscle
		}
		return out[i].Date < out[j].Date
	})
	return out
}

func (s *historyService) GetPivotByExercise(ctx context.Context, discipleID string, days int, includeCatalog bool, metric string, tz string) (*PivotResponse, error) {
	metric = normMetric(metric)
	loc := normTZ(tz)

	rows, catalog, err := s.GetDailyByExercise(ctx, discipleID, days, includeCatalog, tz)
	if err != nil {
		return nil, err
	}

	seriesOrder := orderedSet{}
	for _, r := range rows {
		seriesOrder.add(r.ExerciseName)
	}
	cols := append([]string{"date"}, seriesOrder.values()...)

	// index por fecha->serie
	byDate := map[string]map[string]float64{}
	for _, r := range rows {
		if _, ok := byDate[r.Date]; !ok {
			byDate[r.Date] = map[string]float64{}
		}
		byDate[r.Date][r.ExerciseName] = valueByMetricExercise(r, metric)
	}

	// FECHAS explícitas en tz (incluye HOY)
	dates := dateRangeLocal(days, loc)

	outRows := make([]map[string]interface{}, 0, len(dates))
	for _, d := range dates {
		row := map[string]interface{}{"date": d}
		for _, sname := range seriesOrder.values() {
			if v, ok := byDate[d][sname]; ok {
				row[sname] = v
			} else {
				row[sname] = 0.0
			}
		}
		outRows = append(outRows, row)
	}

	resp := &PivotResponse{
		Columns: cols,
		Rows:    outRows,
		Mode:    "by_exercise",
		Days:    clampDays(days),
	}
	if includeCatalog {
		resp.Catalog = catalog
	}
	return resp, nil
}

func (s *historyService) GetPivotByMuscle(ctx context.Context, discipleID string, days int, metric string, tz string) (*PivotResponse, error) {
	metric = normMetric(metric)
	loc := normTZ(tz)

	rows, err := s.GetDailyByMuscle(ctx, discipleID, days, tz)
	if err != nil {
		return nil, err
	}

	seriesOrder := orderedSet{}
	for _, r := range rows {
		seriesOrder.add(r.PrimaryMuscle)
	}
	cols := append([]string{"date"}, seriesOrder.values()...)

	byDate := map[string]map[string]float64{}
	for _, r := range rows {
		if _, ok := byDate[r.Date]; !ok {
			byDate[r.Date] = map[string]float64{}
		}
		byDate[r.Date][r.PrimaryMuscle] = valueByMetricMuscle(r, metric)
	}

	dates := dateRangeLocal(days, loc) // 👈 fechas en tz, incluye HOY

	outRows := make([]map[string]interface{}, 0, len(dates))
	for _, d := range dates {
		row := map[string]interface{}{"date": d}
		for _, m := range seriesOrder.values() {
			if v, ok := byDate[d][m]; ok {
				row[m] = v
			} else {
				row[m] = 0.0
			}
		}
		outRows = append(outRows, row)
	}

	return &PivotResponse{
		Columns: cols,
		Rows:    outRows,
		Mode:    "by_muscle",
		Days:    clampDays(days),
	}, nil
}

func normMetric(metric string) string {
	switch strings.ToLower(metric) {
	case "sets":
		return "sets"
	case "reps":
		return "reps"
	default:
		return "volume" // por defecto
	}
}

func valueByMetricExercise(r repository.DailyExerciseVolume, metric string) float64 {
	switch metric {
	case "sets":
		return float64(r.Sets)
	case "reps":
		return float64(r.Reps)
	default:
		return r.Volume
	}
}

func valueByMetricMuscle(r repository.DailyMuscleVolume, metric string) float64 {
	switch metric {
	case "sets":
		return float64(r.Sets)
	case "reps":
		return float64(r.Reps)
	default:
		return r.Volume
	}
}

func normTZ(tz string) *time.Location {
	if tz == "" {
		tz = "America/Santiago"
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		// fallback seguro
		loc = time.FixedZone("UTC", 0)
	}
	return loc
}

func sinceLocalDate(days int, loc *time.Location) string {
	d := clampDays(days)
	nowLocal := time.Now().In(loc)
	// N días incluyendo HOY => restar (N-1)
	start := time.Date(nowLocal.Year(), nowLocal.Month(), nowLocal.Day(), 0, 0, 0, 0, loc).
		AddDate(0, 0, -(d - 1))
	return start.Format("2006-01-02")
}

// helpers para convertir "YYYY-MM-DD" a time.Time en la tz local
func mustParseDate(yyyyMMdd string, loc *time.Location) time.Time {
	t, err := time.ParseInLocation("2006-01-02", yyyyMMdd, loc)
	if err != nil {
		return time.Now().In(loc)
	}
	return t
}

func dateRangeLocal(days int, loc *time.Location) []string {
	d := clampDays(days)
	nowLocal := time.Now().In(loc)
	start := time.Date(nowLocal.Year(), nowLocal.Month(), nowLocal.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, -(d - 1))
	out := make([]string, 0, d)
	for t := start; !t.After(nowLocal); t = t.Add(24 * time.Hour) {
		out = append(out, t.Format("2006-01-02"))
	}
	return out
}
