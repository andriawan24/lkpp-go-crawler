package models

import (
	"context"

	"github.com/golang-module/carbon/v2"
	"github.com/jackc/pgx/v5"
)

const (
	URL_STATUS_READY   = 0
	URL_STATUS_CRAWLED = 1
	URL_STATUS_ERROR   = 2
)

type UrlFrontier struct {
	ID        int             `json:"id"`
	Url       string          `json:"url"`
	Crawler   string          `json:"crawler"`
	Status    int8            `json:"status"`
	CreatedAt carbon.DateTime `json:"created_at"`
	UpdatedAt carbon.DateTime `json:"updated_at"`
}

func UpsertUrlFrontier(ctx context.Context, tx pgx.Tx, urlFrontier []UrlFrontier) error {
	sql := "INSERT INTO url_frontiers (url, crawler) VALUES ($1, $2) ON CONFLICT (url) DO UPDATE SET url = EXCLUDED.url, crawler = EXCLUDED.crawler, status = EXCLUDED.status, updated_at = EXCLUDED.updated_at"

	batch := &pgx.Batch{}

	for _, url := range urlFrontier {
		batch.Queue(sql, url.Url, url.Crawler)
	}

	res := tx.SendBatch(ctx, batch)

	return res.Close()
}

func GetUrlFrontiers() []UrlFrontier {
	return []UrlFrontier{}
}
