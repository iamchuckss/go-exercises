package crawler

import (
	"fmt"
	"sync"
)

//
// Serial crawler
//

func Serial(url string, fetcher Fetcher, fetched map[string]bool) {
	if fetched[url] {
		return
	}
	fetched[url] = true
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		return
	}
	fmt.Printf("found URL: %s ; title: %q\n", url, body)
	for _, u := range urls {
		Serial(u, fetcher, fetched)
	}
	return
}

//
// Concurrent crawler with shared state and Mutex
//

type fetchState struct {
	mu      sync.Mutex
	fetched map[string]bool
}

func ConcurrentMutex(url string, fetcher Fetcher, f *fetchState) {
	f.mu.Lock()
	already := f.fetched[url]
	f.fetched[url] = true
	f.mu.Unlock()

	if already {
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		return
	}
	fmt.Printf("found URL: %s ; title: %q\n", url, body)
	var done sync.WaitGroup
	for _, u := range urls {
		done.Add(1)
		//     u2 := u
		// go func() {
		// 	defer done.Done()
		// 	ConcurrentMutex(u2, fetcher, f)
		// }()
		go func(u string) {
			defer done.Done()
			ConcurrentMutex(u, fetcher, f)
		}(u)
	}
	done.Wait()
	return
}

func MakeState() *fetchState {
	f := &fetchState{}
	f.fetched = make(map[string]bool)
	return f
}

//
// Concurrent crawler with channels
//

func worker(url string, ch chan []string, fetcher Fetcher) {
	body, urls, err := fetcher.Fetch(url)
	fmt.Printf("found URL: %s ; title: %q\n", urls, body)
	if err != nil {
		ch <- []string{}
	} else {
		ch <- urls
	}
}

func master(ch chan []string, fetcher Fetcher) {
	n := 1
	fetched := make(map[string]bool)
	for urls := range ch {
		for _, u := range urls {
			if fetched[u] == false {
				fetched[u] = true
				n++
				go worker(u, ch, fetcher)
			}
		}
		n--
		if n == 0 {
			break
		}
	}
}

func ConcurrentChannel(url string, fetcher Fetcher) {
	ch := make(chan []string)
	go func() {
		ch <- []string{url}
	}()
	master(ch, fetcher)
}
