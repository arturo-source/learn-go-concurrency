package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

/* From the original problem:

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	if depth <= 0 {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		Crawl(u, depth-1, fetcher)
	}
	return
}
*/

type VisitedUrls struct {
	mu   sync.Mutex
	urls map[string]bool
}

func (vu *VisitedUrls) MarkAsVisited(url string) {
	vu.mu.Lock()
	defer vu.mu.Unlock()
	vu.urls[url] = true
}

func (vu *VisitedUrls) IsVisited(url string) bool {
	return vu.urls[url]
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, vUrls *VisitedUrls) {
	if depth <= 0 {
		return
	}

	urlsc := make(chan string)

	go func() {
		defer close(urlsc)

		vUrls.MarkAsVisited(url)
		body, urls, err := fetcher.Fetch(url)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("found: %s %q\n", url, body)

		// Append urls to channel if not visited
		for _, u := range urls {
			if !vUrls.IsVisited(u) {
				urlsc <- u
			}
		}
	}()

	for u := range urlsc {
		Crawl(u, depth-1, fetcher, vUrls)
	}
}

func main() {
	vUrls := VisitedUrls{urls: map[string]bool{}}
	Crawl("https://golang.org/", 4, fetcher, &vUrls)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
