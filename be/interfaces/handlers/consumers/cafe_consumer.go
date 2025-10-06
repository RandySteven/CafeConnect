package consumer_interfaces

import "context"

type CafeConsumer interface {
	GetCafesByRadius(ctx context.Context) error
}
