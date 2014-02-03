package main

import (
	"bytes"
	"fmt"
	"github.com/mrjones/oauth"
	"html/template"
	"log"
)

const (
	tweetsCount = "100"
)

// TwitterConfig is for twitter
type TwitterConfig struct {
	UserID         string   `json:"user_id"`
	WhiteListUsers []string `json:"white_list_users"`
	ScreenName     string   `json:"screen_name"`
	ConsumerKey    string   `json:"consumerKey"`
	ConsumerSecret string   `json:"consumerSecret"`
	AccessToken    string   `json:"accessToken"`
	AccessSecret   string   `json:"accessSecret"`
}

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
		if stringInSlice(tweet.User.ScreenName, config.Twitter.WhiteListUsers) {
			group[tweet.User.Name] = append(group[tweet.User.Name], tweet)
		}
	}

	return group
}

func auth() (*oauth.Consumer, *oauth.AccessToken) {
	c := oauth.NewConsumer(
		config.Twitter.ConsumerKey,
		config.Twitter.ConsumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		})

	if true {
		t := oauth.AccessToken{
			Token:  config.Twitter.AccessToken,
			Secret: config.Twitter.AccessSecret,
			AdditionalData: map[string]string{
				"user_id":     config.Twitter.UserID,
				"screen_name": config.Twitter.ScreenName}}

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
