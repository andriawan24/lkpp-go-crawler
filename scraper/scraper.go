package scraper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"lexicon/lkpp-go-crawler/common"
	"lexicon/lkpp-go-crawler/crawler/services"
	"lexicon/lkpp-go-crawler/scraper/models"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"gopkg.in/guregu/null.v4"
)

func StartScraper() error {
	unscrapedUrls, err := services.GetUnscrapedUrl()
	if err != nil {
		return err
	}

	fmt.Println("[started]: Unscraped urls", len(unscrapedUrls))

	queue, err := queue.New(2, &queue.InMemoryQueueStorage{MaxSize: 10000})
	if err != nil {
		return err
	}

	scraper, err := buildScraper(queue)
	if err != nil {
		return err
	}

	for _, url := range unscrapedUrls {
		queue.AddURL(url.Url)
	}

	queue.Run(scraper)
	scraper.Wait()

	return nil
}

func buildScraper(queue *queue.Queue) (*colly.Collector, error) {
	var currentUrl string
	injunctions := []models.Injunction{}

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

	c.OnHTML("input", func(h *colly.HTMLElement) {
		fmt.Println("Data ke ", currentUrl, h.Attr("value"))
	})

	c.OnHTML("textarea", func(h *colly.HTMLElement) {
		fmt.Println("Data ke ", currentUrl, h.Text)
	})

	// c.OnHTML("table.table-list > tbody", func(h *colly.HTMLElement) {
	// 	injunction := models.Injunction{}
	// 	h.ForEach("tr", func(i int, h *colly.HTMLElement) {
	// 		h.DOM.Find("td:nth-child(1)").Contents().Each(func(i int, s *goquery.Selection) {
	// 			if i == 0 {
	// 				number := strings.TrimSpace(s.Text())
	// 				injunction.Number = number
	// 			} else if i == 1 {
	// 				rule := strings.TrimSpace(s.Find(".header").Text())
	// 				description := strings.TrimSpace(s.Find(".description").Text())
	// 				injunction.Rule = rule
	// 				injunction.Description = description
	// 			}
	// 		})

	// 		h.DOM.Find("td:nth-child(2)").Contents().Each(func(i int, s *goquery.Selection) {
	// 			if i == 0 {
	// 				startDate := strings.TrimSpace(s.Text())
	// 				injunction.StartDate = startDate
	// 			} else if i == 2 {
	// 				endDate := strings.TrimSpace(s.Text())
	// 				injunction.EndDate = endDate
	// 			}
	// 		})

	// 		injunction.PublishedDate = h.ChildText("tr:nth-child(1) > td:nth-child(3)")
	// 	})
	// 	injunctions = append(injunctions, injunction)
	// })

	c.OnRequest(func(r *colly.Request) {
		currentUrl = r.URL.String()
		fmt.Println("[visiting]:", currentUrl)
		queue.AddRequest(r)
	})

	c.OnScraped(func(r *colly.Response) {
		rawPageUrl := r.Request.URL.String()
		frontierId := sha256.Sum256([]byte(rawPageUrl))

		newExtraction.RawPageLink = null.StringFrom(rawPageUrl)
		newExtraction.Id = hex.EncodeToString(frontierId[:])
		newExtraction.UrlFrontierId = hex.EncodeToString(frontierId[:])
		newExtraction.Metadata.Injunctions = injunctions

		// fmt.Println("[started]: Upserting extraction", newExtraction.RawPageLink)
		// err := scraper_service.UpsertExtraction(newExtraction)
		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }
		// fmt.Println("[finished] Upserting extraction", newExtraction.RawPageLink)

		// fmt.Println("[started]: Upserting crawler", newExtraction.RawPageLink)
		// err = services.UpdateUrlFrontierStatus(newExtraction.Id, crawler_model.URL_STATUS_CRAWLED)
		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }
		// fmt.Println("[finished] Upserting crawler", newExtraction.RawPageLink)
	})

	return c, nil
}
