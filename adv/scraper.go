package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func scrape(id int, urls <-chan string, results chan<- string, wg *sync.WaitGroup) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	defer wg.Done()
	for url := range urls {
		res, err := client.Get(url)
		if err != nil {
			fmt.Printf("Worker %d: Error scraping %s: %v\n", id, url, err)
			results <- fmt.Sprintf("Worker %d: Error scraping %s: %v", id, url, err)
			continue
		}
		res.Body.Close()
		fmt.Printf("Worker %d: Scraped %s: %s\n", id, url, res.Status)
		results <- fmt.Sprintf("Worker %d: Scraped %s: %s", id, url, res.Status)
	}
}
