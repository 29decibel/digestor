package main

import (
	"bytes"
	"fmt"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/css"
	"github.com/moovweb/gokogiri/xml"
	"html/template"
)

// HNLink is a struct for a single link
type HNLink struct {
	LinkMarkup template.HTML
	Excerpt    template.HTML
}

// HNConfig is for hackernews config
type HNConfig struct {
	ReadabilityKey string `json:"readability_key"`
}

func (link HNLink) String() string {
	return fmt.Sprintf("%s\n%s", link.LinkMarkup, link.Excerpt)
}

func hackerNewsMarkup() string {
	fmt.Println("Fetching hackernews.....")

	body := body("https://news.ycombinator.com/")
	doc, _ := gokogiri.ParseHtml(body)

	xpath := css.Convert("td.title a", css.GLOBAL)
	links, _ := doc.Search(xpath)

	var newLinks []HNLink
	var chanLinks = make(chan HNLink)
	for _, l := range links[:len(links)-1] {
		go parseLink(l, chanLinks)
	}

	// collect back links info
	for i := 0; i < len(links)-1; i++ {
		newLinks = append(newLinks, <-chanLinks)
	}

	tmpl, err := template.New("hackernews").Parse(`
    <ul>
      {{range .links}}
        <li>{{.LinkMarkup}}</li>
        <p>{{.Excerpt}}</p>
      {{end}}
    </ul>
  `)
	checkErr(err, "template error")
	var results bytes.Buffer

	tmpl.Execute(&results, map[string]interface{}{
		"links": newLinks,
	})

	return string(results.Bytes())
}

func parseLink(node xml.Node, linksChan chan HNLink) {
	hnlink := HNLink{
		LinkMarkup: template.HTML(node.String()),
		Excerpt:    template.HTML(excerpt(node.Attribute("href").Value()))}
	linksChan <- hnlink
}

func excerpt(url string) string {
	article := Article{URL: url}
	article.parse()

	return article.Excerpt
}
