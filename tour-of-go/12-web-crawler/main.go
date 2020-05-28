package main

import (
	"fmt"
	crawler "go-exercises/tour-of-go/12-web-crawler/crawler"
)

//
// main
//

func main() {

	fmt.Printf("=== Serial===\n")
	crawler.Serial("http://golang.org/", crawler.FetcherImpl, make(map[string]bool))

	fmt.Printf("=== ConcurrentMutex ===\n")
	crawler.ConcurrentMutex("http://golang.org/", crawler.FetcherImpl, crawler.MakeState())

	fmt.Printf("=== ConcurrentChannel ===\n")
	crawler.ConcurrentChannel("http://golang.org/", crawler.FetcherImpl)

	fmt.Printf("=== CrawlChannel ===\n")
	crawler.CrawlChannel("https://golang.org/", 4, crawler.FetcherImpl)
}
