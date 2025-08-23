package topics_interfaces

import (
	"context"
)

type Topic interface {
	WriteMessage(ctx context.Context, value string) (err error)
	ReadMessage(ctx context.Context) (value string, err error)
}
