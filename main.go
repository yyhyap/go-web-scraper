package main

import (
	"fmt"
	"go-web-scraper/logger"

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

	c.Visit("https://scrapeme.live/shop/")

	for _, p := range PokemonProducts {
		logger.Logger.Info(fmt.Sprintf("Pokemon: %v", p))
	}
}
