package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func main() {
    res, err := http.Get("http://paulgraham.com/articles.html")
    if err != nil {
        log.Fatal("Couldn't GET: \n", err)
    }
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        log.Fatal("res.Status not OK:\n", res.Status)
    }

    document, err := html.Parse(res.Body)
    if err != nil {
        log.Fatal("Couldn't read body: \n", err)
    }

    articleLinks := getArticleLinks(document)
    fmt.Println(articleLinks[2])
    fmt.Println(getArticle(articleLinks[2]))
}

func getArticleLinks(node *html.Node) []string {
    var links []string

    if node.Type == html.ElementNode && node.Data == "a" {
        for _, attribute := range node.Attr {
            link, err := parseAnchorTag(attribute)
            if err != nil {}
            links = append(links, link)
        }
    }

    for child := node.FirstChild; child != nil; child = child.NextSibling {
        childLinks := getArticleLinks(child)
        links = append(links, childLinks...)
    }

    return links
}

func parseAnchorTag(tag html.Attribute) (string, error) {
    if tag.Key != "href" {
        return "", errors.New("Expected Anchor Tag, but received something else.")
    }

    if strings.Contains(tag.Val, "/") {
        return "", errors.New("Not article link.")
    }
    return tag.Val, nil
}


type Article struct {
    link    string
    title   string
    date    string
}

func getArticle(articleUrl string) Article {
    res, err := http.Get("http://paulgraham.com/" + articleUrl)
    if err != nil {
        log.Println("Could not get article:", err.Error())
    }
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        log.Fatal("res.Status not OK:\n", res.Status)
    }

    articleNode, err := html.Parse(res.Body)
    if err != nil {
        log.Fatal("Couldn't read body: \n", err)
    }

    date, err := getArticleDate(articleNode)
    if err != nil {
        log.Fatal("No date.")
    }

    title, err := getArticleTitle(articleNode)
    if err != nil {
        log.Fatal("No title.")
    }

    return Article{
        link: articleUrl,
        date: date,
        title: title,
    }
}

func getArticleDate(articleNode *html.Node) (string, error) {
    fmt.Println("starting")
    if articleNode.Type == html.ElementNode && articleNode.Data == "font" {
        for node := articleNode.FirstChild; node != nil; node = node.NextSibling {
            if node.Type == html.TextNode {
                fmt.Println(node.Data, "hele")
                return node.Data, nil
            }
        }
    }

    for child := articleNode.FirstChild; child != nil; child = child.NextSibling {
        getArticleDate(child)
    }

    return "", errors.New("Could not get date")
}

func getArticleTitle(articleNode *html.Node) (string, error) {
    if articleNode.Type == html.ElementNode && articleNode.Data == "img" {
        for _, attribute := range articleNode.Attr {
            if attribute.Key == "alt" {
                return attribute.Val, nil
            }
        }
    }

    for child := articleNode.FirstChild; child != nil; child = child.NextSibling {
        getArticleTitle(child)
    }

    return "", errors.New("Could not get title.")
}
