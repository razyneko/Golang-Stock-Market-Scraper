package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

func main() {

	type Stock struct {
		company, price, change string
	}

	ticker := []string{
		"MSFT",
		"IBM",
		"GE",
		"UNP",
		"COST",
		"MCD",
		"V",
		"WMT",
		"DIS",
		"MMM",
		"INTC",
		"AXP",
		"AAPL",
		"BA",
		"CSCO",
		"GS",
		"JPM",
		"CRM",
		"VZ",
	}

	stocks := []Stock{}

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnHTML("section.yf-xxbei9", func(e *colly.HTMLElement) {

		stock := Stock{}

		stock.company = e.ChildText("h1")
		fmt.Println("Company:", stock.company)

		c.OnHTML("div.price", func(e *colly.HTMLElement) {
			stock.price = e.ChildText("fin-streamer[data-field='regularMarketPrice']")
			fmt.Println("Price:", stock.price)
			stock.change = e.ChildText("fin-streamer[data-field='regularMarketChangePercent']")
			fmt.Println("Change:", stock.change)
		})

		stocks = append(stocks, stock)
	})

	c.Wait()

	for _, t := range ticker {
		c.Visit("https://finance.yahoo.com/quote/" + t + "/")
	}

	fmt.Println(stocks)

	file, err := os.Create("Stocks.csv")

	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}

	defer file.Close()
	writer := csv.NewWriter(file)

	headers := []string{
		"company",
		"price",
		"change",
	}

	writer.Write(headers)

	for _, stock := range stocks {
		record := []string{
			stock.company,
			stock.price,
			stock.change,
		}
		writer.Write(record)
	}

	defer writer.Flush()
}
