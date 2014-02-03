package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/css"
	"github.com/mrjones/oauth"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os/user"
)

const (
	mailSubject = "Today's Digestor "
	tweetsCount = "100"
)

var twitterUserWhitelist = []string{"dhh", "brainpicker", "Explore", "CatChen", "Livid"}

// Config is for the whole config
type Config struct {
	Mail    map[string]string `json:"mail"`
	Twitter map[string]string `json:"twitter"`
}

var config *Config

// Tweet is a tweet structure
type Tweet struct {
	CreatedAt string `json:"created_at"`
	Text      string
	User      User `json:"user"`
}

func (tweet Tweet) String() string {
	return tweet.Text
}

// User is the struct for Tweeter user
type User struct {
	Name       string `json:"name"`
	Avatar     string `json:"profile_image_url"`
	ScreenName string `json:"screen_name"`
}

func main() {
	initConfig()

	client, accessToken := auth()

	fmt.Println("Fetching tweets.....")
	response, err := client.Get(
		"https://api.twitter.com/1.1/statuses/home_timeline.json",
		map[string]string{"count": tweetsCount},
		accessToken)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	bits, err := ioutil.ReadAll(response.Body)

	var tweets []Tweet
	json.Unmarshal(bits, &tweets)

	fmt.Println("Sending email...")
	// get the email contents
	contents := emailContents(tweetsMarkup(groupByUser(tweets)))
	sendEmail(contents)
	fmt.Println("Email sent.")
}

func initConfig() {
	usr, _ := user.Current()
	dir := usr.HomeDir

	configFile, err := ioutil.ReadFile(dir + "/.digestor.json")
	checkErr(err, "Can not load config file, please edit your config file: ~/.digestor.json")
	json.Unmarshal(configFile, &config)
}

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

func tweetsMarkup(gTweets map[string][]Tweet) []byte {
	tmpl, err := template.New("tweets").Parse(`
    {{range $name, $tweets := .group}}
      <div class="user-tweets">
        <h4>{{$name}}</h4>
        <ul>
          {{ range $tweets }}
            <li>{{.}}</li>
          {{ end }}
        </ul>
      </div>
    {{end}}
  `)
	checkErr(err, "template error")
	var doc bytes.Buffer

	tmpl.Execute(&doc, map[string]interface{}{
		"group": gTweets,
	})

	return doc.Bytes()
}

func groupByUser(tweets []Tweet) map[string][]Tweet {
	var group = make(map[string][]Tweet)
	for _, tweet := range tweets {
		if stringInSlice(tweet.User.ScreenName, twitterUserWhitelist) {
			group[tweet.User.Name] = append(group[tweet.User.Name], tweet)
		}
	}

	return group
}

// send email to user
func sendEmail(contents []byte) {
	e := email.NewEmail()
	e.From = config.Mail["from"]
	e.To = []string{config.Mail["to"]}
	e.Subject = mailSubject

	e.HTML = contents
	e.Send(config.Mail["host"]+":587", smtp.PlainAuth("", config.Mail["user"], config.Mail["password"], config.Mail["host"]))
}

func emailContents(tweets []byte) []byte {
	// read the templat
	emailTemplateString, err := ioutil.ReadFile("./email.html")
	checkErr(err, "can not load email template ")
	tmpl, err := template.New("tweets").Parse(string(emailTemplateString))

	var doc bytes.Buffer
	tmpl.Execute(&doc, map[string]interface{}{
		"tweetsMarkup":     template.HTML(string(tweets)),
		"githubMarkup":     template.HTML(githubMarkup()),
		"hackerNewsMarkup": template.HTML(hackerNewsMarkup())})

	return doc.Bytes()
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func hackerNewsMarkup() string {

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

func auth() (*oauth.Consumer, *oauth.AccessToken) {
	c := oauth.NewConsumer(
		config.Twitter["consumerKey"],
		config.Twitter["consumerSecret"],
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		})

	if true {
		t := oauth.AccessToken{
			Token:          config.Twitter["accessToken"],
			Secret:         config.Twitter["accessSecret"],
			AdditionalData: map[string]string{"user_id": config.Twitter["user_id"], "screen_name": config.Twitter["screen_name"]}}

		return c, &t
	}

	// Get request token
	requestToken, url, err := c.GetRequestTokenAndUrl("oob")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("(1) Go to: " + url)
	fmt.Println("(2) Grant access, you should get back a verification code.")
	fmt.Println("(3) Enter that verification code here: ")

	verificationCode := ""
	fmt.Scanln(&verificationCode)

	// Get access token
	accessToken, err := c.AuthorizeToken(requestToken, verificationCode)
	if err != nil {
		log.Fatal(err)
	}

	return c, accessToken

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func body(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body
}

/*
 *
 *func getTciwiweets() type {
 *  anaconda.SetConsumerKey("your-consumer-key")
 *  anaconda.SetConsumerSecret("your-consumer-secret")
 *  api := anaconda.NewTwitterApi("your-access-token", "your-access-token-secret")
 *
 *  api.GetUserTimeline
 *
 *}
 */
