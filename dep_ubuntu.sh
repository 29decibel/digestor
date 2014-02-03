#!/bin/sh

# dep of rubex
sudo apt-get install -y libonig-dev

# go libs
go get github.com/moovweb/gokogiri
go get github.com/jordan-wright/email
go get github.com/mrjones/oauth
go get github.com/moovweb/rubex
