package topics_interfaces

import (
	"context"
)

type Topic interface {
	ReadMessage(ctx context.Context, key string) (result string, err error)
	WriteMessage(ctx context.Context, key string, value string) (err error)
}
