package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	//"net/http/httputil"
	//"net/url"
)

var severCount = 0

// These constant is used to define server
const (
	SERVER1 = "https://www.google.com"
	SERVER2 = "https://www.facebook.com"
	SERVER3 = "https://www.yahoo.com"
	PORT    = "1338"
)

var SERVER_LIST []string = []string{SERVER1, SERVER2, SERVER3}
var wg sync.WaitGroup
var mutex = sync.Mutex{}

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, slice *[]string) {

	resp, err := http.Get(target)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	mutex.Lock()
	*slice = append(*slice, string(body))
	mutex.Unlock()

	wg.Done()
}

// Log the typeform payload and redirect url
func logRequestPayload(proxyURL string) {
	log.Printf("proxy_url: %s\n", proxyURL)
}

// Given a request send it to the appropriate url
func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {

	wg.Add(len(SERVER_LIST))

	var slice []string

	i := 0
	for ; i < len(SERVER_LIST); i++ {
		url := SERVER_LIST[i]
		logRequestPayload(url)
		go serveReverseProxy(url, &slice)
	}

	wg.Wait()
	//
	// change below code
	//
	i = 0
	for ; i < len(slice); i++ {
		println(slice[i])
	}

	//
	io.WriteString(res, "Hello, world!\n")

}

func main() {
	// start server
	http.HandleFunc("/", handleRequestAndRedirect)

	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
