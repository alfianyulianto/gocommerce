package elasticsearch

import (
	"context"
	"embed"
	_ "embed"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

//go:embed mappings/*
var mappings embed.FS

func (c *Client) ensureIndices(ctx context.Context) error {
	dir, err := mappings.ReadDir("mappings")
	if err != nil {
		c.log.WithError(err).Error("Failed to read mappings")
		return err
	}

	for _, entry := range dir {
		if !entry.IsDir() {
			indexName := fmt.Sprintf("%s_%s", strings.ToLower(c.cfg.App.Name), strings.TrimSuffix(entry.Name(), "_mapping.json"))

			content, err := mappings.ReadFile("mappings/" + entry.Name())
			if err != nil {
				c.log.WithError(err).Error("Failed to read mappings")
				return err
			}

			exists, _ := c.indexExists(ctx, indexName)
			if exists {
				continue
			}

			err = c.createIndex(ctx, indexName, string(content))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Client) indexExists(ctx context.Context, index string) (bool, error) {
	indexes := []string{index}
	res, err := c.es.Indices.Exists(
		indexes,
		c.es.Indices.Exists.WithContext(ctx),
	)
	if err != nil {
		c.log.WithField("indexes", indexes).WithError(err).Errorf("Failed to check if index exists")
		return false, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return false, fmt.Errorf("Failed to check if indexes %s exists with status %s", indexes, res.Status())
	}

	c.log.WithFields(logrus.Fields{
		"indexes":  indexes,
		"status":   res.Status(),
		"response": res.String(),
	}).Info("Index exists successfully")
	return true, nil
}

func (c *Client) createIndex(ctx context.Context, index, mapping string) error {
	res, err := c.es.Indices.Create(
		index,
		c.es.Indices.Create.WithContext(ctx),
		c.es.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		c.log.WithField("index", index).WithError(err).Errorf("Failed to create index")
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		c.log.WithFields(logrus.Fields{
			"index":    index,
			"status":   res.Status(),
			"response": res.String(),
		}).WithError(err).Errorf("Failed to create index")
		return fmt.Errorf("Failed to create index %s with status %s", index, res.Status())
	}

	c.log.WithFields(logrus.Fields{
		"index":    index,
		"status":   res.Status(),
		"response": res.String(),
	}).Info("Created index successfully")
	return nil
}
