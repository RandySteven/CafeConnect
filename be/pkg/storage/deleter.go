package storage_client

import (
	"context"
	"fmt"
)

func (g *googleStorage) DeleteFile(ctx context.Context, objectFileName string) (err error) {
	if err = g.bkt.Object(objectFileName).Delete(ctx); err != nil {
		return fmt.Errorf(`error while delete %v`, err)
	}
	return nil
}
