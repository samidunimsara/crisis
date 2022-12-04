package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "net/http"
    "golang.org/x/net/html"
    "strings"
)

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
    if depth <= 0 {
        return
    }
    body, urls, err := fetcher.Fetch(url)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Check if the page contains a JavaScript file
    if strings.Contains(body, ".js") {
        // If the page contains a JavaScript file, search for endpoint URLs
        doc, err := html.Parse(strings.NewReader(body))
        if err != nil {
            fmt.Println(err)
            return
        }
        var f func(*html.Node)
        f = func(n *html.Node) {
            if n.Type == html.ElementNode && n.Data == "a" {
                for _, a := range n.Attr {
                    if a.Key == "href" && strings.HasPrefix(a.Val, "/") {
                        // If an endpoint URL is found, check if it is valid
                        endpoint := url + a.Val
                        resp, err := http.Head(endpoint)
                        if err != nil {
                            fmt.Println(err)
                            return
                        }
                        if resp.StatusCode == 200 {
                            // If the endpoint URL is valid, check if it returns JSON
                            if resp.Header.Get("Content-Type") == "application/json" {
                                // If the endpoint returns JSON, print the response body,
                                // HTTP status code, and content length
                                resp, err := http.Get(endpoint)
                                if err != nil {
                                    fmt.Println(err)
                                    return
                                }
                                defer resp.Body.Close()
                                var data map[string]interface{}
                                if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
                                    fmt.Println(err)
                                    return
                                }
                                fmt.Printf("Endpoint: %s\nHTTP Status: %d\nContent Length: %d\n",
                                    endpoint, resp.StatusCode, resp.ContentLength)
                                fmt.Printf("\x1b[31m%s\x1b[0m\n", data)
                            }
                        }
                    }
                }
            }
            for c := n.FirstChild; c != nil; c = c.NextSibling {
                f(c
