package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	link "gophercises/ex4-link"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const xmlNs = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	XMLNs string `xml:"xmlns,attr"`
}

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the url that you want to build a sitemap for")
	maxDepth := flag.Int("depth", 3, "the maximum number of links deep to traverse")
	flag.Parse()

	pages := bfs(*urlFlag, *maxDepth)

	toXML := urlset{
		XMLNs: xmlNs,
	}
	for _, r := range pages {
		toXML.Urls = append(toXML.Urls, loc{r})
	}

	fmt.Print(xml.Header)
	xmlEnc := xml.NewEncoder(os.Stdout)
	xmlEnc.Indent("", "\t")
	if err := xmlEnc.Encode(toXML); err != nil {
		log.Fatal(err)
	}
}

func get(uStr string) []string {
	resp, err := http.Get(uStr)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Host:   reqURL.Host,
		Scheme: reqURL.Scheme,
	}
	base := baseURL.String()
	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func bfs(urlStr string, maxDepth int) (ret []string) {
	visited := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: struct{}{},
	}
	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})
		for url := range q {
			if _, ok := visited[url]; ok {
				continue
			}
			visited[url] = struct{}{}
			for _, link := range get(url) {
				if _, ok := visited[link]; !ok {
					nq[link] = struct{}{}
				}
			}
		}
	}

	for url := range visited {
		ret = append(ret, url)
	}
	return
}

func hrefs(r io.Reader, base string) (ret []string) {
	links, _ := link.Parse(r)
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}
	return
}

func filter(links []string, keepFn func(string) bool) (ret []string) {
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}
	return
}

func withPrefix(prefix string) func(string) bool {
	return func(s string) bool {
		return strings.HasPrefix(s, prefix)
	}
}
