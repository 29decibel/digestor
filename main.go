package main

import (
	"bytes"
	"encoding/json"
	"flag"
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
	RSS     RSSConfig         `json:"rss"`
}

var config *Config
var defaultConfigPath = "/.digestor.json"
var configFile = flag.String("c", defaultConfigPath, "config file path")

func main() {
	flag.Parse()

	initConfig()

	// get the email contents
	contents := emailContents()

	sendEmail(contents)
	fmt.Println("Email sent.")
}

func initConfig() {
	var path string
	if *configFile == defaultConfigPath {
		usr, _ := user.Current()
		dir := usr.HomeDir
		path = dir + "/.digestor.json"
	} else {
		path = *configFile
	}

	configFile, err := ioutil.ReadFile(path)
	checkErr(err, "Can not load config file "+path+", please edit your config file: ~/.digestor.json or provide -c your-config.json")
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

func emailContents() []byte {
	// read the templat
	tmpl, err := template.New("tweets").Parse(mailTemplate)
	checkErr(err, "mail template create failed")

	var doc bytes.Buffer
	tmpl.Execute(&doc, map[string]interface{}{
		"tweetsMarkup":     template.HTML(tweetsMarkup()),
		"githubMarkup":     template.HTML(githubMarkup()),
		"rssMarkup":        template.HTML(rssMarkup()),
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
