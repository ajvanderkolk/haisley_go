package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

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

	// convert the response buffer to bytes
	byteBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error while reading the response buffer", err)
	}

	// convert the byte HTML content to string and
	// print it
	html := string(byteBody)
	fmt.Println(html)

}
