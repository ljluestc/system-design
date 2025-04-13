package extractor

import (
    "strings"
    "golang.org/x/net/html"
    "crawler/internal/utils"
)

func ParseHTML(content []byte) ([]string, string, error) {
    urls := []string{}
    var text strings.Builder

    node, err := html.Parse(strings.NewReader(string(content)))
    if err != nil {
        return nil, "", utils.Errorf("parse HTML: %v", err)
    }

    var walk func(*html.Node)
    walk = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "a" {
            for _, attr := range n.Attr {
                if attr.Key == "href" {
                    urls = append(urls, attr.Val)
                }
            }
        }
        if n.Type == html.TextNode {
            text.WriteString(n.Data)
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            walk(c)
        }
    }
    walk(node)

    return urls, text.String(), nil
}