package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
	"sync"
)

var fetched = struct {
	m map[string]error
	sync.Mutex
}{m: make(map[string]error)}



type urlType struct {
	url_name string
	m map[string]error
	sync.Mutex
}

type Edge struct {
	url_start string
	url_end string
}


var jobs = make(chan Edge, 10)


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



func get_hrefs(url string) (error, []string) {
	urls := make([]string, 0)
	resp, err := http.Get(url)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	if err != nil {
		panic(err)
	}
	// reads html as a slice of bytes
	html_body := resp.Body
	if err != nil {
		panic(err)
	}

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
				//fmt.Println(href)
			}
		}
	}
}
var loading = errors.New("url load in progress")

func crawl(url string, depth int, ch chan Edge){
	if depth <= 0{
		//fmt.Println("Done", url)
		return
	}

	fetched.Lock()
	if _, ok := fetched.m[url]; ok {
		fetched.Unlock()
		return
	}

	fetched.m[url] = loading
	fetched.Unlock()

	err, urls := get_hrefs(url)

	// And update the status in a synced zone.
	fetched.Lock()
	fetched.m[url] = err
	fetched.Unlock()

	if err != nil {
		//fmt.Printf("<- Error on %v: %v\n", url, err)
		return
	}
	//fmt.Printf("Found: %s\n", url)
	done := make(chan bool)
	for _, u := range urls {
		//fmt.Printf("-> Crawling child %v/%v of %v : %v.\n", i, len(urls), url, u)
		go func(url string) {
			crawl(url, depth-1, ch)
			done <- true
		}(u)
	}
	for _, u := range urls {
		//fmt.Printf("<- [%v] %v/%v Waiting for child %v.\n", url, i, len(urls), u)
		newEdge := Edge{url, u}
		ch <- newEdge
		<-done
	}
	//fmt.Printf("<- Done with %v\n", url)
}


func main(){

	url_string := "https://golang.org/"

	go func(url string) {
		crawl(url, 3, jobs)
		close(jobs)
	}(url_string)

	for pair := range jobs {
		fmt.Println(pair.url_start, pair.url_end)
	}
}


