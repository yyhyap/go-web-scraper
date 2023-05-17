package main

import (
	"encoding/csv"
	"fmt"
	"go-web-scraper/logger"
	"os"
	"path/filepath"
	"sync"

	"github.com/gocolly/colly"
)

var (
	PokemonProducts []PokemonProduct
	mutex           sync.Mutex
)

type PokemonProduct struct {
	Url   string
	Image string
	Name  string
	Price string
}

func main() {
	var pageToScrape string
	pagesToScrape := []string{}
	pagesDiscovered := map[string]struct{}{}
	pagesDiscovered["https://scrapeme.live/shop/"] = struct{}{}

	// current iteration
	i := 1
	// max pages to scrape
	maxLimit := 5

	c := colly.NewCollector(
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		// limit the parallel requests to 4 request at a time
		Parallelism: 4,
	})

	// setting a valid User-Agent header
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	c.OnRequest(func(r *colly.Request) {
		logger.Logger.Info(fmt.Sprintf("Visiting: %v", r.URL))
	})

	c.OnError(func(_ *colly.Response, err error) {
		logger.Logger.Error(fmt.Sprintf("error when crawling: %v", err.Error()))
	})

	c.OnResponse(func(r *colly.Response) {
		logger.Logger.Info(fmt.Sprintf("Page visited: %v", r.Request.URL))
	})

	// iterating over the list of pagination links to implement the crawling logic
	c.OnHTML("a.page-numbers", func(e *colly.HTMLElement) {
		// discovering a new page
		newPaginationLink := e.Attr("href")

		mutex.Lock()
		if _, exist := pagesDiscovered[newPaginationLink]; !exist {
			pagesToScrape = append(pagesToScrape, newPaginationLink)
			pagesDiscovered[newPaginationLink] = struct{}{}
		}
		mutex.Unlock()
	})

	// iterating over the list of HTML product elements
	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		pokemonProduct := PokemonProduct{}

		pokemonProduct.Url = e.ChildAttr("a", "href")
		pokemonProduct.Image = e.ChildAttr("img", "src")
		pokemonProduct.Name = e.ChildText("h2.woocommerce-loop-product__title")
		pokemonProduct.Price = e.ChildText("span.price")

		mutex.Lock()
		PokemonProducts = append(PokemonProducts, pokemonProduct)
		mutex.Unlock()
	})

	c.OnScraped(func(res *colly.Response) {
		var newPageToScrape *string = nil

		mutex.Lock()
		if len(pagesToScrape) > 0 && i < maxLimit {
			pageToScrape = pagesToScrape[0]
			pagesToScrape[0] = ""
			pagesToScrape = pagesToScrape[1:]

			newPageToScrape = &pageToScrape

			// increment the counter
			i++
		}
		mutex.Unlock()

		if newPageToScrape != nil {
			c.Visit(*newPageToScrape)
		}

		logger.Logger.Info(fmt.Sprintf("%v scraped!", res.Request.URL))
	})

	c.Visit("https://scrapeme.live/shop/")

	c.Wait()

	AddToCSVFile()
}

func AddToCSVFile() {
	csvFilePath := filepath.Join(".", "csv_file")
	err := os.MkdirAll(csvFilePath, os.ModePerm)
	if err != nil {
		logger.Logger.Panic("unable to create csv file folder")
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
