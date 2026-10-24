package mvc_test

import (
	"context"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	configuration "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/configurations"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	service_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	db_configuration "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/configurations"
	registry_clients "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_agent_identity_registry/clients"
	registry_configuration "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_agent_identity_registry/configurations"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_mvc"
	mvc_configuration "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_mvc/configurations"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/container_builders"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type httpSuite struct {
	dbName        string
	app           *fiber.App
	registryStub  *httptest.Server
	registryAgent *registry_clients.ResponsesAgentResponse
}

// Stubs satisfy the broker's node-discovery and cache ports; the http flow
// under test never reaches an agent node.
type stubEnsUrlResolver struct{}

func (s *stubEnsUrlResolver) ResolveUrl(_ context.Context, ensDomain string) (string, error) {
	return "", fmt.Errorf("no agent node in integration tests")
}

type stubRemoteAgentCardRequester struct{}

func (s *stubRemoteAgentCardRequester) GetCard(
	_ context.Context,
	agentUrl string,
) (*models.RemoteAgentCard, []byte, error) {
	return nil, nil, fmt.Errorf("no agent node in integration tests")
}

type stubCacheStorage struct{}

func (s *stubCacheStorage) Store(_ context.Context, key, value string, expiration *time.Duration) error {
	return nil
}
func (s *stubCacheStorage) Read(_ context.Context, key string) (string, error) {
	return "", fmt.Errorf("cache miss")
}
func (s *stubCacheStorage) Delete(_ context.Context, key string) error { return nil }

func newHttpSuite() *httpSuite {
	dbName, connStr := container_builders.CreatePostgreSQLDatabase()
	migrationFolder, err := common_tests.FindConversationBrokerMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	suite := &httpSuite{dbName: dbName}

	suite.registryStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if suite.registryAgent == nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(suite.registryAgent)
	}))

	k := koanf.New(".")
	if err := k.Load(confmap.Provider(map[string]interface{}{
		"database.url":                     connStr,
		"migrations.folder_path":           *migrationFolder,
		"agent_identity_registry.base_url": suite.registryStub.URL,
	}, "."), nil); err != nil {
		panic(err)
	}

	container := dig.New()
	if err := container.Provide(func() *zap.Logger { return zap.NewNop() }); err != nil {
		panic(err)
	}
	if err := container.Provide(func() *koanf.Koanf { return k }); err != nil {
		panic(err)
	}

	if err := container.Provide(
		func() service_interfaces.EnsUrlResolverInterface { return &stubEnsUrlResolver{} },
	); err != nil {
		panic(err)
	}
	if err := container.Provide(
		func() service_interfaces.RemoteAgentCardRequesterInterface { return &stubRemoteAgentCardRequester{} },
	); err != nil {
		panic(err)
	}
	if err := container.Provide(
		func() service_interfaces.CacheStorageInterface { return &stubCacheStorage{} },
	); err != nil {
		panic(err)
	}

	db_configuration.AddConversationBrokerDbConfiguration(container)
	configuration.AddConversationBrokerConfiguration(container)
	mvc_configuration.AddConversationBrokerMvcConfiguration(container)
	registry_configuration.AddConversationBrokerHttpAgentIdentityRegistryConfiguration(container)

	if err := container.Invoke(func(client db_clients.DatabaseClientInterface) {
		if err := client.Migrate(); err != nil {
			panic(err)
		}
	}); err != nil {
		panic(err)
	}

	app := fiber.New()
	type ctrlsParams struct {
		dig.In
		Controllers []conversation_broker_mvc.ControllerInterface `group:"controllers"`
	}
	if err := container.Invoke(func(p ctrlsParams) {
		for _, c := range p.Controllers {
			c.InitRoutes(app)
		}
	}); err != nil {
		panic(err)
	}

	suite.app = app
	return suite
}

func (s *httpSuite) teardown() {
	if s.registryStub != nil {
		s.registryStub.Close()
	}
	container_builders.DestroyPostgreSQLDatabase(s.dbName)
}

func (s *httpSuite) do(method, path string, body any) (*http.Response, []byte) {
	var reader io.Reader
	if body != nil {
		buf, _ := json.Marshal(body)
		reader = bytes.NewReader(buf)
	}
	req, _ := http.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := s.app.Test(req, fiber.TestConfig{Timeout: 0, FailOnTimeout: false})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	out, _ := io.ReadAll(resp.Body)
	return resp, out
}

func TestWhenCallingEndpoints(t *testing.T) {
	t.Parallel()

	t.Run("Given a clean broker", func(t *testing.T) {
		t.Run("Should support full register/attach/ingest/get/conclude flow", func(t *testing.T) {
			suite := newHttpSuite()
			defer suite.teardown()

			// 1. Register bridge
			resp, body := suite.do(http.MethodPost, "/bridges", &requests.RegisterBridgeRequest{
				Kind:    enums.BridgeKindSlack,
				BaseUrl: "http://slack-bridge.local",
			})
			assert.Equal(t, http.StatusCreated, resp.StatusCode)
			var bridge models.Bridge
			assert.NoError(t, json.Unmarshal(body, &bridge))

			// 2. Register an agent participant (broker fetches data from registry stub)
			//    + a user participant
			agentId := uuid.New()
			tokenId := uint64(1)
			suite.registryAgent = &registry_clients.ResponsesAgentResponse{
				Id:                lo.ToPtr(agentId.String()),
				Attributes:        lo.ToPtr(map[string]string{"role": "clone"}),
				EnsDomain:         lo.ToPtr("demo.eth"),
				CollectionAddress: lo.ToPtr("0xabc00000000000000000000000000000000000de"),
				CollectionTokenId: lo.ToPtr(int(tokenId)),
			}
			resp, _ = suite.do(http.MethodPost, "/participants/agents", &requests.RegisterAgentRequest{
				CollectionAddress: "0xAbC00000000000000000000000000000000000dE",
				CollectionTokenId: tokenId,
			})
			assert.Equal(t, http.StatusCreated, resp.StatusCode)

			resp, body = suite.do(http.MethodPost, "/participants/users", &requests.RegisterUserRequest{
				DisplayName: "Pierre",
			})
			assert.Equal(t, http.StatusCreated, resp.StatusCode)
			var human models.Participant
			assert.NoError(t, json.Unmarshal(body, &human))

			// 3. Register conversation, attaching both participants in one call
			desc := "Test channel"
			resp, body = suite.do(http.MethodPost, "/conversations", &requests.RegisterConversationRequest{
				BridgeID:       bridge.ID,
				Description:    &desc,
				ParticipantIDs: []uuid.UUID{agentId, human.ID},
			})
			assert.Equal(t, http.StatusCreated, resp.StatusCode)
			var conv models.Conversation
			assert.NoError(t, json.Unmarshal(body, &conv))

			// 4. Detach + re-attach the human
			resp, _ = suite.do(
				http.MethodDelete,
				fmt.Sprintf("/conversations/%s/participants/%s", conv.ID, human.ID),
				nil,
			)
			assert.Equal(t, http.StatusNoContent, resp.StatusCode)

			resp, _ = suite.do(http.MethodPost, fmt.Sprintf("/conversations/%s/participants", conv.ID), map[string]any{
				"participantId": human.ID,
			})
			assert.Equal(t, http.StatusCreated, resp.StatusCode)

			// 5. Ingest a message from the human
			resp, _ = suite.do(
				http.MethodPost,
				fmt.Sprintf("/conversations/%s/messages", conv.ID),
				&requests.IngestMessageRequest{
					ParticipantID: human.ID,
					Content:       "hello",
				},
			)
			assert.Equal(t, http.StatusCreated, resp.StatusCode)

			// 6. Get conversation
			resp, body = suite.do(http.MethodGet, fmt.Sprintf("/conversations/%s", conv.ID), nil)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			var got models.Conversation
			assert.NoError(t, json.Unmarshal(body, &got))
			assert.Equal(t, conv.ID, got.ID)
			assert.Equal(t, enums.ConversationStatusActive, got.Status)

			// 7. Force conclude
			resp, _ = suite.do(http.MethodPost, fmt.Sprintf("/conversations/%s/force-conclude", conv.ID), nil)
			assert.Equal(t, http.StatusNoContent, resp.StatusCode)

			resp, body = suite.do(http.MethodGet, fmt.Sprintf("/conversations/%s", conv.ID), nil)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.NoError(t, json.Unmarshal(body, &got))
			assert.Equal(t, enums.ConversationStatusConcluded, got.Status)
			assert.NotNil(t, got.ConcludedAt)
		})

		t.Run("Should reject invalid registerBridge with 400", func(t *testing.T) {
			suite := newHttpSuite()
			defer suite.teardown()

			resp, _ := suite.do(http.MethodPost, "/bridges", &requests.RegisterBridgeRequest{
				Kind:    enums.BridgeKindSlack,
				BaseUrl: "",
			})
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("Should return 404 when getting unknown conversation", func(t *testing.T) {
			suite := newHttpSuite()
			defer suite.teardown()

			resp, _ := suite.do(http.MethodGet, "/conversations/"+uuid.New().String(), nil)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	})
}
