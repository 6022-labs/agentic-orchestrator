package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync_event_listeners"
	_ "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync_host/docs"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync_host/settings"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync_mvc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func ConfigureServer(container *dig.Container) {
	container.Provide(newHttpServer)

	err := container.Invoke(useRestApi)
	if err != nil {
		panic(err)
	}

	err = container.Invoke(useListeners)
	if err != nil {
		panic(err)
	}
}

func newHttpServer(hostSettings *settings.HostSettings) *fiber.App {
	if !hostSettings.UseApi {
		return nil
	}

	return fiber.New(fiber.Config{
		AppName:   hostSettings.AppName,
		BodyLimit: hostSettings.BodyLimitBytes,
	})
}

type useListenersParams struct {
	dig.In

	HostSettings              *settings.HostSettings
	AgentCollectionRepository repositories.AgentCollectionRepositoryInterface
	Listeners                 []orchestrator_sc_sync_event_listeners.ListenerInterface `group:"listeners"`
}

// RegisterListeners
func useListeners(p useListenersParams) {
	if !p.HostSettings.UseListeners {
		return
	}

	agentCollections, err := p.AgentCollectionRepository.GetAll()
	if err != nil {
		panic(err)
	}

	for _, listener := range p.Listeners {
		// If the listener is a CollectionListenerInterface, subscribe to the collections
		if collectionListener, ok := listener.(orchestrator_sc_sync_event_listeners.CollectionListenerInterface); ok {
			for _, collection := range agentCollections {
				err := collectionListener.Subscribe(common.HexToAddress(collection.Address))
				if err != nil {
					panic(err)
				}
			}
		}

		go func() {
			err := listener.Listen()
			if err != nil {
				panic(err)
			}
		}()
	}
}

type useRoutesParams struct {
	dig.In

	HostSettings *settings.HostSettings
	App          *fiber.App
	Logger       *zap.Logger
	Controllers  []orchestrator_sc_sync_mvc.ControllerInterface `group:"controllers"`
}

// RegisterRoutes hooks up the routes and uses Fx to create new controller instances per request
func useRestApi(p useRoutesParams) {
	if !p.HostSettings.UseApi {
		return
	}

	p.App.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	p.App.Use(fiberzap.New(fiberzap.Config{
		Logger: p.Logger,
	}))

	p.App.Get("/swagger/*", swagger.New())

	for _, controller := range p.Controllers {
		controller.InitRoutes(p.App)
	}
}
