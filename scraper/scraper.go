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
	"github.com/go-rod/rod"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"gopkg.in/guregu/null.v4"
)

func StartScraper() error {
	unscrapedUrls, err := services.GetUnscrapedUrl()
	if err != nil {
		return err
	}

	var nonActiveUrls []crawler_model.UrlFrontier
	var activeUrls []crawler_model.UrlFrontier
	for _, url := range unscrapedUrls {
		if strings.Contains(url.Url, "non-aktif") {
			nonActiveUrls = append(nonActiveUrls, url)
		} else {
			activeUrls = append(activeUrls, url)
		}
	}

	fmt.Println("[started]: Unscraped non-active urls", len(activeUrls))
	fmt.Println("[started]: Unscraped active urls", len(nonActiveUrls))

	// not active
	for _, url := range nonActiveUrls {
		url := url.Url
		browser := rod.New().ControlURL("ws://host.docker.internal:3000")
		err = browser.Connect()
		if err != nil {
			return err
		}

		page := browser.MustPage(url).MustWaitStable()
		mainTable := page.MustElement(".ui.modal.large").MustElement("table.definition > tbody")
		foulTable := page.MustElement(".ui.modal.large").MustElements("table#injunctions > tbody > tr")

		injunctions := []models.Injunction{}
		for _, el := range foulTable {
			validityPeriodValue := el.MustElement("td:nth-child(3) table tr td:nth-child(2), td:nth-child(3) table tr td:nth-child(3)")
			validityPeriodText := validityPeriodValue.MustText()
			validities := strings.Split(validityPeriodText, "s/d")

			cols := el.MustElements("td")
			injunction := models.Injunction{
				Number:      strings.TrimSpace(cols.First().MustText()),
				Rule:        strings.TrimSpace(cols[1].MustElement(".header").MustText()),
				Description: strings.TrimSpace(cols[1].MustElement(".description").MustText()),
				StartDate:   strings.TrimSpace(validities[0]),
				EndDate:     strings.TrimSpace(validities[len(validities)-1]),
			}
			injunctions = append(injunctions, injunction)
		}

		title := strings.TrimSpace(mainTable.MustElement("#nama-penyedia").MustText())
		id := sha256.Sum256([]byte(title))
		metadata := models.Metadata{
			ID:          hex.EncodeToString(id[:]),
			Title:       title,
			NPWP:        strings.TrimSpace(mainTable.MustElement("#npwp").MustText()),
			Address:     strings.TrimSpace(mainTable.MustElement("#alamat").MustText()),
			Province:    strings.TrimSpace(mainTable.MustElement("#propinsi").MustText()),
			City:        strings.TrimSpace(mainTable.MustElement("#kabupaten-kota").MustText()),
			Status:      "inactive",
			Injunctions: injunctions,
		}

		frontierId := sha256.Sum256([]byte(url))
		newExtraction := models.Extraction{
			Id:            hex.EncodeToString(frontierId[:]),
			UrlFrontierId: hex.EncodeToString(frontierId[:]),
			RawPageLink:   null.StringFrom(url),
			Metadata:      metadata,
		}

		fmt.Println("[started]: Upserting extraction", newExtraction.RawPageLink)
		err := scraper_service.UpsertExtraction(newExtraction)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("[finished] Upserting extraction", newExtraction.RawPageLink)

		fmt.Println("[started]: Upserting crawler", newExtraction.RawPageLink)
		err = services.UpdateUrlFrontierStatus(newExtraction.Id, crawler_model.URL_STATUS_CRAWLED)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("[finished] Upserting crawler", newExtraction.RawPageLink)
	}

	// Active
	queue, err := queue.New(2, &queue.InMemoryQueueStorage{MaxSize: 10000})
	if err != nil {
		return err
	}

	scraper, err := buildScraper(queue)
	if err != nil {
		return err
	}

	for _, url := range activeUrls {
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

	c.OnHTML("table.table-list > tbody", func(h *colly.HTMLElement) {
		injunction := models.Injunction{}
		h.ForEach("tr", func(i int, h *colly.HTMLElement) {
			h.DOM.Find("td:nth-child(1)").Contents().Each(func(i int, s *goquery.Selection) {
				if i == 0 {
					number := strings.TrimSpace(s.Text())
					injunction.Number = number
				} else if i == 1 {
					rule := strings.TrimSpace(s.Find(".header").Text())
					description := strings.TrimSpace(s.Find(".description").Text())
					injunction.Rule = rule
					injunction.Description = description
				}
			})

			h.DOM.Find("td:nth-child(2)").Contents().Each(func(i int, s *goquery.Selection) {
				if i == 0 {
					startDate := strings.TrimSpace(s.Text())
					injunction.StartDate = startDate
				} else if i == 2 {
					endDate := strings.TrimSpace(s.Text())
					injunction.EndDate = endDate
				}
			})

			injunction.PublishedDate = h.ChildText("tr:nth-child(1) > td:nth-child(3)")
		})
		injunctions = append(injunctions, injunction)
	})

	c.OnRequest(func(r *colly.Request) {
		currentUrl = r.URL.String()
		fmt.Println("[visiting]:", currentUrl)
		queue.AddRequest(r)
	})

	c.OnScraped(func(r *colly.Response) {
		if strings.Contains(currentUrl, "penundaan") {
			newExtraction.Metadata.Status = "pending"
		} else if strings.Contains(currentUrl, "batal") {
			newExtraction.Metadata.Status = "cancelled"
		} else {
			newExtraction.Metadata.Status = "active"
		}

		rawPageUrl := r.Request.URL.String()
		frontierId := sha256.Sum256([]byte(rawPageUrl))

		newExtraction.RawPageLink = null.StringFrom(rawPageUrl)
		newExtraction.Id = hex.EncodeToString(frontierId[:])
		newExtraction.UrlFrontierId = hex.EncodeToString(frontierId[:])
		newExtraction.Metadata.Injunctions = injunctions

		fmt.Println("[started]: Upserting extraction", newExtraction.RawPageLink)
		err := scraper_service.UpsertExtraction(newExtraction)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("[finished] Upserting extraction", newExtraction.RawPageLink)

		fmt.Println("[started]: Upserting crawler", newExtraction.RawPageLink)
		err = services.UpdateUrlFrontierStatus(newExtraction.Id, crawler_model.URL_STATUS_CRAWLED)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("[finished] Upserting crawler", newExtraction.RawPageLink)
	})

	return c, nil
}
