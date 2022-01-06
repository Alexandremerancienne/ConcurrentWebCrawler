package searcher

import (
	"bufio"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"math"
	"os"
	"strconv"
	"strings"
	"webCrawler/indexer"
)

type PairSlice []Pair

type Pair struct {
	Key   string
	Value int
}

// ValidateInput validates user's input and returns a slice of strings.
func ValidateInput() []string {
	fmt.Printf("Search in browser: ")
	scanner := bufio.NewScanner(os.Stdin)
	s := []string{""}
	for scanner.Scan() {
		s = strings.Fields(scanner.Text())
		break
	}
	return s
}

// Len returns the length of an instance of PairSlice
func (p PairSlice) Len() int {
	return len(p)
}

// Swap swaps two Pair instances in PairSlice
func (p PairSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Less sorts two Pair instances in PairSlice
func (p PairSlice) Less(i, j int) bool {
	return p[i].Value < p[j].Value
}

// CountTextInSlice counts the occurrences of a word in a slice of strings.
// Counter is incremented by x for each match.
func CountTextInSlice(userInput []string, field *[]string, counter *int, x int) {
	for _, j := range userInput {
		for _, k := range *field {
			if strings.Contains(strings.ToLower(k), strings.ToLower(j)) {
				*counter = *counter + x
			}
		}
	}
}

// ScorePage returns a score for an HTML page.
// It reflects the relevance of a page, given search words
// And the weighting criteria defined by function CountTextInSlice.
func ScorePage(r *indexer.Response, userInput []string, counter *int) {
	var title = []string{r.Title}
	CountTextInSlice(userInput, &r.Keyword, counter, 6)
	CountTextInSlice(userInput, &r.Name, counter, 5)
	CountTextInSlice(userInput, &r.H1, counter, 4)
	CountTextInSlice(userInput, &r.H2, counter, 3)
	CountTextInSlice(userInput, &title, counter, 2)
	CountTextInSlice(userInput, &r.Alt, counter, 1)
}

// maxCountTextInSlice returns the maximum possible score for a tag field.
// The maximum is calculated admitting each string related to a tag
// Has at least one match with each word of user's input.
// For instance, all subheaders h2 have a match with each word of user's input.
func maxCountTextInSlice(userInput []string, field *[]string, x int) int {
	counter := 0
	for i := 0; i < len(userInput); i++ {
		for j := 0; j < len(*field); j++ {
			counter = counter + x
		}
	}
	return counter
}

// MaxScorePage returns the maximum possible score for a page.
// The maximum is calculated aggregating all maximum scores
// For tags screened with maxCountTextInSlice function.
func MaxScorePage(r *indexer.Response, userInput []string) int {
	finalCount := 0
	var title = []string{r.Title}
	count1 := maxCountTextInSlice(userInput, &r.Keyword, 6)
	count2 := maxCountTextInSlice(userInput, &r.Name, 5)
	count3 := maxCountTextInSlice(userInput, &r.H1, 4)
	count4 := maxCountTextInSlice(userInput, &r.H2, 3)
	count5 := maxCountTextInSlice(userInput, &title, 2)
	count6 := maxCountTextInSlice(userInput, &r.Alt, 1)
	finalCount = finalCount + count1 + count2 + count3 + count4 + count5 + count6
	return finalCount
}

// PrintResults prints the 10 most relevant resources.
func PrintResults(r PairSlice, pages map[string]int) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Url", "Relevance (%)"})

	fmt.Printf("\nMost relevant results:\n\n")
	for i := 0; i < 10; i++ {
		x := float64((r[len(r)-1-i]).Value)
		y := float64(pages[(r[len(r)-1-i]).Key])
		l := math.Floor((x/y)*100000) / 1000
		str := strconv.FormatFloat(l, 'f', 2, 64)
		t.AppendRow([]interface{}{(r[len(r)-1-i]).Key, str})
		t.AppendSeparator()
	}

	t.SetStyle(table.StyleLight)
	t.Render()
}