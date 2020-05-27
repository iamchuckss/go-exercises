package main

import (
	"fmt"
	crawler "go-exercises/tour-of-go/12-web-crawler"
)

//
// main
//

func main() {

	fmt.Printf("=== Serial===\n")
	crawler.Serial("http://golang.org/", fetcher, make(map[string]bool))

	fmt.Printf("=== ConcurrentMutex ===\n")
	crawler.ConcurrentMutex("http://golang.org/", fetcher, makeState())

	fmt.Printf("=== ConcurrentChannel ===\n")
	crawler.ConcurrentChannel("http://golang.org/", fetcher)

	fmt.Printf("=== CrawlChannel ===\n")
	crawler.CrawlChannel("https://golang.org/", 4, fetcher)
}
