package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "strings"
    "encoding/xml" // use to unmarshal xml structure
)

type SitemapIndex struct { // https://a7.co/sitemap.xml
    // slice of strings type
    Locations []string `xml:"sitemap>loc"` // used for unmarshalling (tag)
}

// todo probably make an image struct for image:loc node
type Image struct {
    Locations []string `xml:"loc"`
    Titles    []string `xml:"title"`
    Captions  []string `xml:"caption"`
}

type Product struct { // https://a7.co/sitemap_products_1.xml?from=3783194241&to=4519876460653
    Locations       []string `xml:"url>loc"`
    LastModified    []string `xml:"url>lastmod"`
    ChangeFrequency []string `xml:"url>changefreq"`
    ImageField      []Image  `xml:"url>image"`
}

type Pages struct { // https://a7.co/sitemap_pages_1.xml
    Titles          []string `xml:"url>loc"`
    LastModified    []string `xml:"url>lastmod"`
    ChangeFrequency []string `xml:"url>changefreq"`
}

func main() {
    var s SitemapIndex
    var p Product

    ASevenRootSitemapUrl := "https://a7.co/sitemap.xml"
    resp, err := http.Get(ASevenRootSitemapUrl)
    if err != nil {
        fmt.Println("received err when requesting a7 sitemap root", err)
    }
    bytes, err := ioutil.ReadAll(resp.Body) // returns as bytes 
    if err != nil {
        fmt.Println("err reading bytes from response body", err)
    }
    // stringBody := string(bytes) // convert response bytes to string
    resp.Body.Close() // free up resources that made the request
    // fmt.Println(stringBody)

    // parse xml
    err = xml.Unmarshal(bytes, &s) // unmarshal bytes from response into memoery address of our sitemapindex struct 
    if err != nil {
        fmt.Println("err unmarshalling bytes into our SitemapIndex struct")
    }
    fmt.Println(s.Locations) // print out slice
    fmt.Printf("There are %d <loc>s in %s\n\n", len(s.Locations), ASevenRootSitemapUrl)


    for _, loc := range s.Locations {
        if strings.Contains(loc, "product") { // if .xml page has "product" in its url 
            resp, err := http.Get(loc)
            if err != nil {
                fmt.Printf("received err when requesting %s - %s\n", loc, err)
            }
            bytes, err := ioutil.ReadAll(resp.Body) // returns as bytes 
            if err != nil {
                fmt.Println("err reading bytes from response body", err)
            }
            resp.Body.Close()
            err = xml.Unmarshal(bytes, &p) // unmarshal bytes from response into memoery address of our product struct 
            if err != nil {
                fmt.Println("err unmarshalling bytes into our product struct")
            }

            // remove "https://a7.co/"
            var productUrls []string
            for _, productUrl := range p.Locations {
                if strings.Contains(productUrl, "/products/") { 
                    productUrls = append(productUrls, productUrl)
                }
            }
            fmt.Println(productUrls)
        }
    }
}
