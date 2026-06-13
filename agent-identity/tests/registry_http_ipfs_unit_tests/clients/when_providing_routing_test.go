package clients_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_http_ipfs/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_http_ipfs/models"
	ipfs_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_http_ipfs/settings"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenProvidingRoutingTestingSuite struct {
	sut    *clients.IPFSClient
	server *httptest.Server
}

func WhenProvidingRoutingBeforeEach(t *testing.T, handler http.HandlerFunc, bearerToken string) *WhenProvidingRoutingTestingSuite {
	t.Helper()

	server := httptest.NewServer(handler)
	sut := clients.NewIPFSClient(clients.NewIPFSClientParams{
		Logger:     zap.NewNop(),
		HTTPClient: http.DefaultClient,
		IPFSSettings: &ipfs_settings.IPFSSettings{
			NodeHttpUrl: server.URL,
			BearerToken: bearerToken,
		},
	})

	return &WhenProvidingRoutingTestingSuite{sut: sut, server: server}
}

func WhenProvidingRoutingAfterEach(suite *WhenProvidingRoutingTestingSuite) {
	suite.server.Close()
}

func TestWhenProvidingRouting(t *testing.T) {
	t.Parallel()

	t.Run("Given the IPFS node accepts the provide and a bearer token is configured", func(t *testing.T) {
		t.Parallel()

		var seenMethod, seenPath, seenQuery, seenAuth string
		handler := func(w http.ResponseWriter, r *http.Request) {
			seenMethod = r.Method
			seenPath = r.URL.Path
			seenQuery = r.URL.RawQuery
			seenAuth = r.Header.Get("Authorization")
			w.WriteHeader(http.StatusOK)
		}

		t.Run("Should POST to /api/v0/routing/provide with the cid arg and the bearer header", func(t *testing.T) {
			t.Parallel()

			suite := WhenProvidingRoutingBeforeEach(t, handler, "secret-token123")
			defer WhenProvidingRoutingAfterEach(suite)

			err := suite.sut.RoutingProvideWithoutResponse(&models.ExternalRoutingProvideRequest{Arg: "QmHashValue"})

			assert.NoError(t, err)
			assert.Equal(t, http.MethodPost, seenMethod)
			assert.Equal(t, "/api/v0/routing/provide", seenPath)
			assert.Contains(t, seenQuery, "arg=QmHashValue")
			assert.Equal(t, "Bearer secret-token123", seenAuth)
		})
	})

	t.Run("Given the IPFS node returns a non-2xx", func(t *testing.T) {
		t.Parallel()

		handler := func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenProvidingRoutingBeforeEach(t, handler, "")
			defer WhenProvidingRoutingAfterEach(suite)

			err := suite.sut.RoutingProvideWithoutResponse(&models.ExternalRoutingProvideRequest{Arg: "QmX"})

			assert.Error(t, err)
		})
	})
}
