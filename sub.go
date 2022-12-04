package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
)

// Result is the structure of a single result returned by the crt.sh API
type Result struct {
    NameValue string `json:"name_value"`
}

// Results is the structure of the response from the crt.sh API
type Results struct {
    Results []Result `json:"results"`
}

func main() {
    // define the command-line flag for the domain name
    var domainFlag = flag.String("u", "", "domain name to search for subdomains")
    flag.Parse()

    // check if the domain name flag was provided
    if *domainFlag == "" {
        fmt.Println("error: domain name not provided")
        return
    }

    // construct the URL for the crt.sh API
    url := fmt.Sprintf("https://crt.sh/?q=%.%s&output=json", *domainFlag)

    // make an HTTP request to the crt.sh API
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("error:", err)
        return
    }
    defer resp.Body.Close()

    // read the response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("error:", err)
        return
    }

    // parse the JSON response
    var results Results
    err = json.Unmarshal(body, &results)
    if err != nil {
        fmt.Println("error:", err)
        return
    }

    // extract the subdomains from the response
    subdomains := make([]string, 0)
    for _, result := range results.Results {
        subdomain := strings.TrimPrefix(result.NameValue, "*.")
        subdomains = append(subdomains, subdomain)
    }

    // print the subdomains to the terminal
    for _, subdomain := range subdomains {
        fmt.Println(subdomain)
    }
}
