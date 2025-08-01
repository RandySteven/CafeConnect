package jobs

import job_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/jobs"

type Jobs struct {
	DevJob         job_interfaces.DevJob
	TransactionJob job_interfaces.TransactionJob
}

func NewJobs() *Jobs {
	return &Jobs{
		DevJob: newDevJob(),
	}
}
