package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/iarenzana/urbanobot/objects"
)

const version = "0.5"

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v1/word", getWord)
	router.HandleFunc("/v1/random", getRandomWord)
	log.Printf("Starting up urbanobot %v...\n", version)
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	log.Fatal(http.ListenAndServe(":"+port, router))

	log.Print("Server started on port " + port)
}

//GetWord
func getWord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	word := r.URL.Query().Get("text")

	slackUser := r.URL.Query().Get("user_name")
	slackChannel := r.URL.Query().Get("channel_name")
	slackTeam := r.URL.Query().Get("team_id")
	log.Print("Request received for " + word + " from " + slackUser + ", from team " + slackTeam + ", on channel " + slackChannel)

	wordDefinition, err := getWordDefinition(word)
	if fmt.Sprintf("%s", err) == "NOTFOUND" {
		log.Println("Word " + word + " not found.")
		w.WriteHeader(http.StatusNotFound)

		response := objects.Response{Response: "Word not Found", BotVersion: version}
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
		log.Print("Error - %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Print("Returning definition of " + word + " to " + slackUser + " from team " + slackTeam + " on channel " + slackChannel)

	w.WriteHeader(http.StatusOK)
	response := objects.SlackResponse{}
	response.Text = wordDefinition.Definition
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
		log.Print("Error - %v\n", err)
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
