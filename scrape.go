package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"time"

	//"io"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

// custom type to keep of the target object to scrape
type Product struct {
	name, price string
}

func new() {
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

func main() {
	// setup options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.DisableGPU,
		chromedp.NoSandbox,
	)

	// create chrome instance
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// also set up a custom logger
	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// create a timeout
	taskCtx, cancel = context.WithTimeout(taskCtx, 15*time.Second)
	defer cancel()

	// ensure that the browser process is started
	if err := chromedp.Run(taskCtx); err != nil {
		log.Fatal(err)
	}

	// ## DO WORK ##
	// The Navigate action already waits until a page loads.
	err := chromedp.Run(taskCtx, chromedp.Navigate(`https://www.dibbs.bsm.dla.mil/Awards/`))
	if err != nil {
		log.Fatal(err)
	}

	// However, actions like Click don't always trigger a page navigation,
	// so they don't wait for a page load directly. Wrapping them with
	// RunResponse does that waiting, and also obtains the HTTP response.
	resp, err := chromedp.RunResponse(taskCtx, chromedp.Click(`butAgree`, chromedp.NodeVisible, chromedp.ByID))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("consent page status code:", resp.Status)

	// get page 1 table
	var html string
	var html2 string
	err = chromedp.Run(taskCtx,
		chromedp.Navigate(`https://www.dibbs.bsm.dla.mil/Awards/AwdDates.aspx?category=awddt`),
		chromedp.Navigate(`https://www.dibbs.bsm.dla.mil/Awards/AwdRecs.aspx?Category=awddt&TypeSrch=cq&Value=08-13-2024`),
		chromedp.InnerHTML(`ctl00_cph1_grdAwardSearch`, &html, chromedp.ByID),
	)
	if err != nil {
		log.Fatal(err)
	}

	// go to page 2
	resp, err = chromedp.RunResponse(taskCtx, chromedp.Evaluate(`javascript:__doPostBack('ctl00$cph1$grdAwardSearch','Page$2')`, nil))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("page 2 status code:", resp.Status)

	// get page 2 table
	err = chromedp.Run(taskCtx,
		chromedp.InnerHTML(`ctl00_cph1_grdAwardSearch`, &html2, chromedp.ByID),
	)
	if err != nil {
		log.Fatal(err)
	}

	// loop through the links with a pool
	// "javascript:__doPostBack('ctl00$cph1$grdAwardSearch','Page$X')" 1-100 etc.
	fmt.Println(html[:100])
	fmt.Println("nexttttttt...")
	fmt.Println(html2[:100])
}
