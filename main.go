package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"webCrawler/crawler"
	"webCrawler/progressbar"
	"webCrawler/searcher"

	"os"

	"sync"
	"time"
	"webCrawler/indexer"
	"webCrawler/semaphor"
)

// Example of command line request: go run main.go https://www.economist.com/ 100 10 results.json 50

func main() {

	// Command line arguments
	Url := os.Args[1]
	MaxRoutines, err := strconv.Atoi(os.Args[2])
	MaxIndexingTime, err := strconv.Atoi(os.Args[3])
	FileName := os.Args[4]

	// Set timer and timeout context
	now := time.Now()
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(MaxIndexingTime)*time.Second)

	// Crawl
	s := progressbar.Crawl()
	visited := crawler.New()
	visited.Add(Url)
	url, err := crawler.ExtractNode(Url)
	if err != nil {
		fmt.Printf("Error getting page %s %s\n", url, err)
		return
	}
	out := crawler.CheckDuplicates(ctx, url, 50, visited)
	s.Stop()

	// Index results
	p := progressbar.Index()
	var concurrentJobs = semaphor.New(MaxRoutines)
	var wg sync.WaitGroup

	c := make(chan *indexer.Response)
	var responses []*indexer.Response

	for k := range out {
		wg.Add(1)
		concurrentJobs.Acquire()
		go indexer.ScrapePage(&wg, c, k, concurrentJobs)
	}

	go func() {
		wg.Wait()
		close(c)
	}()
	for val := range c {
		responses = append(responses, val)
	}

	p.Stop()
	progressbar.Finish(now)

	// Export indexed results to JSON file
	file, _ := json.MarshalIndent(responses, "", " ")

	if strings.HasSuffix(FileName, "json") {
		_ = ioutil.WriteFile(FileName, file, 0644)
	} else {
		_ = ioutil.WriteFile("results.json", file, 0644)
	}

	// Search by keyword
	progressbar.Search()
	userInput := searcher.ValidateInput()

	var data []indexer.Response

	if strings.HasSuffix(FileName, "json") {
		jsonFile, _ := ioutil.ReadFile("./" + FileName)
		_ = json.Unmarshal(jsonFile, &data)

	} else {
		jsonFile, _ := ioutil.ReadFile("./results.json")
		_ = json.Unmarshal(jsonFile, &data)
	}

	maps := make(map[string]int)

	for _, d := range data {
		counter := 0
		url := d.Url
		maps[url] = 0

		searcher.ScorePage(&d, userInput, &counter)
		maps[url] = counter
	}

	resultPair := make(searcher.PairSlice, len(maps))

	j := 0
	for k, v := range maps {
		resultPair[j] = searcher.Pair{k, v}
		j++
	}

	// Calculate MaxRelevance Score for each page.
	pagesMaxRelevance := make(map[string]int)
	for _, response := range responses {
		maxScore := searcher.MaxScorePage(response, userInput)
		pagesMaxRelevance[response.Url] = maxScore
	}

	// Sort pages by number of matches.
	sort.Sort(resultPair)

	// Pretty print sorted list.
	searcher.PrintResults(resultPair, pagesMaxRelevance)
}