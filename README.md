# ConcurrentWebCrawler
A Golang concurrent Web crawler with data aggregation using goroutines.

# Command-line arguments
The application receives from command line: 
1) The URL of a web page; 
2) A MaxRoutines number - the maximal number of goroutines;
3) A MaxIndexingTime - the maximal allowed time for indexing;
4) A base result files name (default "results.json”);
5) A MaxResults number - the maximum number of resources processed (crawled and indexed)

Example of command: `go run main.go https://www.economist.com/ 100 10 results.json 50`

# Working steps

The application encompasses two phases:

1) The crawler indexes all relevant results based on the parameters given as command-line arguments;
2) The application asks for keywords and returns the most relevant results. 

## Phase I: Indexing
- The program traverses recursively the connected pages starting from the URL provided and following the hyperlinks in these web pages. 
- Separate goroutines are used to speed up the traversal process. 
- The program uses Breadth-first-search to prioritize the search pages.
- The program  extracts information about the keywords mentioned in each page using a custom extraction and weighting criteria.
- Keywords and names found in headers and sub-headers (h1 and h2), for instance, are more relevant than the others.
- Other weighting criteria have been defined for the first 5000 characters of the web page body.
- The program aggregates the extracted information about the found web pages. 
- A JSON file is automaticcaly generated with collected information. The file is named using the name provided in command-line arguments (default name "results.json”).

## Phase II: Search
- The program takes a list of search keywords given as input and find the most relevant web pages to keywords. 
- Result is given as a sorted list of most relevant top 10 resources with relevance percentage. 
