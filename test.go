package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"

    "golang.org/x/net/html"
)

func Test() {
    res, err := http.Get("http://paulgraham.com/articles.html")
    if err != nil {
        log.Fatal("Couldn't GET: \n", err)
    }
    defer res.Body.Close()

    doc, err := html.Parse(res.Body)
    if err != nil {
        log.Fatal("Couldn't parse HTML: \n", err)
    }

    var titles []string
    var urls []string

    var traverse func(n *html.Node)
    traverse = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "a" {
            for _, a := range n.Attr {
                if a.Key == "href" && strings.HasPrefix(a.Val, "") {
                    urls = append(urls, a.Val)
                    for c := n.FirstChild; c != nil; c = c.NextSibling {
                        if c.Type == html.TextNode {
                            titles = append(titles, c.Data)
                        }
                    }
                }
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            traverse(c)
        }
    }
    traverse(doc)

    fmt.Println("Titles:")
    for _, title := range titles {
        fmt.Println("- " + title)
    }

    fmt.Println("\nURLs:")
    for _, url := range urls {
        fmt.Println("- " + url)
    }
}

