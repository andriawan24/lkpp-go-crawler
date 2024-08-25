package services

import (
	"context"
	"lexicon/lkpp-go-crawler/common"
	"lexicon/lkpp-go-crawler/crawler/models"
)

func UpsetUrl(urlFrontiers []models.UrlFrontier) error {
	ctx := context.Background()

	tx, err := common.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	err = models.UpsertUrlFrontier(ctx, tx, urlFrontiers)
	if err != nil {
		return err
	}

	tx.Commit(ctx)

	return nil
}

func GetUnscrapedUrl() ([]models.UrlFrontier, error) {
	ctx := context.Background()

	tx, err := common.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	urls, err := models.GetUrlFrontiersUnscraped(ctx, tx)
	if err != nil {
		return nil, err
	}

	return urls, nil
}
