package main

import (
	"bytes"
	"fmt"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/css"
	"html/template"
)

func hackerNewsMarkup() string {
	fmt.Println("Fetching hackernews.....")

	body := body("https://news.ycombinator.com/")
	doc, _ := gokogiri.ParseHtml(body)

	xpath := css.Convert("td.title a", css.GLOBAL)
	links, _ := doc.Search(xpath)

	var newLinks []interface{}
	for _, l := range links[:len(links)-1] {
		newLinks = append(newLinks, template.HTML(l.String()))
	}

	tmpl, err := template.New("hackernews").Parse(`
    <ul>
      {{range .links}}
        <li>{{.}}</li>
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
