package scraper

import (
	"fmt"
	"lexicon/lkpp-go-crawler/common"
	"lexicon/lkpp-go-crawler/crawler/services"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
)

func StartScraper() {
	unscraped_urls, err := services.GetUnscrapedUrl()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Unscraped urls", len(unscraped_urls))

	queue, err := queue.New(2, &queue.InMemoryQueueStorage{MaxSize: 10000})
	if err != nil {
		log.Fatalln(err)
	}

	scraper, err := buildScraper(queue)
	if err != nil {
		log.Fatalln(err)
	}

	if len(unscraped_urls) > 0 {
		queue.AddURL(unscraped_urls[0].Url)
	}

	queue.Run(scraper)

	scraper.Wait()
}

func buildScraper(queue *queue.Queue) (*colly.Collector, error) {
	// newExtraction := models.Extraction{}

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

	c.OnHTML("div.large.modal > .content > table.definition", func(h *colly.HTMLElement) {
		retryCount := 0
		maxRetries := 5
		var text string
		for retryCount < maxRetries {
			// Check if the popup is visible
			text = strings.TrimSpace(h.ChildText("#nama-penyedia"))
			fmt.Println("Current value", text)
			if text != "-" {
				fmt.Println("Popup content:", text)
				break
			}
			retryCount++
			fmt.Printf("Retrying... Attempt %d of %d\n", retryCount, maxRetries)
			time.Sleep(5 * time.Second) // Wait for 2 seconds before retrying
		}

		if text == "-" {
			fmt.Println("Popup content not found after retries")
		}
		// text := strings.TrimSpace(h.ChildText("td"))
		// if len(text) > 0 {
		// 	fmt.Println(text)
		// }
	})

	c.OnRequest(func(r *colly.Request) {
		// *newExtraction.RawPageLink = r.URL.String()
		fmt.Println("Visiting", r.URL.String())
		queue.AddRequest(r)
	})

	c.OnScraped(func(r *colly.Response) {
		// frontierId := sha256.Sum256([]byte(r.Request.URL.String()))
		// newExtraction.UrlFrontierId = hex.EncodeToString(frontierId[:])
		// newExtraction.Id = hex.EncodeToString(frontierId[:])
	})

	return c, nil
}
