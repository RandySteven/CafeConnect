package jobs

import job_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/jobs"

type Jobs struct {
	DevJob job_interfaces.DevJob
}

func NewJobs() *Jobs {
	return &Jobs{
		DevJob: newDevJob(),
	}
}
