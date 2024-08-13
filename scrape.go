package main

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	//"io"
	"log"
	"net/http"
	"os"
)

// custom type to keep of the target object to scrape
type Product struct {
	name, price string
}

func main() {
	// download the target HTML document
	// Get() method performs an HTTP GET request to the destination page because net/http acts as an HTTP
	// client. Server will respond with the HTTP document of the page in the response body.
	response, err := http.Get("https://www.scrapingcourse.com/ecommerce/")
	if err != nil {
		log.Fatal("Failed to connect to the target page", err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Fatalf("HTTP Error %d: %s", response.StatusCode, response.Status)
	}

	// parse the HTML document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Failed to parse the HTML document", err)
	}

	// where to store the scraped data
	var products []Product

	// retrieve name and price from each product and print it
	doc.Find("li.product").Each(func(i int, p *goquery.Selection) {
		// scraping logic
		product := Product{}
		product.name = p.Find("h2").Text()
		product.price = p.Find("span.price").Text()

		// store the scraped product
		products = append(products, product)
	})
	fmt.Println(products)

	// print the scraped data
	// initialize the output CSV file
	file, err := os.Create("products.csv")
	if err != nil {
		log.Fatal("Failed to create the output CSV file", err)
	}
	defer file.Close()

	// initialize a file writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// define the CSV headers
	headers := []string{
		"image",
		"price",
	}
	// write the column headers
	writer.Write(headers)

	// add each product to the CSV file
	for _, product := range products {
		// convert a Product to an array of strings
		record := []string{
			product.name,
			product.price,
		}

		// write a new CSV record
		writer.Write(record)
	}
}
