package scraper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"lexicon/lkpp-go-crawler/common"
	"lexicon/lkpp-go-crawler/crawler/services"
	"lexicon/lkpp-go-crawler/scraper/models"
	"os"
	"strings"
	"time"

	crawler_model "lexicon/lkpp-go-crawler/crawler/models"
	scraper_service "lexicon/lkpp-go-crawler/scraper/services"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"gopkg.in/guregu/null.v4"
)

func StartScraper() error {
	unscraped_urls, err := services.GetUnscrapedUrl()
	if err != nil {
		return err
	}

	fmt.Println("[started]: Unscraped urls", len(unscraped_urls))

	queue, err := queue.New(2, &queue.InMemoryQueueStorage{MaxSize: 10000})
	if err != nil {
		return err
	}

	scraper, err := buildScraper(queue)
	if err != nil {
		return err
	}

	// for _, url := range unscraped_urls {
	queue.AddURL(unscraped_urls[0].Url)
	// }

	queue.Run(scraper)

	scraper.Wait()

	return nil
}

func buildScraper(queue *queue.Queue) (*colly.Collector, error) {
	newExtraction := models.Extraction{
		Metadata: models.Metadata{},
		Language: "id",
	}

	c := colly.NewCollector(
		colly.AllowedDomains(common.CRAWLER_DOMAIN),
		colly.MaxDepth(1),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  common.CRAWLER_DOMAIN,
		Parallelism: 10,
		Delay:       time.Second * 2,
		RandomDelay: time.Second * 2,
	})

	c.SetRequestTimeout(time.Minute * 2)

	c.OnHTML("table.definition > tbody > tr", func(h *colly.HTMLElement) {
		text := strings.TrimSpace(h.ChildText("td:nth-child(2)"))
		switch h.Index {
		case 0:
			id := sha256.Sum256([]byte(text))
			newExtraction.Id = hex.EncodeToString(id[:])
			newExtraction.Metadata.ID = hex.EncodeToString(id[:])
			newExtraction.Metadata.Title = text
		case 1:
			newExtraction.Metadata.NPWP = text
		case 2:
			newExtraction.Metadata.Address = text
		case 3:
			newExtraction.Metadata.City = text
		case 4:
			newExtraction.Metadata.Province = text
		}
	})

	c.OnHTML("table.table-list > tbody > tr:nth-child(1) > td:nth-child(1)", func(h *colly.HTMLElement) {
		h.DOM.Contents().Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				cleanedText := strings.TrimSpace(s.Text())
				newExtraction.Metadata.Number = cleanedText
			} else if i == 1 {
				rule := strings.TrimSpace(s.Find(".header").Text())
				description := strings.TrimSpace(s.Find(".description").Text())
				newExtraction.Metadata.Rule = rule
				newExtraction.Metadata.Description = description
			}
		})
	})

	c.OnHTML("table.table-list > tbody > tr:nth-child(1) > td:nth-child(2)", func(h *colly.HTMLElement) {
		h.DOM.Contents().Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				startDate := strings.TrimSpace(s.Text())
				newExtraction.Metadata.StartDate = startDate
			} else if i == 2 {
				endDate := strings.TrimSpace(s.Text())
				newExtraction.Metadata.EndDate = endDate
			}
		})
	})

	c.OnHTML("table.table-list > tbody > tr:nth-child(1) > td:nth-child(3)", func(h *colly.HTMLElement) {
		newExtraction.Metadata.PublishedDate = h.Text
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("[visiting]:", r.URL.String())
		queue.AddRequest(r)
	})

	c.OnScraped(func(r *colly.Response) {
		rawPageUrl := r.Request.URL.String()
		frontierId := sha256.Sum256([]byte(rawPageUrl))

		newExtraction.RawPageLink = null.StringFrom(rawPageUrl)
		newExtraction.Id = hex.EncodeToString(frontierId[:])
		newExtraction.UrlFrontierId = hex.EncodeToString(frontierId[:])

		fmt.Println("[started]: Upserting extraction", newExtraction)
		err := scraper_service.UpsertExtraction(newExtraction)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("[finished] Upserting extraction", newExtraction.Id)

		fmt.Println("[started]: Upserting crawler", newExtraction)
		err = services.UpdateUrlFrontierStatus(newExtraction.Id, crawler_model.URL_STATUS_CRAWLED)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("[finished] Upserting crawler", newExtraction.Id)
	})

	return c, nil
}
