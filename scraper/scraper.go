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
	"sync"
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

	pagePool := rod.NewPagePool(20)
	defer pagePool.Cleanup(func(p *rod.Page) {
		err := p.Close()
		if err != nil {
			log.Error().Err(err).Msg("Error closing page")
		}
	})

	create := func() (*rod.Page, error) {
		incognito, err := c.browser.Incognito()
		if err != nil {
			log.Error().Err(err).Msg("Error creating incognito page")
			return nil, err
		}
		return incognito.MustPage(), nil
	}

	job := func(ctx context.Context, urlFrontier repository.UrlFrontier, c chan<- repository.Extraction) error {
		page, err := pagePool.Get(create)
		if err != nil {
			log.Error().Err(err).Msg("Error getting page")
			return err
		}
		defer pagePool.Put(page)
		extraction, err := scrapeUrlFrontiers(ctx, page, urlFrontier)
		if err != nil {
			log.Error().Err(err).Msg("Error scraping url frontier")
			return err
		}
		c <- extraction
		return nil
	}

	for ok := true; ok; ok = len(unscrappedUrlFrontiers) > 0 {

		// Check if context is cancelled before starting new chunk
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var scraperErrors []error
		unscrappedUrlFrontiers, err := services.GetUnscrappedUrlFrontiers(ctx, 3500)
		if err != nil {
			log.Error().Err(err).Msg("Error fetching unscrapped url frontier")
		}
		log.Info().Msgf("Unscrapped URLs: %d", len(unscrappedUrlFrontiers))

		chunks := lo.Chunk(unscrappedUrlFrontiers, 20)

		for _, urlFrontiers := range chunks {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			wg := sync.WaitGroup{}
			extractionChan := make(chan repository.Extraction, len(urlFrontiers))
			errChan := make(chan error, len(urlFrontiers))

			for _, urlFrontier := range urlFrontiers {
				wg.Add(1)
				go func(urlFrontier repository.UrlFrontier) {
					defer wg.Done()

					select {
					case <-ctx.Done():
						errChan <- ctx.Err()
					default:
						if err := job(ctx, urlFrontier, extractionChan); err != nil {
							errChan <- fmt.Errorf("error crawling %s: %w", urlFrontier.Url, err)
						}
					}
				}(urlFrontier)
			}

			// Wait for all goroutines to finish
			wg.Wait()

			// Close channels after all goroutines are done
			close(errChan)
			close(extractionChan)

			// Collect errors and extractions
			for err := range errChan {
				scraperErrors = append(scraperErrors, err)
			}

			var extractions []repository.Extraction
			for extraction := range extractionChan {
				extractions = append(extractions, extraction)
			}

			if len(scraperErrors) > 0 {
				log.Error().Err(scraperErrors[0]).Msg("Error scraping url frontier")
			}
			log.Info().Msgf("Upserting extractions")
			err := scraperServices.UpsertExtraction(ctx, extractions)
			if err != nil {
				log.Error().Err(err).Msg("Error upserting extractions")
			}
			log.Info().Msgf("Updating url frontier statuses")
			err = services.UpdateFrontierStatuses(ctx, lo.Map(urlFrontiers, func(urlFrontier repository.UrlFrontier, _ int) lo.Tuple2[string, int16] {
				return lo.Tuple2[string, int16]{A: urlFrontier.ID, B: crawlerModel.URL_FRONTIER_STATUS_CRAWLED}
			}))
			if err != nil {
				log.Error().Err(err).Msg("Error updating url frontier status")
			}
			log.Info().Msgf("Finished scraping chunk")

		}

	}
	log.Info().Msgf("Finished scraping all url frontiers")
	return nil
}

func scrapeUrlFrontiers(ctx context.Context, page *rod.Page, urlFrontier repository.UrlFrontier) (repository.Extraction, error) {
	select {
	case <-ctx.Done():
		return repository.Extraction{}, ctx.Err()
	default:
	}
	log.Info().Msgf("Scraping url: %s", urlFrontier.Url)

	extraction := repository.Extraction{}
	metadata := models.Metadata{}

	pageUrl := fmt.Sprintf("https://%s%s", urlFrontier.Domain, urlFrontier.Url)

	rpCtx := page.MustNavigate(pageUrl)

	inputs, err := rpCtx.Timeout(5 * time.Second).Elements("input")
	if err != nil {
		log.Error().Msg("Failed to get procuremenet details " + err.Error())
	}

	textareas, err := rpCtx.Timeout(5 * time.Second).Elements("textarea")
	if err != nil {
		log.Error().Msg("Failed to get procuremenet details " + err.Error())
	}

	rule, err := rpCtx.Timeout(5 * time.Second).Element("div.select__single-value")
	if err != nil {
		log.Error().Msg("Failed to get procuremenet details " + err.Error())
	}

	if rule != nil {
		metadata.Injunction.Rule = rule.Timeout(5 * time.Second).MustText()
	}

	for index, textarea := range textareas {
		text, err := textarea.Timeout(5 * time.Second).Text()
		if err != nil {
			log.Error().Msg("Failed to get textareas " + err.Error())
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
			log.Error().Msg("Failed to get inputs " + err.Error())
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

	detailsDiv, err := rpCtx.Timeout(5 * time.Second).Elements("dl > div")
	if err != nil {
		log.Error().Msg("Failed to get inputs " + err.Error())
	}

	var details []models.ProcurementDetail
	detail := models.ProcurementDetail{}

	for index, text := range detailsDiv {
		t, err := text.Timeout(5 * time.Second).Element("dd")
		if err != nil {
			log.Error().Msg("Failed to get procuremenet details " + err.Error())
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
			log.Error().Msg("Failed to get procuremenet details " + err.Error())
		}
	}

	// Add details to metadata
	details = append(details, detail)
	metadata.ProcurementDetails = details

	id := sha256.Sum256([]byte(pageUrl))
	extraction.ID = hex.EncodeToString(id[:])
	extraction.Metadata = metadata
	extraction.ArtifactLink = &pageUrl
	extraction.UrlFrontierID = urlFrontier.ID
	extraction.Language = "ID"
	siteContent, err := rpCtx.HTML()
	if err != nil {
		log.Error().Msg("Failed to get site contents " + err.Error())
	}

	extraction.SiteContent = &siteContent

	log.Debug().Msg("Finished scraped for url " + pageUrl)

	return extraction, nil
}
