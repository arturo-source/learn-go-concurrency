# Exercise: Web Crawler

> This exercise is taken from <https://go.dev/tour/concurrency/10>, in `main.go` you can find my solution.
>
> This solution uses goroutines recursively, be quiet, you may allocate more memory than is available.

In this exercise you'll use Go's concurrency features to parallelize a web crawler.

Modify the `Crawl` function to fetch URLs in parallel without fetching the same URL twice.

Hint: you can keep a cache of the URLs that have been fetched on a map, but maps alone are not safe for concurrent use!
