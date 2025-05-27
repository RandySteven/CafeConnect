package consumer_interfaces

import (
	"context"
)

type DummyConsumer interface {
	CheckHealth(ctx context.Context)
}
