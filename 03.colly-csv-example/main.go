package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly"
)

const (
	fileName = "coinmarketcap.csv"
	target   = "https://coinmarketcap.com/all/views/all/"
)

func crawlAndWrite(writer *csv.Writer) error {
	// Write CSV header
	writer.Write([]string{"Name", "Symbol", "Price (USD)", "Volume (USD)", "Market capacity (USD)", "Change (1h)", "Change (24h)", "Change (7d)"})

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnHTML(".cmc-table__table-wrapper-outer tbody .cmc-table-row", func(e *colly.HTMLElement) {
		writer.Write([]string{
			e.ChildText(".cmc-table__column-name--name"),
			e.ChildText(".cmc-table__cell--sort-by__symbol"),
			e.ChildText(".cmc-table__cell--sort-by__price"),
			e.ChildText(".cmc-table__cell--sort-by__volume-24-h"),
			e.ChildText(".cmc-table__cell--sort-by__market-cap"),
			e.ChildText(".cmc-table__cell--sort-by__percent-change-1-h"),
			e.ChildText(".cmc-table__cell--sort-by__percent-change-24-h"),
			e.ChildText(".cmc-table__cell--sort-by__percent-change-7-d"),
		})
	})

	return c.Visit(target)
}

func main() {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fileName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = crawlAndWrite(writer)
	if err != nil {
		log.Fatalf("Cannot visit url %q: %s\n", target, err)
		return
	}
	log.Printf("Scraping finished, check file %q for results\n", fileName)
}
