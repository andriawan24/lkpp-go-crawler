package crawler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"lexicon/lkpp-go-crawler/common"
	"lexicon/lkpp-go-crawler/crawler/models"
	"lexicon/lkpp-go-crawler/crawler/services"
	"lexicon/lkpp-go-crawler/repository"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

func StartCrawlingUrl() error {
	var urlFrontiers []repository.UrlFrontier

	lastPage := getLastPage("")
	for currentPage := 1; currentPage <= lastPage; currentPage++ {
		c := colly.NewCollector(colly.AllowedDomains(common.CRAWLER_DOMAIN))

		c.OnHTML("table.border-tertiary50 > tbody", func(h *colly.HTMLElement) {
			h.ForEach("tr", func(i int, h *colly.HTMLElement) {
				title := h.ChildText("td:nth-child(2) span")
				scenario := h.ChildText("td:nth-child(3) span")
				packageNumber := h.ChildText("td:nth-child(4) span")
				packageName := h.ChildText("td:nth-child(5) span")
				startDate := h.ChildText("td:nth-child(6) span")
				endDate := h.ChildText("td:nth-child(7) span")
				duration := h.ChildText("td:nth-child(8) span")
				status := h.ChildText("td:nth-child(9) span")

				url := h.ChildAttr("td:nth-child(10) a", "href")
				id := sha256.Sum256([]byte(url))

				urlFrontiers = append(urlFrontiers, repository.UrlFrontier{
					ID:      hex.EncodeToString(id[:]),
					Url:     url,
					Crawler: common.CRAWLER_NAME,
					Domain:  common.CRAWLER_DOMAIN,
					Status:  int16(models.URL_FRONTIER_STATUS_NEW),
					Metadata: models.UrlFrontierMetadata{
						Title:         title,
						Scenario:      scenario,
						PackageNumber: packageNumber,
						PackageName:   packageName,
						StartDate:     startDate,
						EndDate:       endDate,
						Duration:      duration,
						Status:        status,
					},
				})
			})
		})

		c.OnScraped(func(r *colly.Response) {
			fmt.Println("[finished] Finished scrap", r.Request.URL.String())
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("[started]: Visiting URL", r.URL.String())
		})

		c.Visit(fmt.Sprintf("https://%s/?p=%d", common.CRAWLER_DOMAIN, currentPage))

		fmt.Println("[finished]: Successfully crawled endpoint ", fmt.Sprintf("https://%s/?p=%d", common.CRAWLER_DOMAIN, currentPage))
	}

	err := services.UpsertUrl(context.Background(), urlFrontiers)
	if err != nil {
		return err
	}

	return nil
}

func getLastPage(endpoint string) int {
	c := colly.NewCollector(
		colly.AllowedDomains(common.CRAWLER_DOMAIN),
	)

	lastPage := 0

	c.OnHTML("button.Pagination_button__m7YBp", func(h *colly.HTMLElement) {

		text := strings.TrimSpace(h.Text)

		num, err := strconv.Atoi(text)
		if err == nil {
			if num > lastPage {
				lastPage = num
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Error:", err)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("[started]: Visiting", r.URL.String())
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("[finished]: Get last page", lastPage)
	})

	c.Visit(fmt.Sprintf("https://%s/%s", common.CRAWLER_DOMAIN, endpoint))

	return lastPage
}
