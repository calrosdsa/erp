package cron

import (
	"context"
	"github.com/go-co-op/gocron/v2"
)

type scheduler struct {

}

type Scheduler interface {
	RegisterScheduler(ctx context.Context,)
}

func NewSchedule() gocron.Scheduler{
	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
		// handle error
	}

	return s
}

// func (s *scheduler) RegisterJob() {

// }