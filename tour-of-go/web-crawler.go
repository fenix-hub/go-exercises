// Ref: https://go.dev/tour/concurrency/10

package main

import (
	"fmt"
	"sync"
)


type Fetcher interface {
	Fetch(url string) (string, []string, error)
}

type FetchResult struct { 
	body string
	urls []string
	err error 
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	defer wg.Done() // gets executed at the end of Crawl
	
	if depth <= 0 {
		return
	}
	if cache.Check(url) {
		fmt.Printf("Hit! %s\n", url)
		return
	}
	
	body, urls, err := fetcher.Fetch(url)
	cache.Set(url)
	
	if err != nil {
		fmt.Println(err)
		return
	}
	
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		wg.Add(1)
		go Crawl(u, depth-1, fetcher)
	}
	return
}


type Cache struct {
	v map[string]bool // Url -> Body
	m sync.Mutex
}

func (c *Cache) Check(key string) bool {
	c.m.Lock()
	defer c.m.Unlock()
	return c.v[key]  
}

func (c *Cache) Set(key string) {
	c.m.Lock()
	c.v[key] = true
	c.m.Unlock()
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

var cache Cache = Cache{ v: make(map[string]bool) }
var wg sync.WaitGroup     // required to let all launched goroutines finish before main() exits

func main() {
	wg.Add(1)
	go Crawl("https://golang.org/", 4, fetcher)
	wg.Wait()
}
