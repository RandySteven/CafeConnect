package cron_client

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/RandySteven/CafeConnect/be/handlers/jobs"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

type (
	job func(ctx context.Context) error

	Scheduler interface {
		RunAllJobs(ctx context.Context)
	}

	scheduler struct {
		cron *cron.Cron
		jobs *jobs.Jobs
	}
)

func (s *scheduler) RunAllJobs(ctx context.Context) {
	s.cron.Start()
	for _, job := range s.registeredJobs() {
		if err := job(ctx); err != nil {
			log.Fatalln(err)
		}
	}
}

func (s *scheduler) register(registerJob ...job) (jobs []job) {
	for _, job := range registerJob {
		jobs = append(jobs, job)
	}
	return jobs
}

func (s *scheduler) registeredJobs() (jobs []job) {
	jobs = s.register(
		s.jobs.DevJob.CheckHealth,
	)
	return jobs
}

var _ Scheduler = &scheduler{}

func NewScheduler(config *configs.Config) (*scheduler, error) {
	cron2 := config.Config.Cron
	jakartaTime, _ := time.LoadLocation(cron2.Time)
	return &scheduler{
		cron: cron.New(cron.WithSeconds(), cron.WithLocation(jakartaTime)),
	}, nil
}

func (s *scheduler) SetJobs(jobs *jobs.Jobs) {
	s.jobs = jobs
}
