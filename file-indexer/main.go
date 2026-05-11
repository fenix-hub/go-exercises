package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"
)

var wordRegex = regexp.MustCompile(`[a-zA-Z0-9]+`)

type Index struct {
	mu    sync.Mutex
	words map[string]int
}

type WordCount struct {
	word  string
	count int
}

func (i *Index) incr(word string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.words[word]++
}

func (i *Index) topN(n int) []WordCount {
	i.mu.Lock()
	defer i.mu.Unlock()

	var t []WordCount
	for word, count := range i.words {
		t = append(t, WordCount{word, count})
	}

	sort.Slice(t, func(i, j int) bool {
		return t[i].count > t[j].count
	})

	return t[0:n]
}

func index(id int, index *Index, chFiles <-chan string) error {
	for filePath := range chFiles {
		fmt.Printf("Worker %d: indexing %s\n", id, filePath)

		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Worker %d: error reading %s: %v\n", id, filePath, err)
			continue // Non bloccare per un singolo file
		}

		text := strings.ToLower(string(content))
		words := wordRegex.FindAllString(text, -1)

		for _, word := range words {
			if len(word) > 1 {
				index.incr(word)
			}
		}
	}
	return nil
}

var extensions = []string{".md", ".txt", ".json", ".log"}

func crawl(folder string, chFiles chan<- string) error {
	fmt.Println("Crawling:", folder)

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return nil
		}

		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if slices.Contains(extensions, ext) {
				chFiles <- path
			}
		}
		return nil
	})

	return err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <directory>")
		os.Exit(1)
	}

	folder := os.Args[1]
	ind := Index{words: make(map[string]int)}
	chFiles := make(chan string, 100)

	var wg sync.WaitGroup

	t0 := int64(0)

	// Crawler goroutine
	wg.Go(func() {
		defer close(chFiles)

		t0 = time.Now().UnixMilli()
		err := crawl(folder, chFiles)
		if err != nil {
		}

	})

	// Indexer goroutines
	for i := 1; i <= 5; i++ {
		workerID := i
		wg.Go(func() {
			err := index(workerID, &ind, chFiles)
			if err != nil {
			}
		})
	}

	wg.Wait()

	elapsed := time.Now().UnixMilli() - t0
	fmt.Printf("\nElapsed time: %.2f s\n", float64(elapsed)/1000.0)

	// Top 10
	top := ind.topN(10)
	fmt.Println("\n==== COMPLETED ====")
	fmt.Printf("Unique words: %d\n", len(ind.words))

	fmt.Println("\nTop 10 words:")
	for i := 0; i < 10 && i < len(top); i++ {
		fmt.Printf("%2d. %-20s %d\n", i+1, top[i].word, top[i].count)
	}
}
