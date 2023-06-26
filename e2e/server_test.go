package e2e_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/addetz/testing-strategies-demo/data"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestServerEvents(t *testing.T) {
	if os.Getenv("E2E") == "" {
		t.Skip("Skipping TestServerEvents in short mode.")
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
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "classicaddetz/conf-talks-server:latest",
		ExposedPorts: []string{"8000/tcp"},
		WaitingFor:   wait.ForListeningPort(nat.Port("8000")),
	}
	serverC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.Nil(t, err)
	defer func() {
		if err := serverC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err.Error())
		}
	}()
	url, err := serverC.Endpoint(ctx, "http")
	require.Nil(t, err)
	r, err := http.Get(url + "/events")
	require.Nil(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	require.Nil(t, err)

	var resp data.Events
	err = json.Unmarshal(body, &resp)
	require.Nil(t, err)
	assert.Len(t, resp.Events, len(expectedEvents))
	for _, e := range expectedEvents {
		assert.Contains(t, resp.Events, e)
	}
}
