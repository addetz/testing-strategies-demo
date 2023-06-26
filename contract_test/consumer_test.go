package contract_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/addetz/testing-strategies-demo/data"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsumerIndex_Local(t *testing.T) {
	if os.Getenv("CONTRACT") == "" {
		t.Skip("Skipping TestConsumerIndex_Local in short mode.")
	}
	expectedEvents := []data.Event{
		{
			ID:        "ewit-2023",
			Name:      "European Women in Tech",
			DateStart: "28/06/2023",
			DateEnd:   "29/06/2023",
			Location:  "Amsterdam",
		},
		{
			ID:        "devbcn-2023",
			Name:      "DevBcn - The Barcelona Developers Conference",
			DateStart: "03/07/2023",
			DateEnd:   "05/07/2023",
			Location:  "Barcelona",
		},
		{
			ID:        "cphdevfest-2023",
			Name:      "Copenhagen Developers Festival",
			DateStart: "30/08/2023",
			DateEnd:   "01/09/2023",
			Location:  "Copenhagen",
		},
	}

	// Initialize
	pact := dsl.Pact{
		Consumer: "Consumer",
		Provider: "ConfTalks",
	}
	pact.Setup(true)

	// Test case - makes the call to the provider
	var test = func() (err error) {
		url := fmt.Sprintf("%s:%d/events", "http://localhost", pact.Server.Port)
		req, err := http.NewRequest("GET", url, nil)
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		return
	}

	t.Run("get events", func(t *testing.T) {
		pact.
			AddInteraction().
			Given("ConfTalksServer is up").
			UponReceiving("GET /events request").
			WithRequest(dsl.Request{
				Method: "GET",
				Path:   dsl.String("/events"),
				Headers: dsl.MapMatcher{
					"Content-Type": dsl.String("application/json"),
				},
			}).
			WillRespondWith(dsl.Response{
				Status: http.StatusOK,
				Body: dsl.Like(data.Events{
					Events: expectedEvents,
				}),
			})
		require.Nil(t, pact.Verify(test))
	})

	// Clean up
	require.Nil(t, pact.WritePact())
	pact.Teardown()
}
