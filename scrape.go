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
	"strings"

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

type Awards_Row struct {
	row_num                 string
	award_num               string
	delv_order_num          string
	delv_order_cnt          string
	last_mod_posting_date   string
	a_cage_code             string
	total_contract_price    string
	award_date              string
	posted_date             string
	nsn_part_number         string
	nomenclature            string
	purchase_request_number string
	solicitation_number     string
}

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	fmt.Printf("%v: %v\n", msg, time.Since(start))
}

func get_awards_page(browserCtx context.Context) (Awards_Page [][]string) {
	/*
		get the table data from the html string
		implement goquery to get the table data
		implement goquery object to csv
	*/

	defer duration(track("get_awards_page"))

	var html string
	var html_reader *strings.Reader

	// get the html string
	err := chromedp.Run(browserCtx,
		chromedp.OuterHTML("ctl00_cph1_grdAwardSearch", &html, chromedp.ByID),
	)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Got the page")
	}

	//fmt.Println(html[:1000])

	// read the html string
	html_reader = strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(html_reader)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Got the document")
	}

	// selection of rows then for each row select the cells and store the data
	doc.Find("tr.BgWhite, tr.BgSilver").Each(func(i int, row *goquery.Selection) {
		parsed_row := Awards_Row{}

		// selection of cells (td) in the row
		row = row.Find("td")
		// first cell in the row, span, text (might have issue with multiple a tags)
		parsed_row.row_num = row.Eq(0).Find("span").Text()
		// second cell in the row - span - second a - text
		parsed_row.award_num = row.Eq(1).Find("span").Find("a").Next().Text()
		// third cell in the row - span - second a - text
		parsed_row.delv_order_num = row.Eq(2).Find("span").Find("a").Next().Text()
		// fourth cell in the row - span - text
		parsed_row.delv_order_cnt = row.Eq(3).Find("span").Text()
		// fifth cell in the row - span - text
		parsed_row.last_mod_posting_date = row.Eq(4).Find("span").Text()
		// sixth cell in the row - span - a - text
		parsed_row.a_cage_code = strings.TrimSpace(row.Eq(5).Find("span").Find("a").Text())
		// seventh cell in the row - span - text
		parsed_row.total_contract_price = row.Eq(6).Find("span").Text()
		// eighth cell in the row - span - text
		parsed_row.award_date = row.Eq(7).Find("span").Text()
		// ninth cell in the row - span - text
		parsed_row.posted_date = row.Eq(8).Find("span").Text()
		// tenth cell in the row - span - text
		parsed_row.nsn_part_number = row.Eq(9).Find("span").Text()
		// eleventh cell in the row - span - text
		parsed_row.nomenclature = row.Eq(10).Find("span").Text()
		// twelfth cell in the row - span - text
		parsed_row.purchase_request_number = row.Eq(11).Find("span").Text()
		// thirteenth cell in the row - span - text
		parsed_row.solicitation_number = row.Eq(12).Find("span").Text()

		//fmt.Println("Row: ", parsed_row)

		parsed_row_list := []string{
			parsed_row.row_num,
			parsed_row.award_num,
			parsed_row.delv_order_num,
			parsed_row.delv_order_cnt,
			parsed_row.last_mod_posting_date,
			parsed_row.a_cage_code,
			parsed_row.total_contract_price,
			parsed_row.award_date,
			parsed_row.posted_date,
			parsed_row.nsn_part_number,
			parsed_row.nomenclature,
			parsed_row.purchase_request_number,
			parsed_row.solicitation_number,
		}
		// store the scraped row
		Awards_Page = append(Awards_Page, parsed_row_list)
	})

	/* finder := "tr.BgWhite, tr.BgSilver"
	fmt.Println("Finder: ", finder)

	//sel := doc.Find(finder).Find("td")
	sel := doc.Find("tr.BgWhite, tr.BgSilver").Find("td").First()
	for i := range sel.Nodes {
		single := sel.Eq(i)
		h, err := single.Html()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Node ", i, " text: ", h)
	}
	*/

	return Awards_Page
}

// TODO: implement the pagination
// for pagination, run the javascript then call the get_awards_page function
// and append what returns
// switch from struct type to list of strings
// TODO: implement Award_Page to CSV

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
	err := chromedp.Run(taskCtx, chromedp.Navigate("https://www.dibbs.bsm.dla.mil/Awards/"))
	if err != nil {
		log.Fatal(err)
	}

	// However, actions like Click don't always trigger a page navigation,
	// so they don't wait for a page load directly. Wrapping them with
	// RunResponse does that waiting, and also obtains the HTTP response.
	resp, err := chromedp.RunResponse(taskCtx, chromedp.Click("butAgree", chromedp.NodeVisible, chromedp.ByID))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("consent page status code:", resp.Status)

	// declare variables
	var date string
	var target_url string
	var last_page string
	date = "08-13-2024"
	target_url = "https://www.dibbs.bsm.dla.mil/Awards/AwdRecs.aspx?Category=awddt&TypeSrch=cq&Value=" + date

	// navigate to the target url then target date page and get table 1 of records
	err = chromedp.Run(taskCtx,
		chromedp.Navigate("https://www.dibbs.bsm.dla.mil/Awards/AwdDates.aspx?category=awddt"),
		chromedp.Navigate(target_url),
		chromedp.Evaluate(`javascript:__doPostBack('ctl00$cph1$grdAwardSearch','Page$Last')`, nil),
		chromedp.InnerHTML(".pagination span", &last_page, chromedp.BySearch),
		chromedp.Navigate(target_url),
	)
	if err != nil {
		log.Fatal(err)
		fmt.Println("Error navigating to the target url for ", date, " or table 1 ")
	} else {
		fmt.Println("Last page: ", last_page)
		fmt.Println("Navigated to the first page of ", date)
		page_1 := get_awards_page(taskCtx)
		fmt.Println("Page 1: ", page_1[:1])
	}

	// get to table 2 of records
	resp, err = chromedp.RunResponse(taskCtx, chromedp.Evaluate(`javascript:__doPostBack('ctl00$cph1$grdAwardSearch','Page$2')`, nil))
	if err != nil || resp.Status != 200 {
		log.Fatal(err)
		fmt.Println("Error retrieving table 2. status code: ", resp.Status)
	} else {
		fmt.Println("Page 2 status code: ", resp.Status)
		page_2 := get_awards_page(taskCtx)
		fmt.Println("Page 2: ", page_2[:1])
	}

}
