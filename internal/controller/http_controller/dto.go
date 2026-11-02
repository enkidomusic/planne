package httpcontroller

import (
	"time"

	"github.com/enkidomusic/planne/internal/domain"
)

type Task struct {
	ID                       int         `json:"id"`
	Description              string      `json:"description"`
	Date                     *time.Time  `json:"date"`
	IsDone                   bool        `json:"is_done"`
	ExpectedDurationMinutes  *int        `json:"expected_duration_minutes"`
	RepeatCount              *int        `json:"repeat_count"`
	RepeatInterval           *int        `json:"repeat_interval"`
	DoneDates                []time.Time `json:"done_dates"`
	ExcludedWeekdays         []int       `json:"excluded_weekdays"`
	EisenhowerMatrixCategory *int        `json:"eisenhower_matrix_category"`
}

func (t *Task) ToDomain() *domain.Task {
	var expectedDuration *time.Duration = nil
	if t.ExpectedDurationMinutes != nil {
		duration := time.Duration(*t.ExpectedDurationMinutes) * time.Minute
		expectedDuration = &duration
	}

	return &domain.Task{
		ID:                       t.ID,
		Description:              t.Description,
		Date:                     t.Date,
		IsDone:                   t.IsDone,
		ExpectedDuration:         expectedDuration,
		RepeatCount:              t.RepeatCount,
		RepeatInterval:           t.RepeatInterval,
		DoneDates:                t.DoneDates,
		ExcludedWeekdays:         t.ExcludedWeekdays,
		EisenhowerMatrixCategory: t.EisenhowerMatrixCategory,
	}
}

func DTOFromDomain(task *domain.Task) *Task {
	var expectedDurationMinutes *int = nil
	if task.ExpectedDuration != nil {
		duration := int(task.ExpectedDuration.Minutes())
		expectedDurationMinutes = &duration
	}

	return &Task{
		ID:                       task.ID,
		Description:              task.Description,
		Date:                     task.Date,
		IsDone:                   task.IsDone,
		ExpectedDurationMinutes:  expectedDurationMinutes,
		RepeatCount:              task.RepeatCount,
		RepeatInterval:           task.RepeatInterval,
		DoneDates:                task.DoneDates,
		ExcludedWeekdays:         task.ExcludedWeekdays,
		EisenhowerMatrixCategory: task.EisenhowerMatrixCategory,
	}
}
