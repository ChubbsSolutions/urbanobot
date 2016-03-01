package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/iarenzana/urbanobot/objects"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v1/word", getWord)
	log.Print("Starting up...")
	log.Fatal(http.ListenAndServe(":60000", router))

}

//GetWord
func getWord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Println("Ping!")

	word := r.URL.Query().Get("text")

	slackUser := r.URL.Query().Get("user_name")
	slackChannel := r.URL.Query().Get("channel_name")
	slackTeam := r.URL.Query().Get("team_id")
	log.Print("Request received from " + slackUser + ", from team " + slackTeam + ", on channel " + slackChannel)
	wordDefinition, err := getWordDefinition(word)
	if fmt.Sprintf("%s", err) == "NOTFOUND" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Word not Found.\n"))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Print("Returning definition of " + word + " to " + slackUser + " from team " + slackTeam + " on channel " + slackChannel)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(wordDefinition.Definition))
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
