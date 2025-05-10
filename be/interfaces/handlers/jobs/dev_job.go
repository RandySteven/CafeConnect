package job_interfaces

import "context"

type DevJob interface {
	CheckHealth(ctx context.Context) error
}
