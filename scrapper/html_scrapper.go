package scrapper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type ShoeInfo struct {
	NAME, PRICE, LINK string
}

type WebsiteSelectors struct {
	NameSelector  string
	PriceSelector string
	LinkSelector  string
	ItemSelector  string
	PageToScrape  string
}

var websiteSelectorsMap = map[string]WebsiteSelectors{
	"Adidas": {
		NameSelector:  ".glass-product-card__title",
		PriceSelector: ".gl-price-item",
		LinkSelector:  "a[href]",
		ItemSelector:  ".grid-item",
		PageToScrape:  "https://www.adidas.pt/calcado-futebol%7Cmontanhismo-homem?price_max=387&price_min=50&sort=newest-to-oldest",
	},
	"New Balance": {
		NameSelector:  ".pname",
		PriceSelector: ".sales",
		LinkSelector:  "a",
		ItemSelector:  ".product-tile",
		PageToScrape:  "https://www.newbalance.pt/pt/homens/tenis/",
	},
}

func extractName(e *colly.HTMLElement, selector string) string {
	return strings.TrimSpace(e.ChildText(selector))
}

func extractPrice(e *colly.HTMLElement, selector string) string {
	return strings.TrimSpace(e.ChildText(selector))
}

func extractLink(e *colly.HTMLElement, selector string) string {
	link := e.ChildAttr(selector, "href")
	if strings.HasPrefix(link, "http") {
		return link
	}
	return e.Request.AbsoluteURL(link)
}

func ScrapeProducts(selectorKey string) []ShoeInfo {
	websiteSelectors, exists := websiteSelectorsMap[selectorKey]
	if !exists {
		fmt.Printf("Website selector key '%s' not found\n", selectorKey)
		return nil
	}

	var Shoes []ShoeInfo

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	c.OnHTML(websiteSelectors.ItemSelector, func(e *colly.HTMLElement) {
		ShoeInfo := ShoeInfo{}
		ShoeInfo.NAME = extractName(e, websiteSelectors.NameSelector)
		ShoeInfo.PRICE = extractPrice(e, websiteSelectors.PriceSelector)
		ShoeInfo.LINK = extractLink(e, websiteSelectors.LinkSelector)

		Shoes = append(Shoes, ShoeInfo)
	})

	c.Visit(websiteSelectors.PageToScrape)
	return Shoes
}
