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

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v1/word", getWord)
	log.Print("Starting up...")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	log.Fatal(http.ListenAndServe(":"+port, router))
}

//GetWord
func getWord(w http.ResponseWriter, r *http.Request) {
	//	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Type", "application/json")

	word := r.URL.Query().Get("text")

	slackUser := r.URL.Query().Get("user_name")
	slackChannel := r.URL.Query().Get("channel_name")
	slackTeam := r.URL.Query().Get("team_id")
	log.Print("Request received for " + word + " from " + slackUser + ", from team " + slackTeam + ", on channel " + slackChannel)
	word = strings.Replace(word, "\"", "", -1)
	wordDefinition, err := getWordDefinition(word)
	if fmt.Sprintf("%s", err) == "NOTFOUND" {
		log.Println("Word " + word + " not found.")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Word not Found.\n"))
		return
	}
	if err != nil {
		log.Println("Error!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Print("Returning definition of " + word + " to " + slackUser + " from team " + slackTeam + " on channel " + slackChannel)
	w.WriteHeader(http.StatusOK)
	response := objects.SlackResponse{}
	response.Text = wordDefinition.Definition
	response.ResponseType = "in_channel"
	//	w.Write([]byte(wordDefinition.Definition))
	resp, err := json.Marshal(response)
	if err != nil {
		log.Println("Error Marshalling response!")
		w.WriteHeader(http.StatusInternalServerError)
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
