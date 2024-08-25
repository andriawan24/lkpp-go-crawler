package main

import (
	"context"
	"fmt"
	"lexicon/lkpp-go-crawler/common"
	"lexicon/lkpp-go-crawler/scraper"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load .env: %v\n", err)
		os.Exit(1)
	}

	context := context.Background()
	dbpool, err := pgxpool.New(context, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create a new pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	err = common.SetDatabase(dbpool)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to set database: %v\n", err)
		os.Exit(1)
	}

	// err = crawler.StartCrawlingUrl()
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Failed to crawl url: %v\n", err)
	// 	os.Exit(1)
	// }

	scraper.StartScraper()
	os.Exit(0)
}
