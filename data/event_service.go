package data

import (
	"errors"
	"fmt"
	"log"
	"time"
)

const dateFormat = "02/01/2006"

var ErrEventServiceInitialisation = errors.New("cannot initialise event service with nil events or talks")

type EventService struct {
	// uuid is key to events map
	events map[string]Event
}

// NewEventService initialises and returns and instance to EventService give, slices of events and talks,
// or an error if either events or talks are nil.
func NewEventService(ev []Event, talks []Talk) (*EventService, error) {
	if ev == nil || talks == nil {
		return nil, ErrEventServiceInitialisation
	}
	es := &EventService{
		events: make(map[string]Event),
	}
	for _, e := range ev {
		es.events[e.ID] = e
	}
	for _, t := range talks {
		event, ok := es.events[t.EventID]
		if !ok {
			log.Printf("key %s not found; dropping invalid talk\n", t.EventID)
			continue
		}
		event.Talks = append(event.Talks, t)
		es.events[t.EventID] = event
	}

	return es, nil
}

// GetEvents returns the full list of events.
func (es *EventService) GetEvents() Events {
	var events []Event
	for _, ev := range es.events {
		events = append(events, ev)
	}

	return Events{
		Events: events,
	}
}

// GetEvents returns the event corresponding to the given ID,
// or an error if no event is found.
func (es *EventService) GetEvent(id string) (*Event, error) {
	event, ok := es.events[id]
	if !ok {
		return nil, fmt.Errorf("no event for id %s", id)
	}

	return &event, nil
}

// GetEventTalks returns all the talks of the event corresponding to the given id,
// or an error if no event is found.
func (es *EventService) GetEventTalks(id string) (*Talks, error) {
	event, err := es.GetEvent(id)
	if err != nil {
		return nil, err
	}
	return &Talks{
		Talks: event.Talks,
	}, nil
}

// GetEventFilteredTalks returns all the talks of the event corresponding to the given id and day count,
// or an error if no event is found.
func (es *EventService) GetEventFilteredTalks(id string, day int) (*Talks, error) {
	if day < 1 {
		return nil, fmt.Errorf("day must be > 1, but was %d", day)
	}
	event, err := es.GetEvent(id)
	if err != nil {
		return nil, err
	}
	startDate, err := time.Parse(dateFormat, event.DateStart)
	if err != nil {
		return nil, err
	}
	endDate, err := time.Parse(dateFormat, event.DateEnd)
	if err != nil {
		return nil, err
	}

	// minus 1 to count start date as day 1
	filteredDate := startDate.Add(time.Hour * 24 * time.Duration(day-1))
	searchDate := filteredDate.Format(dateFormat)
	if filteredDate.After(endDate) {
		return nil, fmt.Errorf("filtered date %v is after event end date %v", searchDate, event.DateEnd)
	}

	filteredTalks := &Talks{}
	for _, t := range event.Talks {
		if t.Date == searchDate {
			filteredTalks.Talks = append(filteredTalks.Talks, t)
		}
	}

	return filteredTalks, nil
}
