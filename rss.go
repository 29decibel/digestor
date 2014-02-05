package main

import (
	"bytes"
	"fmt"
	rss "github.com/jteeuwen/go-pkg-rss"
	"html/template"
	"time"
)

var rssItems []RSSItem
var timeFormats = []string{
	time.RFC1123Z,
	time.RFC1123,
	time.RFC3339,
	time.RFC822,
	"Mon, 02 Jan 06 15:04:05 -0700",
	"02 Jan 2006"} // must use this date please, otherwise go time won't recognize

// RSSItem is the simple version of rss item
type RSSItem struct {
	Href  string
	Title string
}

// RSSConfig is for rss config information
type RSSConfig struct {
	Feeds []string `json:"feeds"`
}

func rssMarkup() string {
	fmt.Println("Fetching rss feeds.....")

	for _, url := range config.RSS.Feeds {
		pollFeed(url, 5)
	}

	tmpl, err := template.New("rss").Parse(`
    <ul>
    {{range .links}}
        <li><a href="{{.Href}}">{{.Title}}</a></li>
      {{end}}
    </ul>
  `)
	checkErr(err, "template error")
	var results bytes.Buffer

	tmpl.Execute(&results, map[string]interface{}{
		"links": rssItems,
	})
	return string(results.Bytes())
}

func p(any interface{}) {
	fmt.Println(any)
}

func parseTime(t string) time.Time {
	for _, f := range timeFormats {
		tObj, err := time.Parse(f, t)
		if err == nil {
			return tObj
			break
		}
	}

	return time.Time{}
}

func isToday(tString string) bool {
	today := time.Now()
	pTime := parseTime(tString)
	return today.Year() == pTime.Year() &&
		today.Month() == pTime.Month() && today.Day() == pTime.Day()
}

func pollFeed(uri string, timeout int) {
	feed := rss.New(timeout, true, chanHandler, itemHandler)
	feed.Fetch(uri, nil)
}

func chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	/*
	 *fmt.Printf("%d new channel(s) in %s\n", len(newchannels), feed.Url)
	 */
}

func itemHandler(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
	for _, item := range newItems {
		if isToday(item.PubDate) {
			rssItems = append(rssItems, RSSItem{Title: item.Title, Href: item.Links[0].Href})
		}
	}
}
