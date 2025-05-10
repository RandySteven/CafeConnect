package cron_client

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/robfig/cron/v3"
	"time"
)

type (
	job func(ctx context.Context) error

	Scheduler interface {
		RunAllJobs(ctx context.Context)
	}

	scheduler struct {
		cron *cron.Cron
	}
)

func NewScheduler(config *configs.Config) (*scheduler, error) {
	cron2 := config.Config.Cron
	jakartaTime, _ := time.LoadLocation(cron2.Time)
	return &scheduler{
		cron: cron.New(cron.WithSeconds(), cron.WithLocation(jakartaTime)),
	}, nil
}
