package exporter

import (
	"log"

	"github.com/martonorova/prometheus-elasticsearch-index-exporter/config"

	"github.com/elastic/go-elasticsearch/v7"
)

type ElasticsearchClient struct {
	es            *elasticsearch.Client
	elasticErrors chan error
}

func NewElasticsearchClient(cfg config.ElasticsearchConfig) (*ElasticsearchClient, chan error, error) {
	esCfg := elasticsearch.Config{
		Addresses: cfg.Hosts,
	}

	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, nil, err
	}
	esClient := &ElasticsearchClient{
		es:            es,
		elasticErrors: make(chan error),
	}

	err = esClient.testConnection()
	if err != nil {
		return nil, nil, err
	}

	// other parts of the code do not have access to the errorchannel through ElasticsearchClient type
	return esClient, esClient.elasticErrors, nil
}

func (ec *ElasticsearchClient) testConnection() error {
	res, err := ec.es.Info()
	if err != nil {
		return err
	}

	defer res.Body.Close()
	log.Println("Elasticsearch connection OK")

	return nil
}
