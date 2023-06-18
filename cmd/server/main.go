package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "embed"

	"github.com/addetz/testing-strategies-demo/data"
)

//go:embed events.json
var eventsFile []byte

//go:embed talks.json
var talksFile []byte

func main() {
	log.Println("Initializing Conference Talks Server ... ")
	port := "8000"
	if p := os.Getenv("SERVER_PORT"); p != "" {
		port = p
	}
	events, talks := importInitial()
	eventService := data.NewEventService(events, talks)
	log.Println(eventService)

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	})
	log.Printf("Server listening on :%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func importInitial() ([]data.Event, []data.Talk) {
	var events data.Events
	var talks data.Talks

	err := json.Unmarshal(eventsFile, &events)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(talksFile, &talks)
	if err != nil {
		log.Fatal(err)
	}

	return events.Events, talks.Talks
}
