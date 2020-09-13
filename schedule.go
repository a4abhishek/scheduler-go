package scheduler

import "time"

type Schedule struct {
	ID  int
	sig <-chan time.Time
	timer *time.Timer
}

type ScheduleBuilder struct {
	schedule *Schedule
}

func NewScheduleBuilder(id int) *ScheduleBuilder {
	return &ScheduleBuilder{
		schedule: &Schedule{
			ID:  id,
			sig: nil,
		},
	}
}

func (s *ScheduleBuilder) WithDuration(duration time.Duration) *ScheduleBuilder {
	timer := time.NewTimer(duration)

	s.schedule.sig = timer.C
	s.schedule.timer = timer

	return s
}

func (s *ScheduleBuilder) WithTimestamp(timestamp time.Time) *ScheduleBuilder {
	return s.WithDuration(time.Since(timestamp))
}

func (s *ScheduleBuilder) Build() *Schedule {
	return s.schedule
}
