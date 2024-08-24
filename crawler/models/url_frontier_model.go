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
	ID        int
	Url       string
	Crawler   string
	Status    int8
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
}

func UpsertUrlFrontier(ctx context.Context, tx pgx.Tx, urlFrontier []UrlFrontier) error {
	return nil
}

func GetUrlFrontiers() []UrlFrontier {
	return []UrlFrontier{}
}
