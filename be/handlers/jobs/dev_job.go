package jobs

import (
	"context"
	job_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/jobs"
	"log"
)

type DevJob struct {
}

func (d *DevJob) CheckHealth(ctx context.Context) error {
	log.Println("Scheduler work properly")
	return nil
}

var _ job_interfaces.DevJob = &DevJob{}

func newDevJob() *DevJob {
	return &DevJob{}
}
