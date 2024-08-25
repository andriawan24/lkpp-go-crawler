package models

import (
	"context"
	"lexicon/lkpp-go-crawler/common"

	"github.com/golang-module/carbon/v2"
	"github.com/jackc/pgx/v5"
)

const (
	URL_STATUS_READY   = 0
	URL_STATUS_CRAWLED = 1
	URL_STATUS_ERROR   = 2
)

type UrlFrontier struct {
	ID        string           `json:"id"`
	Url       string           `json:"url"`
	Domain    string           `json:"domain"`
	Crawler   string           `json:"crawler"`
	Status    int8             `json:"status"`
	CreatedAt carbon.DateTime  `json:"created_at"`
	UpdatedAt *carbon.DateTime `json:"updated_at"`
}

func UpsertUrlFrontier(ctx context.Context, tx pgx.Tx, urlFrontier []UrlFrontier) error {
	sql := "INSERT INTO url_frontiers (id, url, domain, crawler) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET url = EXCLUDED.url, crawler = EXCLUDED.crawler, domain = EXCLUDED.domain, status = EXCLUDED.status, updated_at = EXCLUDED.updated_at"

	batch := &pgx.Batch{}

	for _, url := range urlFrontier {
		batch.Queue(sql, url.ID, url.Url, url.Domain, url.Crawler)
	}

	res := tx.SendBatch(ctx, batch)

	return res.Close()
}

func GetUrlFrontiersUnscraped(ctx context.Context, tx pgx.Tx) ([]UrlFrontier, error) {
	query := "SELECT * FROM url_frontiers WHERE crawler = $1 AND status = $2"

	rows, err := tx.Query(ctx, query, common.CRAWLER_NAME, URL_STATUS_READY)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var urls []UrlFrontier
	for rows.Next() {
		var url UrlFrontier
		err = rows.Scan(&url.ID, &url.Url, &url.Domain, &url.Crawler, &url.Status, &url.CreatedAt, &url.UpdatedAt)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	return urls, nil
}
