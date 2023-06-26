package contract_test

import (
	"context"
	"os"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)


func TestProviderEvents(t *testing.T) {
	if os.Getenv("CONTRACT") == "" {
		t.Skip("Skipping TestProviderEvents in short mode.")
	}

	// Initialize
	pact := dsl.Pact{
		Provider: "ConfTalks-Server",
	}
	pact.Setup(true)

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

	// Verify
	if os.Getenv("REMOTE") != "" {
		_, err = pact.VerifyProvider(t, types.VerifyRequest{
			ProviderBaseURL: url,
			BrokerURL:       os.Getenv("PACT_BROKER_BASE_URL"),
			BrokerToken:     os.Getenv("PACT_BROKER_TOKEN"),
			ProviderVersion: os.Getenv("version"),
		})
		require.Nil(t, err)
		return
	}
	_, err = pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL: url,
		PactURLs:        []string{PACTS_PATH},
	})
	require.Nil(t, err)
}
