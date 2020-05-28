package crawler

import (
	"fmt"
	"sync"
)

func crawlChannelWorker(url string, depth int, fetcher Fetcher, quit chan bool, visitedUrls map[string]bool) {
	if depth <= 0 {
		quit <- true
		return
	}

	didIt, hasIt := visitedUrls[url]
	// If we have already visited this link,
	// stop here
	if didIt && hasIt {
		quit <- true
		return
	}
	// Mark it has visited
	visitedUrls[url] = true

	// Fetch children URLs
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		quit <- true
		return
	}
	fmt.Printf("found URL: %s ; title: %q\n", url, body)

	// Crawl children URLs
	var wg sync.WaitGroup
	wg.Add(len(urls))
	for _, childrenURL := range urls {
		defer wg.Done()
		go crawlChannelWorker(childrenURL, depth-1, fetcher, quit, visitedUrls)
		// To exit goroutines. This channel will always be filled
	}

	wg.Wait()

	quit <- true
}

// CrawlChannel uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func CrawlChannel(url string, depth int, fetcher Fetcher) {
	// Goal: Fetch URLs in parallel.
	// Goal: Don't fetch the same URL twice.

	quit := make(chan bool)

	// initial url not already fetched
	visitedUrls := map[string]bool{url: false}

	go crawlChannelWorker(url, depth, fetcher, quit, visitedUrls)

	// We will not quit until we have something
	// in the "quit" channel
	<-quit
}
