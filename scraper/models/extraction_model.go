package models

import (
	"context"
	"encoding/json"

	"github.com/golang-module/carbon"
	"github.com/jackc/pgx/v5"
)

type Extraction struct {
	Id            string
	UrlFrontierId string
	SiteContent   *string
	ArtifactLink  *string
	RawPageLink   *string
	Language      string
	Metadata      Metadata
	CreatedAt     carbon.DateTime
	UpdatedAt     carbon.DateTime
}

func UpsertExtraction(ctx context.Context, tx pgx.Tx, extraction Extraction) error {
	sql := "INSERT INTO extractions (id, url_frontiers_id, site_content, artifact_link, raw_page_link, metadata, language, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (id) DO UPDATE SET site_content = EXCLUDED.site_content, artifact_link = EXCLUDED.artifact_link, raw_page_link = EXCLUDED.raw_page_link, metadata = EXCLUDED.metadata, updated_at = EXCLUDED.updated_at"
	metadataJson, err := json.Marshal(extraction.Metadata)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, extraction.Id, extraction.UrlFrontierId, extraction.SiteContent, extraction.ArtifactLink, extraction.RawPageLink, metadataJson, extraction.Language)
	if err != nil {
		return err
	}

	return nil
}
