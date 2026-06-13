package clients_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_http_ipfs/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_http_ipfs/models"
	ipfs_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_http_ipfs/settings"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenAddingFileTestingSuite struct {
	sut    *clients.IPFSClient
	server *httptest.Server
}

func WhenAddingFileBeforeEach(t *testing.T, handler http.HandlerFunc, bearerToken string) *WhenAddingFileTestingSuite {
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

	return &WhenAddingFileTestingSuite{sut: sut, server: server}
}

func WhenAddingFileAfterEach(suite *WhenAddingFileTestingSuite) {
	suite.server.Close()
}

func TestWhenAddingFile(t *testing.T) {
	t.Parallel()

	pin := true

	t.Run("Given the IPFS node accepts the file and a bearer token is configured", func(t *testing.T) {
		t.Parallel()

		var seenMethod, seenPath, seenQuery, seenContentType, seenAuth string
		var seenBody []byte
		handler := func(w http.ResponseWriter, r *http.Request) {
			seenMethod = r.Method
			seenPath = r.URL.Path
			seenQuery = r.URL.RawQuery
			seenContentType = r.Header.Get("Content-Type")
			seenAuth = r.Header.Get("Authorization")
			seenBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"Name":"agent.png","Hash":"QmHashValue","Size":"123"}`))
		}

		t.Run("Should POST multipart to /api/v0/add with pin, the bearer header, and decode the hash", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingFileBeforeEach(t, handler, "secret-token123")
			defer WhenAddingFileAfterEach(suite)

			resp, err := suite.sut.Add([]byte("image-bytes"), "agent.png", &models.ExternalAddQueryRequest{Pin: &pin})

			assert.NoError(t, err)
			if assert.NotNil(t, resp) {
				assert.Equal(t, "QmHashValue", resp.Hash)
			}
			assert.Equal(t, http.MethodPost, seenMethod)
			assert.Equal(t, "/api/v0/add", seenPath)
			assert.Contains(t, seenQuery, "pin=true")
			assert.Contains(t, seenContentType, "multipart/form-data")
			assert.Equal(t, "Bearer secret-token123", seenAuth)
			assert.Contains(t, string(seenBody), "image-bytes")
			assert.Contains(t, string(seenBody), "agent.png")
		})
	})

	t.Run("Given no bearer token is configured", func(t *testing.T) {
		t.Parallel()

		var seenAuth string
		handler := func(w http.ResponseWriter, r *http.Request) {
			seenAuth = r.Header.Get("Authorization")
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"Hash":"QmHashValue"}`))
		}

		t.Run("Should not send an Authorization header", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingFileBeforeEach(t, handler, "")
			defer WhenAddingFileAfterEach(suite)

			_, err := suite.sut.Add([]byte("x"), "f.png", &models.ExternalAddQueryRequest{Pin: &pin})

			assert.NoError(t, err)
			assert.Empty(t, seenAuth)
		})
	})

	t.Run("Given the IPFS node returns a non-2xx", func(t *testing.T) {
		t.Parallel()

		handler := func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("boom"))
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingFileBeforeEach(t, handler, "")
			defer WhenAddingFileAfterEach(suite)

			_, err := suite.sut.Add([]byte("x"), "f.png", &models.ExternalAddQueryRequest{Pin: &pin})

			assert.Error(t, err)
		})
	})

	t.Run("Given the IPFS node returns a malformed body", func(t *testing.T) {
		t.Parallel()

		handler := func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("not-json"))
		}

		t.Run("Should return a decode error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingFileBeforeEach(t, handler, "")
			defer WhenAddingFileAfterEach(suite)

			_, err := suite.sut.Add([]byte("x"), "f.png", &models.ExternalAddQueryRequest{Pin: &pin})

			assert.Error(t, err)
		})
	})
}
