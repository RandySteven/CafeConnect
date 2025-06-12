package elastic_client

import (
	"bytes"
	"context"
	json2 "encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v9"
	"log"
	"strings"
)

func Search[T any](ctx context.Context, c *elasticsearch.Client, index string, key, value string) (result []*T, err error) {
	query := fmt.Sprintf(
		`
			{
				"query": {
					"match_all": {
						"%s": "%s"
					}
				}
			}
		`, key, value,
	)
	response, err := c.Search(
		c.Search.WithContext(ctx),
		c.Search.WithIndex(index),
		c.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, err
	}
	log.Println(response)
	return
}

func CreateIndexing[T any](ctx context.Context, c *elasticsearch.Client, index string, data []*T) (err error) {
	json, _ := json2.Marshal(data)
	_, err = c.Index(
		index,
		bytes.NewReader(json),
		c.Index.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	return nil
}

func GetIndex[T any](ctx context.Context, c *elasticsearch.Client, index string) (err error) {
	return
}
