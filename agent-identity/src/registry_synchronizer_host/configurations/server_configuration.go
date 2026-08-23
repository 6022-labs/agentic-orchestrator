package configurations

import (
	"go.uber.org/dig"
)

const AppName = "6022-AgentIdentity-RegistrySynchronizer"

func ConfigureServer(container *dig.Container) {
	container.Provide(NewSynchronizer)
}
