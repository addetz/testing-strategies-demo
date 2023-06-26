package data_test

import (
	"errors"
	"testing"

	"github.com/addetz/testing-strategies-demo/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEventService(t *testing.T) {
	events := []data.Event{
		{
			ID: "event-1",
		},
	}
	talks := []data.Talk{
		{
			EventID: "event-1",
			Title:   "event 1 talk 1",
		},
	}
	t.Run("successful initialisation", func(t *testing.T) {
		es, err := data.NewEventService(events, talks)
		assert.Nil(t, err)
		assert.NotNil(t, es)
	})

	t.Run("nil events", func(t *testing.T) {
		es, err := data.NewEventService(nil, talks)
		assert.Nil(t, es)
		assert.Equal(t, data.ErrEventServiceInitialisation, err)
	})

	t.Run("nil talks", func(t *testing.T) {
		es, err := data.NewEventService(events, nil)
		assert.Nil(t, es)
		assert.Equal(t, data.ErrEventServiceInitialisation, err)
	})
}

func TestGetEvents(t *testing.T) {
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
	es, err := data.NewEventService(events, []data.Talk{})
	assert.Nil(t, err)
	assert.NotNil(t, es)

	t.Run("get events", func(t *testing.T) {
		fetched := es.GetEvents()
		assert.Len(t, fetched.Events, len(events))
		for _, e := range events {
			assert.Contains(t, fetched.Events, e)
		}
	})
}

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

	es, err := data.NewEventService(events, talks)
	assert.Nil(t, err)
	assert.NotNil(t, es)

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

func FuzzGetEvent(f *testing.F) {
	eventID := "event-1"
	events := []data.Event{
		{
			ID: eventID,
		},
	}
	f.Add("event 1 talk1")
	f.Add("Comprehensive testing strategies for modern microservice architectures")
	f.Fuzz(func(t *testing.T, name string) {
		talks := []data.Talk{
			{
				EventID: eventID,
				Title:   name,
			},
		}
		es, err := data.NewEventService(events, talks)
		assert.Nil(t, err)
		assert.NotNil(t, es)
		ev, err := es.GetEvent(eventID)
		require.Nil(t, err)
		assert.Len(t, ev.Talks, 1)
		assert.Equal(t, name, talks[0].Title)
	})
}

func TestGetEventTalks(t *testing.T) {
	eventID := "event-1"
	events := []data.Event{
		{
			ID: eventID,
		},
		{
			ID: "event-2",
		},
	}

	talks := []data.Talk{
		{
			EventID: eventID,
			Title:   "event 1 talk 1",
		},
		{
			EventID: eventID,
			Title:   "event 1 talk 2",
			Date:    "01/01/2010",
		},
	}

	es, err := data.NewEventService(events, talks)
	assert.Nil(t, err)
	assert.NotNil(t, es)

	testCases := map[string]struct {
		eventID       string
		expectedTalks []data.Talk
		expectedErr   error
	}{
		"multiple talks": {
			eventID:       eventID,
			expectedTalks: talks[0:2],
		},
		"empty talks": {
			eventID:       "event-2",
			expectedTalks: []data.Talk{},
		},
		"invalid event": {
			eventID:     "invalid-event",
			expectedErr: errors.New("no event for id invalid-event"),
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			talks, err := es.GetEventTalks(tc.eventID)
			if tc.expectedErr != nil {
				assert.Nil(t, talks)
				assert.Equal(t, tc.expectedErr, err)
				return
			}
			assert.Nil(t, err)
			assert.Len(t, talks.Talks, len(tc.expectedTalks))
			for i, expectedTalk := range tc.expectedTalks {
				assert.Equal(t, expectedTalk, talks.Talks[i])
			}
		})
	}
}
func TestGetEventFilteredTalks(t *testing.T) {
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

	es, err := data.NewEventService(events, talks)
	assert.Nil(t, err)
	assert.NotNil(t, es)

	testCases := map[string]struct {
		eventID       string
		day           int
		expectedTalks []data.Talk
		expectedErr   error
	}{
		"multiple talks": {
			eventID:       eventID,
			day:           1,
			expectedTalks: talks[0:2],
		},
		"last day": {
			eventID:       eventID,
			day:           2,
			expectedTalks: []data.Talk{talks[2]},
		},
		"day after end": {
			eventID:     eventID,
			day:         3,
			expectedErr: errors.New("filtered date 03/01/2010 is after event end date 02/01/2010"),
		},
		"invalid event": {
			eventID:     "invalid-event",
			day:         1,
			expectedErr: errors.New("no event for id invalid-event"),
		},
		"negative day": {
			eventID:     eventID,
			day:         -1,
			expectedErr: errors.New("day must be > 1, but was -1"),
		},
		"zero day": {
			eventID:     eventID,
			day:         0,
			expectedErr: errors.New("day must be > 1, but was 0"),
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			talks, err := es.GetEventFilteredTalks(tc.eventID, tc.day)
			if tc.expectedErr != nil {
				assert.Nil(t, talks)
				assert.Equal(t, tc.expectedErr, err)
				return
			}
			assert.Nil(t, err)
			assert.Len(t, talks.Talks, len(tc.expectedTalks))
			for i, expectedTalk := range tc.expectedTalks {
				assert.Equal(t, expectedTalk, talks.Talks[i])
			}
		})
	}
}
