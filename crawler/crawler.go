package crawler

import (
	"context"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	neturl "net/url"
	"os"
	"strings"
	"sync"
)

type ConcurrentHashSet struct {
	mu     sync.Mutex
	values map[string]struct{}
}

func New() *ConcurrentHashSet {
	return &ConcurrentHashSet{values: make(map[string]struct{})}
}

func (s *ConcurrentHashSet) IsMember(val string) bool {
	s.mu.Lock()
	_, ok := s.values[val]
	s.mu.Unlock()
	return ok
}

func (s *ConcurrentHashSet) Add(val string) {
	s.mu.Lock()
	s.values[val] = struct{}{}
	s.mu.Unlock()
}

// ExtractNode fetches a URL passed as a string.
// It returns an HTML Node pointer.
func ExtractNode(url string) (*html.Node, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("cannot get page")
	}
	n, err := html.Parse(r.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot parse page")
	}
	return n, err
}

// Contains returns boolean True if slice s contains string j.
func Contains(s []string, j string) bool {
	for _, i := range s {
		if i == j {
			return true
		}
	}
	return false
}

var base = os.Args[1]

// FindLinks traverses recursively the pages connected to an HTML Node.
// Each recursion returns a list of links with no redundant link.
// Given the graph structure of a website, links may be redundant in the overall list.
func FindLinks(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				if !Contains(links, a.Val) {
					a.Val = CheckURL(a.Val, base)
					links = append(links, a.Val)
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = FindLinks(links, c)
	}
	return links
}

// CheckURL checks if the urls collected include a base path.
// If not, a base path is added to the incomplete URL.
func CheckURL(url, baseurl string) string {
	if strings.HasPrefix(url, baseurl) {
	}
	u, err := neturl.Parse(url)
	if err != nil {
		log.Fatal(err)
	}
	base, err := neturl.Parse(baseurl)
	if err != nil {
		log.Fatal(err)
	}
	url = base.ResolveReference(u).String()
	return url
}

// CheckDuplicates checks if each url in the overall list of collected urls is unique.
// It returns a channel with unique urls.
func CheckDuplicates(ctx context.Context, n *html.Node, maxNumber int, visited *ConcurrentHashSet) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		links := FindLinks(nil, n)
		if maxNumber > len(links) {
			fmt.Printf("Error: maxNumber value too high for the number of links processed.\n")
			os.Exit(1)
		}
		for i := 0; i < maxNumber; i++ {
			if visited.IsMember(links[i]) {
				continue
			}
			visited.Add(links[i])
			select {
			case out <- links[i]:
			case <-ctx.Done():
				fmt.Println("Timeout! Canceling CheckDuplicate")
				os.Exit(1)
			}

		}
	}()
	return out
}