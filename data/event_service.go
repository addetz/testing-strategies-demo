package data

import (
	"fmt"
	"log"
	"time"
)

const dateFormat = "02/01/2006"

type EventService struct {
	// uuid is key to events map
	events map[string]Event
}

func NewEventService(ev []Event, talks []Talk) *EventService {
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

	return es
}

func (es *EventService) GetEvent(id string) (*Event, error) {
	event, ok := es.events[id]
	if !ok {
		return nil, fmt.Errorf("no event for id %s", id)
	}

	return &event, nil
}

func (es *EventService) GetFilteredTalks(id string, day int) (*Talks, error) {
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
			log.Println(t)
			filteredTalks.Talks = append(filteredTalks.Talks, t)
		}
	}

	return filteredTalks, nil

}
