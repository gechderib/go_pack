package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, j)
		// Simulate work like sending an email or processing data
		time.Sleep(time.Second)
		results <- j * 2
	}
}

func main() {

	sampleUrls := []string{
		"https://www.example.com",
		"https://www.google.com",
		"https://www.github.com",
		"https://www.reddit.com",
		"https://www.facebook.com",
		"https://www.twitter.com",
		"https://www.linkedin.com",
		"https://www.instagram.com",
		"https://www.youtube.com",
		"https://www.stackoverflow.com",
		"https://www.tiktok.com",
		"https://www.snapchat.com",
		"https://www.whatsapp.com",
		"https://www.telegram.com",
		"https://www.discord.com",
		"https://www.twitch.com",
		"https://www.pinterest.com",
		"https://www.medium.com",
		"https://www.quora.com",
		"https://www.wikipedia.org",
		"https://www.bing.com",
		"https://www.yahoo.com",
		"https://www.amazon.com",
		"https://www.ebay.com",
		"https://www.netflix.com",
		"https://www.spotify.com",
		"https://www.adobe.com",
		"https://www.dropbox.com",
		"https://www.slack.com",
		"https://www.gitlab.com",
		"https://www.bitbucket.com",
		"https://www.jenkins.io",
		"https://www.jfrog.com",
		"https://www.docker.com",
		"https://www.kubernetes.io",
		"https://www.terraform.io",
		"https://www.ansible.com",
		"https://www.puppet.com",
		"https://www.saltstack.com",
		"https://www.vagrantup.com",
		"https://www.hashicorp.com",
	}

	urls := make(chan string, 10)
	results := make(chan string, 10)
	var wg sync.WaitGroup

	for w := 1; w <= 10; w++ {
		wg.Add(1)
		go scrape(w, urls, results, &wg)
	}

	// for _, url := range sampleUrls {
	// 	fmt.Println("Sending URL to scrape:", url)
	// 	urls <- url
	// }
	go func() {
		for _, url := range sampleUrls {
			urls <- url
		}
		close(urls)
	}()
	// close(urls)
	go func() {
		wg.Wait()
		close(results)
	}()
	for res := range results {
		fmt.Println("Received result:", res)
	}

}
