package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

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

    var wg sync.WaitGroup
    for _, articleLink := range articleLinks[60:] {
        wg.Add(1)
        getArticle(articleLink, &wg)
    }
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

    /** On PG's website internal links have only their route name
        => <a href="getideas.com">...</a> == paulgraham.com/getideas.html
        Also don't scrape rss.html and index.html
        All other links on page are not inside <a></a>; they are inside <area></area> represting a sidebar.
    */
    if strings.Contains(tag.Val, "/") || tag.Val == "rss.html" || tag.Val == "index.html" {
        return "", errors.New("Not article link.")
    }
    return tag.Val, nil
}


type Article struct {
    link    string
    title   string
    date    string
}

func getArticle(articleUrl string, wg *sync.WaitGroup) (Article, error) {
    defer wg.Done()
    if articleUrl == "" {
        return Article{}, errors.New("Empty articleUrl")
    }

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
        log.Println("Couldn't read body: \n", articleUrl, "\n", err)
        return Article{}, err
    }

    date, err := getArticleDate(articleNode)
    if err != nil {
        log.Println("No date: " + articleUrl)
        return Article{}, err
    }

    title, err := getArticleTitle(articleNode)
    if err != nil {
        log.Println("No title." + articleUrl)
        return Article{}, err
    }

    article := Article{
        link: articleUrl,
        date: date,
        title: title,
    }

    fmt.Println(article)
    return article, nil
}


func getArticleDate(articleNode *html.Node) (string, error) {
    if articleNode.Type == html.ElementNode && (articleNode.Data == "font" || articleNode.Data == "p") {
        for node := articleNode.FirstChild; node != nil; node = node.NextSibling {
            if node.Type == html.TextNode {
                if matched, _ := regexp.MatchString(`^\w+ \d{4}$`, node.Data); matched {
                    return node.Data, nil
                } else if matched, _ := regexp.MatchString(`^\n\w+ \d{4}$`, node.Data); matched {
                    return node.Data[1:], nil
                } else {
                    break
                }
            }
        }
    }

    for child := articleNode.FirstChild; child != nil; child = child.NextSibling {
        date, err := getArticleDate(child)
        if err == nil {
            return date, nil
        }
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
        title, err := getArticleTitle(child)
        if err == nil {
            return title, err
        }
    }

    return "", errors.New("Could not get title.")
}
