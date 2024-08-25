package services

import (
	"context"
	"lexicon/lkpp-go-crawler/common"
	"lexicon/lkpp-go-crawler/scraper/models"
)

func UpsertExtraction(extraction models.Extraction) error {
	context := context.Background()

	tx, err := common.Pool.Begin(context)
	if err != nil {
		return err
	}

	err = models.UpsertExtraction(context, tx, extraction)
	if err != nil {
		return err
	}

	return nil
}
