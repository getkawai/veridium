package domutils

import (
	"strings"

	"github.com/JohannesKaufmann/dom"
	"golang.org/x/net/html"
)

// GetTitle extracts title from HTML document
// It checks <title>, <meta property="og:title">, and <meta name="twitter:title">
func GetTitle(n *html.Node) string {
	// 1. Try <title> tag first (standard)
	var title string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil {
				title = n.FirstChild.Data
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	title = strings.TrimSpace(title)
	if title != "" {
		return title
	}

	// 2. Try Open Graph title
	if val := getMetaContent(n, "property", "og:title"); val != "" {
		return val
	}

	// 3. Try Twitter title
	if val := getMetaContent(n, "name", "twitter:title"); val != "" {
		return val
	}

	return ""
}

func getMetaContent(n *html.Node, attrKey, attrValue string) string {
	var content string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if content != "" {
			return
		}
		if n.Type == html.ElementNode && n.Data == "meta" {
			if val, _ := dom.GetAttribute(n, attrKey); val == attrValue {
				content, _ = dom.GetAttribute(n, "content")
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return strings.TrimSpace(content)
}
