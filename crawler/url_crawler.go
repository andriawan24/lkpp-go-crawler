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

	var endpoints = []string{"", "/non-aktif", "/penundaan", "/batal"}

	for _, endpoint := range endpoints {
		lastPage := getLastPage(endpoint)

		for currentPage := 1; currentPage <= lastPage; currentPage++ {
			c := colly.NewCollector(
				colly.AllowedDomains(common.CRAWLER_DOMAIN),
			)

			c.OnHTML("table.celled", func(h *colly.HTMLElement) {
				h.ForEach("a.button-detail", func(i int, h *colly.HTMLElement) {
					url := fmt.Sprintf("https://%s/daftar-hitam%s#%s", common.CRAWLER_DOMAIN, endpoint, h.Attr("data-id"))
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
				fmt.Println("Successfully scrape page", currentPage, "for endpoint", endpoint)
			})

			c.Visit(fmt.Sprintf("https://%s/daftar-hitam%s?page=%d", common.CRAWLER_DOMAIN, endpoint, currentPage))
		}
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

	lastPage := 0
	c.OnHTML(".pagination", func(h *colly.HTMLElement) {
		var err error
		childTexts := h.ChildTexts("a.item")
		lastPage, err = strconv.Atoi(childTexts[len(childTexts)-1])
		if err != nil {
			fmt.Println("Error")
		}
	})
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished get last page", lastPage)
	})
	c.Visit(fmt.Sprintf("https://%s/daftar-hitam%s", common.CRAWLER_DOMAIN, endpoint))
	return lastPage
}
