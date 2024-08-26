package crawler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"lexicon/lkpp-go-crawler/common"
	"lexicon/lkpp-go-crawler/crawler/models"
	"lexicon/lkpp-go-crawler/crawler/services"
	"strconv"

	"github.com/gocolly/colly/v2"
)

func StartCrawlingUrl() error {
	var urlFrontiers []models.UrlFrontier

	lastPage := getLastPage("/non-aktif")
	for currentPage := 1; currentPage <= lastPage; currentPage++ {
		c := colly.NewCollector(
			colly.AllowedDomains(common.CRAWLER_DOMAIN),
		)

		c.OnHTML("table.celled", func(h *colly.HTMLElement) {
			h.ForEach("a.button-detail", func(i int, h *colly.HTMLElement) {
				url := fmt.Sprintf("https://%s/daftar-hitam/non-aktif?page=%d#%s", common.CRAWLER_DOMAIN, currentPage, h.Attr("data-id"))
				id := sha256.Sum256([]byte(url))

				urlFrontiers = append(urlFrontiers, models.UrlFrontier{
					ID:      hex.EncodeToString(id[:]),
					Url:     url,
					Crawler: common.CRAWLER_NAME,
					Domain:  common.CRAWLER_DOMAIN,
				})
			})
		})

		c.OnScraped(func(r *colly.Response) {
			fmt.Println("[finished] Finished scrape", r.Request.URL.String())
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("[started]: Visiting URL", r.URL.String())
		})

		c.Visit(fmt.Sprintf("https://%s/daftar-hitam/non-aktif?page=%d", common.CRAWLER_DOMAIN, currentPage))
		fmt.Println("[finished]: Successfully crawled endpoint non-aktif", fmt.Sprintf("https://%s/daftar-hitam/non-aktif?page=%d", common.CRAWLER_DOMAIN, currentPage))
	}

	var endpoints = []string{"", "/penundaan", "/batal"}

	for _, endpoint := range endpoints {
		lastPage := getLastPage(endpoint)

		for currentPage := 1; currentPage <= lastPage; currentPage++ {
			c := colly.NewCollector(
				colly.AllowedDomains(common.CRAWLER_DOMAIN),
			)

			c.OnHTML("table.celled", func(h *colly.HTMLElement) {
				h.ForEach("a.button-detail", func(i int, h *colly.HTMLElement) {
					url := fmt.Sprintf("https://%s/daftar-hitam%s/%s", common.CRAWLER_DOMAIN, endpoint, h.Attr("data-id"))
					id := sha256.Sum256([]byte(url))

					urlFrontiers = append(urlFrontiers, models.UrlFrontier{
						ID:      hex.EncodeToString(id[:]),
						Url:     url,
						Crawler: common.CRAWLER_NAME,
						Domain:  common.CRAWLER_DOMAIN,
					})
				})
			})

			c.OnRequest(func(r *colly.Request) {
				fmt.Println("[started]: Visiting URL", r.URL.String())
			})

			c.Visit(fmt.Sprintf("https://%s/daftar-hitam%s?page=%d", common.CRAWLER_DOMAIN, endpoint, currentPage))
		}

		fmt.Println("[finished]: Successfully crawled endpoint", endpoint)
	}

	err := services.UpsetUrl(urlFrontiers)
	if err != nil {
		return err
	}

	return nil
}

func getLastPage(endpoint string) int {
	c := colly.NewCollector(
		colly.AllowedDomains(common.CRAWLER_DOMAIN),
	)

	lastPage := 1
	c.OnHTML(".pagination", func(h *colly.HTMLElement) {
		var err error
		childTexts := h.ChildTexts("a.item")
		if len(childTexts) > 0 {
			lastPage, err = strconv.Atoi(childTexts[len(childTexts)-1])
			if err != nil {
				fmt.Println("Error")
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("[started]: Visiting", r.URL.String())
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("[finished]: Get last page", lastPage)
	})

	c.Visit(fmt.Sprintf("https://%s/daftar-hitam%s", common.CRAWLER_DOMAIN, endpoint))

	return lastPage
}
