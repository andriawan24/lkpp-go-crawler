package main

import (
	"context"
	"fmt"
	"lexicon/lkpp-go-crawler/common"
	"lexicon/lkpp-go-crawler/crawler"
	"lexicon/lkpp-go-crawler/scraper"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
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

	rootCommand := &cobra.Command{
		Use:   "lexicon-lkpp-crawler",
		Short: "Crawl LKPP Blacklist of Indonesia",
		Long:  "Crawl LKPP Blacklist of Indonesia",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("[Started]: Start URL Crawler")
			crawler.StartCrawlingUrl()
			fmt.Println("[Started]: Start Web Scraper")
			scraper.StartScraper()
			fmt.Println("[Finished]: Finished Crawling URL and Web Scraping, happy coding!")
		},
	}

	rootCommand.AddCommand(crawlerCommand())
	rootCommand.AddCommand(scraperCommand())

	err = crawlerCommand().Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func crawlerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "crawler",
		Short: "URL Crawler for detail page of LKPP Blacklist of Indonesia website",
		Long:  "URL Crawler for detail page of LKPP Blacklist of Indonesia website",
		Run: func(cmd *cobra.Command, args []string) {
			err := crawler.StartCrawlingUrl()
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			os.Exit(0)
		},
	}
}

func scraperCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "scraper",
		Short: "URL Scraper for detail page of LKPP Blacklist of Indonesia website",
		Long:  "URL Scraper for detail page of LKPP Blacklist of Indonesia website",
		Run: func(cmd *cobra.Command, args []string) {
			err := scraper.StartScraper()
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			os.Exit(0)
		},
	}
}
