package elastic_client

import (
	"crypto/tls"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/elastic/go-elasticsearch/v9"
	"net"
	"net/http"
	"time"
)

type (
	Elastic interface {
		Client() *elasticsearch.Client
	}

	elasticClient struct {
		es *elasticsearch.Client
	}
)

func NewElastic(config *configs.Config) (*elasticClient, error) {
	elasticSearch := config.Config.ElasticSearch
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("%s:%s", elasticSearch.Host, elasticSearch.Port),
		},
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   elasticSearch.Transport.MaxIdleConnsPerHost,
			ResponseHeaderTimeout: time.Second,
			DialContext: (&net.Dialer{
				Timeout:   time.Duration(elasticSearch.Transport.Timeout) * time.Second,
				KeepAlive: time.Duration(elasticSearch.Transport.KeepAlive) * time.Second,
			}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
		MaxRetries:    elasticSearch.MaxRetries,
		EnableMetrics: true,
	})
	if err != nil {
		return nil, err
	}

	return &elasticClient{
		es: es,
	}, nil
}

func (es *elasticClient) Client() *elasticsearch.Client {
	return es.es
}
