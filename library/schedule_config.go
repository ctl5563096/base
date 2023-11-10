package library

import (
	"context"
	"github.com/go-co-op/gocron"
)

type ScheduleConfig struct {
	Schedule   func(s *gocron.Scheduler) *gocron.Scheduler
	HandleFunc func(ctx context.Context) error
}
