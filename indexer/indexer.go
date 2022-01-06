package indexer

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"webCrawler/semaphor"
)

type Response struct {
	Url     string   `json:"url"`
	Name    []string `json:"names"`
	Keyword []string `json:"keywords"`
	H1      []string `json:"h1"`
	H2      []string `json:"h2"`
	Title   string   `json:"title"`
	Alt     []string `json:"alt"`
}

// ScrapePage scrapes an HTML page and collects information for indexing.
// Keywords and names in the headers, subheaders h1 and h2,
// title and alt tags are considered more relevant for indexing.
func ScrapePage(wg *sync.WaitGroup, c chan *Response, url string, jobs semaphor.Semaphor) {

	defer wg.Done()
	res, err := http.Get(url)
	defer res.Body.Close()

	if err != nil || res.StatusCode != 200 {
		fmt.Printf("Error getting page %s %s\n", url, err)
		os.Exit(1)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	body := ""
	response := new(Response)
	response.Url = url

	doc.Find("html").Each(func(i int, s *goquery.Selection) {
		s.Find("body, h1, h2, img").Each(func(i int, t *goquery.Selection) {
			tag := t.Nodes[0].Data
			switch tag {
			case "body":
				bodyText := t.Text()
				if len(bodyText) > 5000 {
					body = bodyText[:5001] // Keep the first 5000 characters of body
				} else {
					body = bodyText
				}
			case "h1":
				h1 := ReplaceStop(t.Text())
				if strings.Contains(body, strings.TrimSpace(h1)) {
					response.H1 = append(response.H1, strings.TrimSpace(h1))
				}
			case "h2":
				h2 := ReplaceStop(t.Text())
				if strings.Contains(body, strings.TrimSpace(h2)) {
					response.H2 = append(response.H2, strings.TrimSpace(h2))
				}
			case "img":
				altImage, _ := t.Attr("alt")
				a := ReplaceStop(altImage)
				if strings.Contains(body, strings.TrimSpace(altImage)) {
					response.Alt = append(response.Alt, a)
				}
			}
		})

		s.Find("head").Each(func(i int, h *goquery.Selection) {
			h.Find("meta, title").Each(func(i int, t *goquery.Selection) {
				tag := t.Nodes[0].Data
				switch tag {
				case "meta":
					if name, _ := t.Attr("name"); strings.TrimSpace(name) != "" {
						n, _ := t.Attr("content")
						if strings.TrimSpace(n) != "" {
							n = ReplaceStop(n)
							response.Name = append(response.Name, n)
						}
					} else if keyword, _ := t.Attr("keywords"); keyword != "" {
						k, _ := t.Attr("content")
						k = ReplaceStop(k)
						response.Keyword = append(response.Keyword, k)
					}
				case "title":
					a := ReplaceStop(t.Text())
					response.Title = strings.TrimSpace(a)
				}
			})
		})
	})

	c <- response
	jobs.Release()
}

// ReplaceStop replaces "stop words" by spaces in a page.
// Ie returns the modified input text.
func ReplaceStop(txt string) string {
	stopWords := []string{" - ", " a ", " is ", " the ", " an ", " & ",
		" and ", " are ", " as ", " at ", " be ", " but ",
		" by ", " for ", " if ", " in ", " into ", " it ",
		" no ", " not ", " of ", " on ", " or ", " such ",
		" that ", " their ", " then ", " there ", " these ",
		" they ", " this ", " to ", " was ", " will ", " with ",
		"A ", "Is ", "The ", "An ",
		"And ", "Are ", "As ", "At ", "Be ", "But ",
		"By ", "For ", "If ", "In ", "Into ", "It ",
		"No ", "Not ", "Of ", "On ", "Or ", "Such ",
		"That ", "Their ", "Then ", "There ", "These ",
		"They ", "This ", "To ", "Was ", "Will ", "With "}
	for _, word := range stopWords {
		r := strings.NewReplacer(word, " ")
		txt = r.Replace(txt)
	}
	return txt
}
