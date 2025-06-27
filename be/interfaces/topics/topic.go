package topics_interfaces

import (
	"context"
)

type Topic interface {
	RegisterConsumer(handler func(string)) error
	WriteMessage(ctx context.Context, value string) (err error)
}
