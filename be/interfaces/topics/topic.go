package topics_interfaces

import (
	"context"
)

type Topic interface {
	RegisterConsumer(handler func(context.Context, string)) error
	WriteMessage(ctx context.Context, value string) (err error)
	ReadMessage(ctx context.Context) (value string, err error)
}
