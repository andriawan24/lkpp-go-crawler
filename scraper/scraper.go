package scraper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	crawlerModel "lexicon/lkpp-go-crawler/crawler/models"
	"lexicon/lkpp-go-crawler/crawler/services"
	"lexicon/lkpp-go-crawler/repository"
	"lexicon/lkpp-go-crawler/scraper/models"
	scraperServices "lexicon/lkpp-go-crawler/scraper/services"
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

type Scraper interface {
	Setup()
	Scrape(ctx context.Context) error
	Teardown()
}

type ScraperImpl struct {
	browser *rod.Browser
}

func (c *ScraperImpl) Setup() {
	c.browser = rod.New().MustConnect()
}

func (c *ScraperImpl) Teardown() {
	c.browser.MustClose()
}

func (c *ScraperImpl) Scrape(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var unscrappedUrlFrontiers []repository.UrlFrontier

	unscrappedUrlFrontiers, err := services.GetUnscrappedUrlFrontiers(ctx, 10)
	if err != nil {
		return err
	}
	log.Info().Msgf("[Started] Unscrapped URLs: %d", len(unscrappedUrlFrontiers))

	var extractions []repository.Extraction

	for _, url := range unscrappedUrlFrontiers {
		extraction := repository.Extraction{}
		metadata := models.Metadata{}

		pageUrl := fmt.Sprintf("https://%s%s", url.Domain, url.Url)

		browser := rod.New().MustConnect()
		page := browser.MustPage(pageUrl)

		log.Debug().Msg("[Started] scraping url " + pageUrl)

		inputs, err := page.Timeout(5 * time.Second).Elements("input")
		if err != nil {
			fmt.Println("[Process] Failed to get inputs", err.Error())
		}

		textareas, err := page.Timeout(5 * time.Second).Elements("textarea")
		if err != nil {
			fmt.Println("[Process] Failed to get inputs", err.Error())
		}

		rule, err := page.Timeout(5 * time.Second).Element("div.select__single-value")
		if err != nil {
			fmt.Println("[Process] Failed to get inputs", err.Error())
		}

		if rule != nil {
			metadata.Injunction.Rule = rule.Timeout(5 * time.Second).MustText()
		}

		for index, textarea := range textareas {
			text, err := textarea.Timeout(5 * time.Second).Text()
			if err != nil {
				fmt.Println("[Process] Failed to get inputs", err.Error())
			}

			if text != "" {
				switch index {
				case 0: // Alamat
					metadata.Address = text
				case 2: // Deskripsi
					metadata.Injunction.Description = text
				}
			}
		}

		for index, input := range inputs {
			text, err := input.Timeout(5 * time.Second).Attribute("value")
			if err != nil {
				fmt.Println("[Process] Failed to get inputs", err.Error())
			}

			if text != nil {
				switch index {
				case 0: // Nama Penyedia
					metadata.Title = *text
				case 2: // Provinsi
					metadata.Province = *text
				case 3: // Kota
					metadata.City = *text
				case 5: // Durasi
					metadata.Injunction.Duration = *text
				case 6: // Tanggal Mulai
					metadata.Injunction.StartDate = *text
				case 7: // Tanggal Selesai
					metadata.Injunction.EndDate = *text
				}
			}
		}

		detailsDiv, err := page.Timeout(5 * time.Second).Elements("dl > div")
		if err != nil {
			fmt.Println("[Process] Failed to get inputs", err.Error())
		}

		var details []models.ProcurementDetail
		detail := models.ProcurementDetail{}

		for index, text := range detailsDiv {
			t, err := text.Timeout(5 * time.Second).Element("dd")
			if err != nil {
				fmt.Println("[Process] Failed to get inputs", err.Error())
			}

			if t != nil {
				switch index {
				case 0: // Tender ID
					detail.TenderID, err = t.Timeout(5 * time.Second).Text()
				case 1: // Nama Paket
					detail.PackageName, err = t.Timeout(5 * time.Second).Text()
				case 2: // Jenis Pengadaan
					detail.ProcurementType, err = t.Timeout(5 * time.Second).Text()
				case 3: // K/L/PD
					detail.InstitutionArea, err = t.Timeout(5 * time.Second).Text()
				case 4: // Satuan Kerja
					detail.Unit, err = t.Timeout(5 * time.Second).Text()
				case 5: // HPS
					detail.EstimatedPrice, err = t.Timeout(5 * time.Second).Text()
				case 6: // Pagu
					detail.Ceiling, err = t.Timeout(5 * time.Second).Text()
				case 7: // Tahun Anggaran
					detail.FiscalYear, err = t.Timeout(5 * time.Second).Text()
				}
			}

			if err != nil {
				fmt.Println("[Process] Failed to get inputs", err.Error())
			}
		}

		// Add details to metadata
		details = append(details, detail)
		metadata.ProcurementDetails = details

		id := sha256.Sum256([]byte(pageUrl))
		extraction.ID = hex.EncodeToString(id[:])
		extraction.Metadata = metadata
		extraction.ArtifactLink = &pageUrl
		extraction.UrlFrontierID = url.ID
		extraction.Language = "ID"

		extractions = append(extractions, extraction)

		browser.MustClose()

		log.Debug().Msg("[Finished] scraped for url " + pageUrl)
	}

	// Add to database
	err = scraperServices.UpsertExtraction(ctx, extractions)
	if err != nil {
		return err
	}

	// Update Frontier Status
	err = services.UpdateFrontierStatuses(
		ctx,
		lo.Map(unscrappedUrlFrontiers, func(urlFrontier repository.UrlFrontier, _ int) lo.Tuple2[string, int16] {
			return lo.Tuple2[string, int16]{A: urlFrontier.ID, B: crawlerModel.URL_FRONTIER_STATUS_CRAWLED}
		}),
	)

	log.Debug().Msg("[Finished] successfully scrape " + strconv.Itoa(len(unscrappedUrlFrontiers)) + " data ")

	return err
}
