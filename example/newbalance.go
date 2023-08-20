package example

import (
	"strings"

	"github.com/gocolly/colly"
)

type ShoeInfo struct {
	NAME,
	PRICE,
	LINK string
}

func ScrapeProducts() []ShoeInfo {

	var Shoes []ShoeInfo
	pageToScrape := "https://www.newbalance.pt/pt/homens/tenis/"

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	c.OnHTML(".product-tile", func(e *colly.HTMLElement) {
		ShoeInfo := ShoeInfo{}
		ShoeInfo.NAME = strings.TrimSpace(e.ChildText(".pname"))
		ShoeInfo.PRICE = strings.TrimSpace(e.ChildText(".sales"))
		ShoeInfo.LINK = e.Request.AbsoluteURL(e.ChildAttr("a", "href"))

		Shoes = append(Shoes, ShoeInfo)

	})
	c.Visit(pageToScrape)
	return Shoes
}
