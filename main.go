package main

import (
	"./graph"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"strings"
	"sync"
)

var fetched = struct {
	m map[string]error
	sync.Mutex
}{m: make(map[string]error)}



func get_href(t html.Token) (ok bool, href string){
		for _, a := range t.Attr {
			if a.Key == "href" {
				href = a.Val
				if ! strings.HasPrefix(href, "https"){
					ok = false
				}else {
					ok = true
				}
			}
		}
	return
	}


func get_html_body(url string) io.ReadCloser{
	resp, err := http.Get(url)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	if err != nil {
		panic(err)
	}
	html_body := resp.Body
	return html_body
}


func get_hrefs(url string) (error, []string) {
	urls := make([]string, 0)
	// reads html as a slice of bytes
	html_body := get_html_body(url)

	z := html.NewTokenizer(html_body)
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return nil, urls
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "a"
			if isAnchor{
				ok, href := get_href(t)
				if !ok{
					continue
				}
				urls = append(urls, href)
			}
		}
	}
}
var loading = errors.New("url load in progress")

func crawl(currentUrl string, depth int, globalGraph *graph.Graph,wg *sync.WaitGroup){
	defer wg.Done()
	urlNode := graph.Node{Name: currentUrl}
	globalGraph.AddNode(urlNode)

	if depth <= 0{
		return
	}

	fetched.Lock()
	if _, ok := fetched.m[currentUrl]; ok {
		fetched.Unlock()
		return
	}

	fetched.m[currentUrl] = loading
	fetched.Unlock()

	err, urls := get_hrefs(currentUrl)

	// And update the status in a synced zone.
	fetched.Lock()
	fetched.m[currentUrl] = err
	fetched.Unlock()

	if err != nil {
		fmt.Printf("<- Error on %v: %v\n", currentUrl, err)
		return
	}
	for _, u := range urls {
		toNode := graph.Node{Name:u}
		globalGraph.AddNode(toNode)
		globalGraph.AddEdge(&urlNode, &toNode)
		//fmt.Printf("-> Crawling child %v/%v of %v : %v.\n", i, len(urls), url, u)
		wg.Add(1)
		go crawl(u, depth-1, globalGraph, wg)
	}

	//fmt.Printf("<- Done with %v\n", url)
}


func main(){

	var globalGraph graph.Graph
	var wg sync.WaitGroup

	url_string := "https://golang.org/"
	entryNode := graph.Node{Name: url_string}
	globalGraph.EntryNode = entryNode

	crawl(url_string, 3, &globalGraph, &wg)

	wg.Wait()
	fmt.Printf("first url has %d edges \n", len(globalGraph.Edges[entryNode]))

	globalGraph.Traverse(func(n *graph.Node) {
		fmt.Printf("%v\n", n)
	})



}



