package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alfianyulianto/gocommerce/config"
	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/sirupsen/logrus"
)

type Client struct {
	es  *elasticsearch.Client
	cfg *config.Config
	log *logrus.Logger
}

func NewClient(cfg *config.Config, log *logrus.Logger) *Client {
	es, err := elasticsearch.New(
		elasticsearch.WithAddresses(cfg.Elasticsearch.Addresses...),
		elasticsearch.WithLogger(&elastictransport.JSONLogger{
			Output:             log.Out,
			EnableRequestBody:  true,
			EnableResponseBody: false,
		}),
	)
	if err != nil {
		log.WithField("addresses", cfg.Elasticsearch.Addresses).WithError(err).Error("Failed to create elasticsearch client")
	}

	res, err := es.Info()
	if err != nil {
		log.WithField("addresses", cfg.Elasticsearch.Addresses).WithError(err).Fatal("Failed to get elasticsearch info")
	}
	defer res.Body.Close()

	log.WithField("addresses", cfg.Elasticsearch.Addresses).Info("Successfully connected to elasticsearch")

	client := &Client{es: es, cfg: cfg, log: log}
	_ = client.ensureIndices(context.Background())

	return client
}

func (c *Client) Index(ctx context.Context, index string, id string, document interface{}) error {
	bytes, err := json.Marshal(document)
	if err != nil {
		c.log.WithField("document", document).WithError(err).Error("Failed to marshal document for elasticsearch")
		return err
	}

	res, err := c.es.Index(
		index,
		strings.NewReader(string(bytes)),
		c.es.Index.WithContext(ctx),
		c.es.Index.WithDocumentID(id),
		c.es.Index.WithRefresh("true"),
	)
	if err != nil {
		c.log.WithFields(logrus.Fields{
			"document": string(bytes),
			"index":    index,
			"_id":      id,
		}).WithError(err).Errorf("Failed to insert document")
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		c.log.WithFields(logrus.Fields{
			"document": string(bytes),
			"index":    index,
			"_id":      id,
			"status":   res.Status(),
			"response": res.String(),
		}).WithError(err).Error("Failed to insert document")
		return fmt.Errorf("Failed to insert document %s to index %s with status %s", string(bytes), index, res.Status())
	}

	c.log.WithFields(logrus.Fields{
		"document": string(bytes),
		"index":    index,
		"_id":      id,
		"status":   res.Status(),
		"response": res.String(),
	}).Info("Successfully inserted document")
	return err
}

func (c *Client) Delete(ctx context.Context, index string, id string) error {
	res, err := c.es.Delete(index, id, c.es.Delete.WithContext(ctx))
	if err != nil {
		c.log.WithFields(logrus.Fields{
			"index": index,
			"_id":   id,
		}).WithError(err).Errorf("Failed to delete document")
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		c.log.WithFields(logrus.Fields{
			"index":    index,
			"_id":      id,
			"status":   res.Status(),
			"response": res.String(),
		}).WithError(err).Errorf("Failed to delete document")
		return fmt.Errorf("Failed to delete document _id %s from index %s with status %s", id, index, res.Status())
	}

	c.log.WithFields(logrus.Fields{
		"index":    index,
		"_id":      id,
		"status":   res.Status(),
		"response": res.String(),
	}).Info("Successfully deleted document")
	return nil
}

func (c *Client) Search(ctx context.Context, index string, query map[string]interface{}) ([]json.RawMessage, int64, error) {
	bytes, err := json.Marshal(query)
	if err != nil {
		c.log.WithField("query", query).WithError(err).Errorf("Failed to marshal query for elasticsearch")
		return nil, 0, err
	}

	res, err := c.es.Search(
		c.es.Search.WithIndex(index),
		c.es.Search.WithContext(ctx),
		c.es.Search.WithBody(strings.NewReader(string(bytes))),
		c.es.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		c.log.WithFields(logrus.Fields{
			"query": query,
			"index": index,
		}).WithError(err).Errorf("Failed to search document")
	}
	defer res.Body.Close()

	if res.IsError() {
		c.log.WithFields(logrus.Fields{
			"query":    query,
			"index":    index,
			"status":   res.Status(),
			"response": res.String(),
		}).WithError(err).Errorf("Failed to search document")
		return nil, 0, fmt.Errorf("Failed to search document from index %s with status %s", index, res.Status())
	}

	var result struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		c.log.WithError(err).Errorf("Failed to decode search result")
		panic(err)
	}

	var docs []json.RawMessage
	for _, h := range result.Hits.Hits {
		docs = append(docs, h.Source)
	}

	return docs, result.Hits.Total.Value, nil
}

func (c *Client) UserIndex() string {
	return c.cfg.Elasticsearch.UserIndex
}

func (c *Client) ProductIndex() string {
	return c.cfg.Elasticsearch.ProductIndex
}

func (c *Client) OrderIndex() string {
	return c.cfg.Elasticsearch.OrderIndex
}
