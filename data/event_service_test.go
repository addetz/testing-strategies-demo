package data_test

import (
	"testing"

	"github.com/addetz/testing-strategies-demo/data"
	"github.com/stretchr/testify/assert"
)

func TestGetEvent(t *testing.T) {
	events := []data.Event{
		{
			ID: "event-1",
		},
		{
			ID: "event-2",
		},
		{
			ID: "event-3",
		},
	}

	talks := []data.Talk{
		{
			EventID: "event-1",
			Title:   "event 1 talk 1",
		},
		{
			EventID: "event-1",
			Title:   "event 1 talk 2",
		},
		{
			EventID: "event-2",
			Title:   "event 2 talk 1",
		},
		{
			EventID: "event-99",
			Title:   "invalid talk",
		},
	}

	es := data.NewEventService(events, talks)

	t.Run("multiple talks", func(t *testing.T) {
		ev, err := es.GetEvent("event-1")
		assert.Nil(t, err)
		assert.Equal(t, "event-1", ev.ID)
		assert.Len(t, ev.Talks, 2)
		assert.Equal(t, "event 1 talk 1", ev.Talks[0].Title)
		assert.Equal(t, "event 1 talk 2", ev.Talks[1].Title)
	})
	t.Run("single talk", func(t *testing.T) {
		ev, err := es.GetEvent("event-2")
		assert.Nil(t, err)
		assert.Equal(t, "event-2", ev.ID)
		assert.Len(t, ev.Talks, 1)
		assert.Equal(t, "event 2 talk 1", ev.Talks[0].Title)
	})
	t.Run("no talks", func(t *testing.T) {
		ev, err := es.GetEvent("event-3")
		assert.Nil(t, err)
		assert.Len(t, ev.Talks, 0)
	})
	t.Run("invalid id", func(t *testing.T) {
		ev, err := es.GetEvent("event-99")
		assert.Nil(t, ev)
		assert.NotNil(t, err)
		assert.Equal(t, "no event for id event-99", err.Error())
	})
}

func TestGetFilteredTalks(t *testing.T) {
	eventID := "event-1"
	events := []data.Event{
		{
			ID:        eventID,
			DateStart: "01/01/2010",
			DateEnd:   "02/01/2010",
		},
	}

	talks := []data.Talk{
		{
			EventID: eventID,
			Title:   "event 1 talk 1",
			Date:    "01/01/2010",
		},
		{
			EventID: eventID,
			Title:   "event 1 talk 2",
			Date:    "01/01/2010",
		},
		{
			EventID: eventID,
			Title:   "event 1 talk 3",
			Date:    "02/01/2010",
		},
	}

	es := data.NewEventService(events, talks)

	t.Run("multiple talks", func(t *testing.T) {
		talks, err := es.GetFilteredTalks(eventID, 1)
		assert.Nil(t, err)
		assert.Len(t, talks.Talks, 2)
		assert.Equal(t, "event 1 talk 1", talks.Talks[0].Title)
		assert.Equal(t, "event 1 talk 2", talks.Talks[1].Title)
	})

	t.Run("last day", func(t *testing.T) {
		talks, err := es.GetFilteredTalks(eventID, 2)
		assert.Nil(t, err)
		assert.Len(t, talks.Talks, 1)
		assert.Equal(t, "event 1 talk 3", talks.Talks[0].Title)
	})

	t.Run("day after end", func(t *testing.T) {
		expectedErr := "filtered date 03/01/2010 is after event end date 02/01/2010"
		talks, err := es.GetFilteredTalks(eventID, 3)
		assert.Nil(t, talks)
		assert.Equal(t, expectedErr, err.Error())		
	})
}
