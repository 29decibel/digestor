package main

import (
	"bytes"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/css"
	"html/template"
)

func githubMarkup() string {
	body := body("https://github.com/trending")
	doc, _ := gokogiri.ParseHtml(body)

	xpath := css.Convert("div.leaderboard-list-content", css.GLOBAL)
	links, _ := doc.Search(xpath)

	var newLinks []interface{}
	for _, l := range links {
		result := ""

		repos, _ := l.Search(css.Convert(".repository-name", css.LOCAL))
		desc, _ := l.Search(css.Convert(".repo-leaderboard-description", css.LOCAL))

		if len(repos) > 0 {
			result += repos[0].String()
		}
		if len(desc) > 0 {
			result += desc[0].String()
		}
		newLinks = append(newLinks, template.HTML(result))
	}

	tmpl, err := template.New("github-trending").Parse(`
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
