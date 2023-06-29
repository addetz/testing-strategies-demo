package handlers_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/addetz/testing-strategies-demo/data"
	"github.com/addetz/testing-strategies-demo/handlers"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEventsIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("Skipping TestGetEventsIntegration in short mode.")
	}

	events := []data.Event{
		{
			ID:   "event-1",
			Name: "Event 1 2023",
		},
		{
			ID:   "event-2",
			Name: "Event 2 2023",
		},
	}

	// Arrange
	eventsService, err := data.NewEventService(events, []data.Talk{})
	assert.Nil(t, err)
	ha := handlers.NewHandler(eventsService)
	svr := httptest.NewServer(http.HandlerFunc(ha.GetEventsHandler))
	defer svr.Close()

	// Act
	r, err := http.Get(svr.URL)

	// Assert
	require.Nil(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	require.Nil(t, err)

	var resp data.Events
	err = json.Unmarshal(body, &resp)
	require.Nil(t, err)
	assert.Len(t, resp.Events, len(events))
	for _, e := range events {
		assert.Contains(t, resp.Events, e)
	}
}

func TestGetEventIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("Skipping TestGetEventIntegration in short mode.")
	}
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
		},
	}
	es, err := data.NewEventService(events, talks)
	assert.Nil(t, err)
	assert.NotNil(t, es)

	testCases := map[string]struct {
		eventID            string
		expectedTalks      []data.Talk
		expectedErr        string
		expectedStatusCode int
	}{
		"multiple talks": {
			eventID:            eventID,
			expectedTalks:      talks[0:2],
			expectedStatusCode: http.StatusOK,
		},
		"empty talks": {
			eventID:            "event-2",
			expectedTalks:      []data.Talk{},
			expectedStatusCode: http.StatusOK,
		},
		"invalid event": {
			eventID:            "invalid-event",
			expectedErr:        "no event for id invalid-event",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	// Arrange
	ha := handlers.NewHandler(es)
	router := mux.NewRouter()
	router.HandleFunc("/events/{id}", ha.GetEventTalksHandler)

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			path := fmt.Sprintf("/events/%s", tc.eventID)
			req, err := http.NewRequest("GET", path, nil)
			require.Nil(t, err)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			require.Equal(t, tc.expectedStatusCode, rr.Code)

			if len(tc.expectedErr) != 0 {
				var respErr handlers.ErrorResponse
				bytes := rr.Body.Bytes()
				err = json.Unmarshal(bytes, &respErr)
				require.Nil(t, err)
				assert.Contains(t, respErr.Error, tc.expectedErr)
				return
			}

			var resp data.Talks
			err = json.Unmarshal(rr.Body.Bytes(), &resp)
			require.Nil(t, err)
			assert.Len(t, resp.Talks, len(tc.expectedTalks))
			for i, expectedTalk := range tc.expectedTalks {
				assert.Equal(t, expectedTalk, resp.Talks[i])
			}
		})
	}
}

func TestGetEventFilteredTalksIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("Skipping TestGetEventFilteredTalksIntegration in short mode.")
	}
	eventID := "event-1"
	events := []data.Event{
		{
			ID:        eventID,
			DateStart: "01/02/2010",
			DateEnd:   "02/02/2010",
		},
		{
			ID: "event-2",
		},
	}

	talks := []data.Talk{
		{
			EventID: eventID,
			Title:   "event 1 talk 1",
			Date:    "01/02/2010",
		},
		{
			EventID: eventID,
			Title:   "event 1 talk 2",
			Date:    "01/02/2010",
		},
		{
			EventID: eventID,
			Title:   "event 1 talk 3",
			Date:    "02/02/2010",
		},
	}
	es, err := data.NewEventService(events, talks)
	assert.Nil(t, err)
	assert.NotNil(t, es)

	testCases := map[string]struct {
		eventID            string
		day                string
		expectedTalks      []data.Talk
		expectedErr        string
		expectedStatusCode int
	}{
		"multiple talks": {
			eventID:            eventID,
			day:                "1",
			expectedTalks:      talks[0:2],
			expectedStatusCode: http.StatusOK,
		},
		"single talk": {
			eventID:            eventID,
			day:                "2",
			expectedTalks:      []data.Talk{talks[2]},
			expectedStatusCode: http.StatusOK,
		},
		"empty talks": {
			eventID:            "event-2",
			expectedTalks:      []data.Talk{},
			expectedStatusCode: http.StatusOK,
		},
		"invalid event": {
			eventID:            "invalid-event",
			expectedErr:        "no event for id invalid-event",
			expectedStatusCode: http.StatusBadRequest,
		},
		"invalid day param": {
			eventID:            eventID,
			day:                "adelina",
			expectedErr:        "strconv.Atoi",
			expectedStatusCode: http.StatusBadRequest,
		},
		"negative day param": {
			eventID:            eventID,
			day:                "-1",
			expectedErr:        "day must be > 1, but was -1",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	// Arrange
	ha := handlers.NewHandler(es)
	router := mux.NewRouter()
	router.HandleFunc("/events/{id}", ha.GetEventTalksHandler)

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			path := fmt.Sprintf("/events/%s?day=%s", tc.eventID, tc.day)
			req, err := http.NewRequest("GET", path, nil)
			require.Nil(t, err)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			require.Equal(t, tc.expectedStatusCode, rr.Code)

			if len(tc.expectedErr) != 0 {
				var respErr handlers.ErrorResponse
				bytes := rr.Body.Bytes()
				err = json.Unmarshal(bytes, &respErr)
				require.Nil(t, err)
				assert.Contains(t, respErr.Error, tc.expectedErr)
				return
			}

			var resp data.Talks
			err = json.Unmarshal(rr.Body.Bytes(), &resp)
			require.Nil(t, err)
			assert.Len(t, resp.Talks, len(tc.expectedTalks))
			for i, expectedTalk := range tc.expectedTalks {
				assert.Equal(t, expectedTalk, resp.Talks[i])
			}
		})
	}
}
