package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	api = "http://readability.com/api/content/v1/parser?token=%s&url=%s"
)

// Article is the representation of single resource
type Article struct {
	URL     string
	Content string `json: "content"`
	Excerpt string `json: "excerpt"`
	Title   string `json: "title"`
}

// parse the url
func (article *Article) parse() {
	response, _ := http.Get(fmt.Sprintf(api, config.Hackernews.ReadabilityKey, article.URL))
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, article)
}
