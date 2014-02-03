package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jordan-wright/email"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os/user"
)

const (
	mailSubject = "Today's Digestor "
)

// Config is for the whole config
type Config struct {
	Mail    map[string]string `json:"mail"`
	Twitter TwitterConfig     `json:"twitter"`
}

var config *Config

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

// send email to user
func sendEmail(contents []byte) {
	e := email.NewEmail()
	e.From = config.Mail["from"]
	e.To = []string{config.Mail["to"]}
	e.Subject = mailSubject

	e.HTML = contents
	e.Send(config.Mail["host"]+":587", smtp.PlainAuth("", config.Mail["user"], config.Mail["password"], config.Mail["host"]))
}

// get mail temaplte from file
// not used for now, cause we have a embeded one
func mailTemplateFromFile() string {
	emailTemplateString, err := ioutil.ReadFile("./email.html")
	checkErr(err, "can not load email template ")
	return string(emailTemplateString)
}

func emailContents(tweets []byte) []byte {
	// read the templat
	tmpl, err := template.New("tweets").Parse(mailTemplate)
	checkErr(err, "mail template create failed")

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
