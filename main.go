package main

import (
	"fmt"
	"go-web-scraper/logger"

	"github.com/gocolly/colly"
)

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

	c.Visit("https://scrapeme.live/shop/")
}
