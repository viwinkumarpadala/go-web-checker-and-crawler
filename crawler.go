
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/jackdanger/collectlinks"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

var mu = &sync.Mutex{}

func main() {
	flag.Parse()
	
	args := flag.Args()
	fmt.Println(args)
	if len(args) < 1 {
		fmt.Println("Please specify base URL to crawl ")
		os.Exit(1)
	}
	 _, err := url.ParseRequestURI(args[0])
	 if err != nil {
	 	fmt.Println("Please enter a valid URL to crawl")
		os.Exit(2)
	}


	queue := make(chan string)
	filteredQueue := make(chan string)

	go func() { queue <- args[0] }()
	go filterQueue(queue, filteredQueue)

	
	done := make(chan bool)

	for i := 0; i < 5; i++ {
		go func() {
			for uri := range filteredQueue {
				addToQueue(uri, queue)
			}
			done <- true
		}()

	}
	<-done
}

func filterQueue(in chan string, out chan string) {
	
	var seen = make(map[string]bool)
	for val := range in {
		if !seen[val] {
			seen[val] = true
			out <- val
		}
	}
}



func addToQueue(uri string, queue chan string) {

	start := time.Now()

	
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	
	client := http.Client{Transport: transport}
	resp, err := client.Get(uri)
	if err != nil {
		return
	}
	
	defer resp.Body.Close()

	links := collectlinks.All(resp.Body)
	foundUrls := []string{}
	for _, link := range links {
		
		absolute := cleanUrl(link, uri)
		foundUrls = append(foundUrls, absolute)       
		if uri != "" {
			go func() { queue <- absolute }()
		}
	}
	stop := time.Now()
	display(uri,foundUrls,start,stop)
}

func display(uri string,found []string,start time.Time,stop time.Time){
 	mu.Lock()
	fmt.Println("Start time of crawl of this URL:",start)
	fmt.Println("Stop time of crawl of this URL :",stop)
	fmt.Println(uri)

	for _,str := range found{
		str, err := url.Parse(str)
		if err== nil{
			if str.Scheme =="http" || str.Scheme == "https"{
				fmt.Println("\t", str)
			}
		}
	}
	mu.Unlock()
}



func cleanUrl(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	return uri.String()
}