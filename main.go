package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ajg/form"
	"github.com/gorilla/mux"
	"gitlab.com/iarenzana/urbanobot/objects"
	"golang.org/x/crypto/acme/autocert"
)

const version = "1.2.2"

func main() {

	useTLS := flag.Bool("https", false, "use https by default.")
	usePort := flag.Int("port", 61000, "Port to use. Ignored if TLS is enabled.")
	flag.Parse()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/urbano/v1/word", getWord)
	router.HandleFunc("/urbano/v1/random", getRandomWord)

	if *useTLS {
		//Check for the domain
		domain := os.Getenv("URBANO_DOMAIN")
		if domain == "" {
			log.Fatal("$URBANO_DOMAIN must be set")
		}
		log.Printf("Starting up urbanobot %v on %v using https...\n", version, domain)
		//Get certificate and store it under /usr/local/etc. Auto-renewed.
		certManager := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domain),
			Cache:      autocert.DirCache("/usr/local/etc/urbanobot/urbanocerts"),
		}

		//Start server on the https port
		server := &http.Server{
			Addr:    ":https",
			Handler: router,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}
		log.Fatal(server.ListenAndServeTLS("", ""))

	} else {
		log.Printf("Starting up urbanobot %v on port %v...\n", version, *usePort)

		log.Fatal(http.ListenAndServe(":"+fmt.Sprintf("%v", *usePort), router))
	}
}

//GetWord
func getWord(w http.ResponseWriter, r *http.Request) {

	var word, slackUser, slackChannel, slackTeam string

	if strings.Contains(r.Header.Get("User-Agent"), "Slackbot") {
		w.Header().Set("Content-Type", "application/json")
		word = r.URL.Query().Get("text")

		slackUser = r.URL.Query().Get("user_name")
		slackChannel = r.URL.Query().Get("channel_name")
		slackTeam = r.URL.Query().Get("team_id")

		log.Print("Slack request received for " + word + " from " + slackUser + ", from team " + slackTeam + ", on channel " + slackChannel)
	} else {
		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		var u objects.SlackIncoming

		d := form.NewDecoder(r.Body)
		if err2 := d.Decode(&u); err2 != nil {
			fmt.Printf("Form could not be decoded - %v", err2)
			return
		}

		word = u.Text
		slackUser = u.SlackUser
		slackChannel = u.SlackChannel
		slackTeam = u.SlackTeam

		log.Print("Other request received for " + word + " from " + slackUser + ", from team " + slackTeam + ", on channel " + slackChannel)
	}

	if word == "" {
		response := objects.SlackResponse{}
		response.Text = "Screw you, @barnes. Happy now?"
		response.ResponseType = "in_channel"
		response.BotVersion = version

		resp, err := json.Marshal(response)
		if err != nil {
			log.Println("Error Marshalling response!")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(resp)
			return

		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
		return

	}
	wordDefinition, err := getWordDefinition(word)
	if fmt.Sprintf("%s", err) == "NOTFOUND" {
		log.Println("Word " + word + " not found.")
		w.WriteHeader(http.StatusOK)

		response := objects.SlackResponse{}
		response.Text = fmt.Sprintf("%s - Word not found", word)
		response.ResponseType = "ephemeral"
		response.BotVersion = version

		resp, err := json.Marshal(response)
		if err != nil {
			log.Println("Error Marshalling response!")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(resp)
		return
	}

	if err != nil {
		log.Print("Error ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Print("Returning definition of " + word + " to " + slackUser + " from team " + slackTeam + " on channel " + slackChannel)

	w.WriteHeader(http.StatusOK)
	response := objects.SlackResponse{}
	response.Text = fmt.Sprintf("%s --> %s", word, wordDefinition.Definition)
	response.ResponseType = "in_channel"
	response.BotVersion = version

	resp, err := json.Marshal(response)
	if err != nil {
		resp, err := json.Marshal(response)
		if err != nil {
			log.Println("Error Marshalling response!")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(resp)
		return

	}
	w.Write(resp)
}

func getWordDefinition(wordToDefine string) (objects.WordData, error) {
	var UDURL = "http://api.urbandictionary.com/v0/define?term=" + strings.Replace(wordToDefine, " ", "", -1)
	wd := objects.WordDataSlice{}
	var word objects.WordData

	resp, err := http.Get(UDURL)
	if err != nil {
		return word, err
	}

	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return word, err
	}

	err = json.Unmarshal([]byte(string(data)), &wd)
	if err != nil {
		return word, err
	}

	for _, element := range wd.List {
		if element.ThumbsUp > word.ThumbsUp {
			word = element
		}
	}
	if word.Definition == "" {
		return word, errors.New("NOTFOUND")
	}
	return word, nil
}

//GetWord
func getRandomWord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slackUser := r.URL.Query().Get("user_name")
	slackChannel := r.URL.Query().Get("channel_name")
	slackTeam := r.URL.Query().Get("team_id")
	log.Print("Request for random word received from " + slackUser + ", from team " + slackTeam + ", on channel " + slackChannel)

	wordDefinition, err := getNewWord()
	if err != nil {
		log.Print("Error ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Print("Returning random to " + slackUser + " from team " + slackTeam + " on channel " + slackChannel)

	w.WriteHeader(http.StatusOK)
	response := objects.SlackResponse{}
	response.Text = fmt.Sprintf("%s --> %s", wordDefinition.Word, wordDefinition.Definition)
	response.ResponseType = "in_channel"
	response.BotVersion = version

	resp, err := json.Marshal(response)
	if err != nil {
		log.Println("Error Marshalling response!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

//getNewWord gets a random UD word
func getNewWord() (objects.WordData, error) {
	var UDURL = "http://api.urbandictionary.com/v0/random"
	wd := objects.WordDataSlice{}
	var word objects.WordData
	var good = false
	tu := 13000

	for good == false {
		resp, err := http.Get(UDURL)
		if err != nil {
			return word, err
		}
		defer resp.Body.Close()

		data, _ := ioutil.ReadAll(resp.Body)
		if err != nil {
			return word, err
		}

		err = json.Unmarshal([]byte(string(data)), &wd)
		if err != nil {
			return word, err
		}

		for _, element := range wd.List {
			if element.ThumbsUp > tu {
				word = element
				good = true
			}
		}
	}
	return word, nil
}
