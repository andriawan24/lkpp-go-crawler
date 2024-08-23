package main

func main() {

	// var products []models.Product

	// c := colly.NewCollector(
	// 	colly.AllowedDomains("www.scrapingcourse.com"),
	// )

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Requesting to visit", r.URL)
	// })

	// c.OnError(func(r *colly.Response, err error) {
	// 	fmt.Println("Error", err)
	// })

	// c.OnHTML("li.product", func(e *colly.HTMLElement) {
	// 	product := models.Product{}

	// 	product.Url = e.ChildAttr("a", "href")
	// 	product.Name = e.ChildText(".product-name")
	// 	product.Price = e.ChildText(".price")
	// 	product.Image = e.ChildAttr("img", "src")

	// 	products = append(products, product)
	// })

	// c.OnScraped(func(r *colly.Response) {
	// 	fmt.Println(r.Request.URL, "scraped!")

	// 	file, err := os.Create("products.csv")
	// 	if err != nil {
	// 		log.Fatalln("Failed to create output CSV", err)
	// 	}

	// 	defer file.Close()

	// 	writer := csv.NewWriter(file)

	// 	headers := []string{
	// 		"Url",
	// 		"Image",
	// 		"Name",
	// 		"Price",
	// 	}
	// 	writer.Write(headers)

	// 	for _, product := range products {
	// 		record := []string{
	// 			product.Url,
	// 			product.Image,
	// 			product.Name,
	// 			product.Price,
	// 		}
	// 		writer.Write(record)
	// 	}

	// 	defer writer.Flush()
	// })

	// c.Visit("https://www.scrapingcourse.com/ecommerce")
}
