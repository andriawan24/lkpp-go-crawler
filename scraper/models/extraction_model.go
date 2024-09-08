package models

import (
	"context"
	"encoding/json"

	"github.com/golang-module/carbon"
	"github.com/jackc/pgx/v5"
	"gopkg.in/guregu/null.v4"
)

type Extraction struct {
	Id            string
	UrlFrontierId string
	ArtifactLink  null.String
	RawPageLink   null.String
	Language      string
	Metadata      Metadata
	CreatedAt     carbon.DateTime
	UpdatedAt     carbon.DateTime
}

func UpsertExtraction(ctx context.Context, tx pgx.Tx, extraction Extraction) error {
	sql := "INSERT INTO public.extraction (id, url_frontiers_id, artifact_link, raw_page_link, metadata, language) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (id) DO UPDATE SET artifact_link = EXCLUDED.artifact_link, raw_page_link = EXCLUDED.raw_page_link, metadata = EXCLUDED.metadata, updated_at = EXCLUDED.updated_at"

	metadataJson, err := json.Marshal(extraction.Metadata)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, extraction.Id, extraction.UrlFrontierId, extraction.ArtifactLink, extraction.RawPageLink, metadataJson, extraction.Language)
	if err != nil {
		return err
	}

	return nil
}
