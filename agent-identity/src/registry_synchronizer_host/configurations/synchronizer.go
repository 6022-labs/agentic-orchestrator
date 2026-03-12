package configurations

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_listeners"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type Synchronizer struct {
	logger *zap.Logger

	eventListeners []registry_synchronizer_listeners.EventListenerInterface

	blockEventsGapEnqueuer   interfaces.BlockEventsGapEnqueuerInterface
	blockEventsDequeueDaemon interfaces.BlockEventsDequeueDaemonInterface
}

type newSynchronizerParams struct {
	dig.In

	Logger *zap.Logger

	EventListeners []registry_synchronizer_listeners.EventListenerInterface `group:"event_listeners"`

	BlockEventsGapEnqueuer   interfaces.BlockEventsGapEnqueuerInterface
	BlockEventsDequeueDaemon interfaces.BlockEventsDequeueDaemonInterface
}

func NewSynchronizer(
	params newSynchronizerParams,
) *Synchronizer {
	return &Synchronizer{
		logger:                   params.Logger,
		eventListeners:           params.EventListeners,
		blockEventsGapEnqueuer:   params.BlockEventsGapEnqueuer,
		blockEventsDequeueDaemon: params.BlockEventsDequeueDaemon,
	}
}

func (app *Synchronizer) Start() error {
	for _, listener := range app.eventListeners {
		if err := listener.SubscribeAll(); err != nil {
			return fmt.Errorf("failed to subscribe listener: %w", err)
		}

		go func(listener registry_synchronizer_listeners.EventListenerInterface) {
			if err := listener.Listen(); err != nil {
				app.logger.Error("Listener stopped", zap.Error(err))
				panic(err)
			}
		}(listener)
	}

	err := app.blockEventsGapEnqueuer.EnqueueGap()
	if err != nil {
		app.logger.Error("Error syncing", zap.Error(err))
		return err
	}

	return app.blockEventsDequeueDaemon.Start()
}
