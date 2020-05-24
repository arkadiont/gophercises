package link

import (
	"io"
	"log"
	"strings"

	"golang.org/x/net/html"
)

// Link is <a href="">...</a> of HTML document
type Link struct {
	Href string
	Text string
}

// Parse HTML and return
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nodes := findLinkNodes(doc)
	var l []Link
	for _, node := range nodes {
		l = append(l, buildLink(node))
	}
	return l, nil
}

func buildLink(n *html.Node) (l Link) {
	for _, att := range n.Attr {
		if att.Key == "href" {
			l.Href = att.Val
			break
		}
	}
	l.Text = text(n)
	return
}

func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var b strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if _, err := b.WriteString(text(c)); err != nil {
			log.Printf("err: %v\n", err)
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

func findLinkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}

	var l []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		l = append(l, findLinkNodes(c)...)
	}
	return l
}
