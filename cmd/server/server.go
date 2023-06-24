package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "embed"

	"github.com/addetz/testing-strategies-demo/data"
	"github.com/addetz/testing-strategies-demo/handlers"
	"github.com/gorilla/mux"
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
	events, talks := importData()
	eventService, err := data.NewEventService(events, talks)
	if err != nil {
		log.Fatal(err)
	}
	handler := handlers.NewHandler(eventService)
	router := configureRouter(handler)

	log.Printf("Server listening on :%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// configureRouter configures the routes of this server and binds handler functions to them
func configureRouter(handler *handlers.Handler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Methods("GET").Path("/events").Handler(http.HandlerFunc(handler.GetEventsHandler))
	router.Methods("GET").Path("/events/{id}").Handler(http.HandlerFunc(handler.GetEventTalksHandler))

	return router
}

func importData() ([]data.Event, []data.Talk) {
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
