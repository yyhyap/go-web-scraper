package main

import (
	"encoding/csv"
	"fmt"
	"go-web-scraper/logger"
	"os"
	"path/filepath"

	"github.com/gocolly/colly"
)

var (
	PokemonProducts []PokemonProduct
)

type PokemonProduct struct {
	Url   string
	Image string
	Name  string
	Price string
}

func main() {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		logger.Logger.Info(fmt.Sprintf("Visiting: %v", r.URL))
	})

	c.OnError(func(_ *colly.Response, err error) {
		logger.Logger.Error(fmt.Sprintf("error when crawling: %v", err.Error()))
	})

	c.OnResponse(func(r *colly.Response) {
		logger.Logger.Info(fmt.Sprintf("Page visited: %v", r.Request.URL))
	})

	// iterating over the list of HTML product elements
	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		pokemonProduct := PokemonProduct{}

		pokemonProduct.Url = e.ChildAttr("a", "href")
		pokemonProduct.Image = e.ChildAttr("img", "src")
		pokemonProduct.Name = e.ChildText("h2.woocommerce-loop-product__title")
		pokemonProduct.Price = e.ChildText("span.price")

		PokemonProducts = append(PokemonProducts, pokemonProduct)
	})

	c.OnScraped(func(r *colly.Response) {
		logger.Logger.Info(fmt.Sprintf("%v scraped!", r.Request.URL))
	})

	c.Visit("https://scrapeme.live/shop/")

	AddToCSVFile()
}

func AddToCSVFile() {
	csvFilePath := filepath.Join(".", "csv_file")
	err := os.MkdirAll(csvFilePath, os.ModePerm)
	if err != nil {
		logger.Logger.Panic("unable to create log file folder")
	}
	csvFile, err := os.Create(fmt.Sprintf("%v/%v", csvFilePath, "products.csv"))
	if err != nil {
		logger.Logger.Panic("unable to create product csv file")
	}
	defer csvFile.Close()

	// initializing a file writer
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// defining the CSV headers
	headers := []string{
		"url",
		"image",
		"name",
		"price",
	}

	// writing the column headers
	writer.Write(headers)

	for _, pokemon := range PokemonProducts {
		record := []string{
			pokemon.Url,
			pokemon.Image,
			pokemon.Name,
			pokemon.Price,
		}

		writer.Write(record)
	}
}
