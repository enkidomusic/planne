package domain

import "time"

type Task struct {
	ID                       int
	Description              string
	Date                     *time.Time
	IsDone                   bool
	ExpectedDuration         *time.Duration
	RepeatCount              *int
	RepeatInterval           *int
	DoneDates                []time.Time
	ExcludedWeekdays         []int
	EisenhowerMatrixCategory *int
}
