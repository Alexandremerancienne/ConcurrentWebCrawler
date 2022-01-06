package progressbar

import (
	"fmt"
	spin "github.com/briandowns/spinner"
	"os"
	"time"
)

// Crawl starts a progressbar for crawler.
func Crawl() *spin.Spinner {
	fmt.Printf("==========   PHASE I: CRAWlING & INDEXING   ==========\n")
	s := spin.New(spin.CharSets[35], 100*time.Millisecond, spin.WithWriter(os.Stderr))
	s.Suffix = "  Crawling urls..."
	s.Start()
	return s
}

// Index starts a progressbar for indexer.
func Index() *spin.Spinner {
	fmt.Printf("\n")
	p := spin.New(spin.CharSets[35], 100*time.Millisecond, spin.WithWriter(os.Stderr))
	p.Suffix = "  Crawling completed! Indexing urls..."
	p.Start()
	return p
}

// Finish stops timer after crawling and indexing.
func Finish(now time.Time) {
	elapsed := time.Since(now)
	fmt.Println("\nProcess completed! Execution time: ", elapsed)

}

// Search starts search step.
func Search() {
	fmt.Printf("\n===============   PHASE II: SEARCHING   ===============\n")
}
