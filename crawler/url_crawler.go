package crawler

import (
	"fmt"
	"lexicon/lkpp-go-crawler/common"
	"lexicon/lkpp-go-crawler/crawler/models"
	"lexicon/lkpp-go-crawler/crawler/services"
	"strconv"

	"github.com/gocolly/colly/v2"
)

func StartCrawlingUrl() error {
	var urlFrontiers []models.UrlFrontier

	lastPage := getLastPage()

	for currentPage := 1; currentPage <= lastPage; currentPage++ {
		c := colly.NewCollector(
			colly.AllowedDomains(common.CRAWLER_DOMAIN),
		)
		c.OnHTML("table.celled", func(h *colly.HTMLElement) {
			h.ForEach("a.button-detail", func(i int, h *colly.HTMLElement) {
				urlFrontiers = append(urlFrontiers, models.UrlFrontier{
					Url:     fmt.Sprintf("https://%s/daftar-hitam/%s", common.CRAWLER_DOMAIN, h.Attr("data-id")),
					Crawler: common.CRAWLER_NAME,
				})
			})
		})
		c.OnScraped(func(r *colly.Response) {
			fmt.Println("Successfully scraped page", currentPage)
		})
		c.Visit(fmt.Sprintf("https://%s/daftar-hitam?page=%d", common.CRAWLER_DOMAIN, currentPage))
	}

	err := services.UpsetUrl(urlFrontiers)
	if err != nil {
		return err
	}

	return nil
}

func getLastPage() int {
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
		fmt.Println("Finished get last page")
	})
	c.Visit("https://www.inaproc.id/daftar-hitam")
	return lastPage
}
