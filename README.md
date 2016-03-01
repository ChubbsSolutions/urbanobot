# urbanobot

[![Build Status](https://travis-ci.org/urbanobot/urbanobot.png)](https://travis-ci.org/iarenzana/urbanobot)

Integrates Urban Dictionary right into Slack

##Download the software

[Download](https://github.com/iarenzana/urbanobot/releases) the latest version of urbano for all major platforms.

##Compile and run the source

Requires Go 1.5 or newer (earlier versions untested). Remember to set the GOPATH variable.

```
git clone https://github.com/iarenzana/urbanobot
cd urbanobot
go get
go build
PORT=60000 ./urbano.go
```

##Usage

Run this service in Heroku (Procfile provided). Go to your Custom Integrations, Slash Commands on Slack and create a GET that points to https://[YOUR_HOST]/v1/word.

##About

Crafted with :heart: in Indiana by [Chubbs Solutions] (http://chubbs.solutions).
