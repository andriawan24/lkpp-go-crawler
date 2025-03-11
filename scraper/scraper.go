package scraper

import (
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
)

type TableData struct {
	NamaPenyedia string
	NPWP         string
}

func StartScraper() error {
	// unscraped_urls, err := services.GetUnscrapedUrl()
	// if err != nil {
	// 	return err
	// }

	// fmt.Println("[started]: Unscraped urls", len(unscraped_urls))

	// url := unscraped_urls[0].Url

	// browser := rod.New().MustConnect()
	// page := browser.MustPage(url)

	// fmt.Println("[started]: Scraping url", url)

	// // Extract data from the page
	// var data TableData

	// // // Find the div containing "Nama Penyedia" label
	// namaPenyediaSection := page.MustElementR("div", "Nama Penyedia")
	// // // Navigate to the closest input element in the adjacent column
	// if namaPenyediaInput := namaPenyediaSection.MustParent().MustParent().MustNext().MustElement("input"); namaPenyediaInput != nil {
	// 	data.NamaPenyedia = namaPenyediaInput.MustProperty("value").String()
	// }

	// // Find the div containing "NPWP Penyedia" label
	// npwpSection := page.MustElementR("div", "NPWP Penyedia")
	// // Navigate to the closest input element in the adjacent column
	// if npwpInput := npwpSection.MustParent().MustParent().MustNext().MustElement("input"); npwpInput != nil {
	// 	data.NPWP = npwpInput.MustProperty("value").String()
	// }

	// fmt.Println("[Finished]: Extracted data:")
	// fmt.Println("Nama Penyedia:", data.NamaPenyedia)
	// fmt.Println("NPWP:", data.NPWP)

	// // Update the extraction model or database if needed
	// // For example:
	// // newExtraction := models.Extraction{
	// //     Metadata: models.Metadata{
	// //         Title: data.NamaPenyedia,
	// //         NPWP: data.NPWP,
	// //     },
	// //     Language: "id",
	// // }
	// // scraper_service.UpsertExtraction(newExtraction)

	// browser.MustClose()
	return nil
}

func buildScraper(queue *queue.Queue) (*colly.Collector, error) {
	// var currentUrl string

	// newExtraction := models.Extraction{
	// 	Metadata: models.Metadata{},
	// 	Language: "id",
	// }

	// c := colly.NewCollector(
	// 	colly.AllowedDomains(common.CRAWLER_DOMAIN),
	// 	colly.MaxDepth(1),
	// )

	// c.Limit(&colly.LimitRule{
	// 	DomainGlob:  common.CRAWLER_DOMAIN,
	// 	Parallelism: 10,
	// 	Delay:       time.Second * 2,
	// 	RandomDelay: time.Second * 2,
	// })

	// c.SetRequestTimeout(time.Minute * 2)

	// c.OnHTML("table.definition > tbody > tr", func(h *colly.HTMLElement) {
	// 	if !strings.Contains(currentUrl, "non-aktif") {
	// 		text := strings.TrimSpace(h.ChildText("td:nth-child(2)"))
	// 		switch h.Index {
	// 		case 0:
	// 			id := sha256.Sum256([]byte(text))
	// 			newExtraction.Id = hex.EncodeToString(id[:])
	// 			newExtraction.Metadata.ID = hex.EncodeToString(id[:])
	// 			newExtraction.Metadata.Title = text
	// 		case 1:
	// 			newExtraction.Metadata.NPWP = text
	// 		case 2:
	// 			newExtraction.Metadata.Address = text
	// 		case 3:
	// 			newExtraction.Metadata.City = text
	// 		case 4:
	// 			newExtraction.Metadata.Province = text
	// 		}
	// 	}
	// })

	// c.OnHTML("table.table-list > tbody > tr:nth-child(1) > td:nth-child(1)", func(h *colly.HTMLElement) {
	// 	if !strings.Contains(currentUrl, "non-aktif") {
	// 		h.DOM.Contents().Each(func(i int, s *goquery.Selection) {
	// 			if i == 0 {
	// 				cleanedText := strings.TrimSpace(s.Text())
	// 				newExtraction.Metadata.Number = cleanedText
	// 			} else if i == 1 {
	// 				rule := strings.TrimSpace(s.Find(".header").Text())
	// 				description := strings.TrimSpace(s.Find(".description").Text())
	// 				newExtraction.Metadata.Rule = rule
	// 				newExtraction.Metadata.Description = description
	// 			}
	// 		})
	// 	}
	// })

	// c.OnHTML("table.table-list > tbody > tr:nth-child(1) > td:nth-child(2)", func(h *colly.HTMLElement) {
	// 	if !strings.Contains(currentUrl, "non-aktif") {
	// 		h.DOM.Contents().Each(func(i int, s *goquery.Selection) {
	// 			if i == 0 {
	// 				startDate := strings.TrimSpace(s.Text())
	// 				newExtraction.Metadata.StartDate = startDate
	// 			} else if i == 2 {
	// 				endDate := strings.TrimSpace(s.Text())
	// 				newExtraction.Metadata.EndDate = endDate
	// 			}
	// 		})

	// 		if strings.Contains(currentUrl, "penundaan") {
	// 			newExtraction.Metadata.Status = "pending"
	// 		} else if strings.Contains(currentUrl, "batal") {
	// 			newExtraction.Metadata.Status = "cancelled"
	// 		} else {
	// 			newExtraction.Metadata.Status = "active"
	// 		}
	// 	}
	// })

	// c.OnHTML("table.table-list > tbody > tr:nth-child(1) > td:nth-child(3)", func(h *colly.HTMLElement) {
	// 	if !strings.Contains(currentUrl, "non-aktif") {
	// 		newExtraction.Metadata.PublishedDate = h.Text
	// 	}
	// })

	// c.OnHTML("table.ui.table.small.celled.very.padded tbody tr", func(e *colly.HTMLElement) {
	// 	if strings.Contains(currentUrl, "non-aktif") {
	// 		id := strings.Split(currentUrl, "#")
	// 		if strings.Contains(e.ChildAttr("a.button-detail", "data-id"), id[len(id)-1]) {
	// 			title := strings.TrimSpace(e.ChildText("td:nth-child(1) h5 a"))
	// 			number := strings.TrimSpace(e.ChildText("td:nth-child(1) div.npwp strong"))
	// 			city := strings.TrimSpace(e.ChildText("td:nth-child(2) .ui.list .item .content .header"))
	// 			address := strings.TrimSpace(e.ChildText("td:nth-child(2) .ui.list .item .content .description"))
	// 			endDate := strings.TrimSpace(e.ChildText("td:nth-child(3) table tbody tr:nth-child(2) td:nth-child(2)"))
	// 			startDate := strings.Split(strings.TrimSpace(e.ChildText("td:nth-child(3) table tbody tr:nth-child(3) td:nth-child(2)")), "s/d")

	// 			id := sha256.Sum256([]byte(title))
	// 			newExtraction.Metadata.ID = hex.EncodeToString(id[:])
	// 			newExtraction.Metadata.Title = title
	// 			newExtraction.Metadata.Number = number
	// 			newExtraction.Metadata.Address = address
	// 			newExtraction.Metadata.City = city
	// 			newExtraction.Metadata.EndDate = endDate
	// 			newExtraction.Metadata.StartDate = strings.TrimRight(startDate[0], "\n\r")
	// 			newExtraction.Metadata.EndDate = strings.Trim(startDate[1], "\n\r")
	// 			newExtraction.Metadata.Status = "not-active"
	// 		}
	// 	}
	// })

	// c.OnRequest(func(r *colly.Request) {
	// 	currentUrl = r.URL.String()
	// 	fmt.Println("[visiting]:", currentUrl)
	// 	queue.AddRequest(r)
	// })

	// c.OnScraped(func(r *colly.Response) {
	// 	rawPageUrl := r.Request.URL.String()
	// 	frontierId := sha256.Sum256([]byte(rawPageUrl))

	// 	newExtraction.RawPageLink = null.StringFrom(rawPageUrl)
	// 	newExtraction.Id = hex.EncodeToString(frontierId[:])
	// 	newExtraction.UrlFrontierId = hex.EncodeToString(frontierId[:])

	// 	fmt.Println("[started]: Upserting extraction", newExtraction.RawPageLink)
	// 	err := scraper_service.UpsertExtraction(newExtraction)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		os.Exit(1)
	// 	}
	// 	fmt.Println("[finished] Upserting extraction", newExtraction.RawPageLink)

	// 	fmt.Println("[started]: Upserting crawler", newExtraction.RawPageLink)
	// 	// err = services.UpdateUrlFrontierStatus(newExtraction.Id, crawler_model.URL_FRONTIER_STATUS_CRAWLED)
	// 	// if err != nil {
	// 	// 	fmt.Println(err)
	// 	// 	os.Exit(1)
	// 	// }
	// 	fmt.Println("[finished] Upserting crawler", newExtraction.RawPageLink)
	// })

	return nil, nil
}
