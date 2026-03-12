package configurations

import (
	"go.uber.org/dig"
)

func ConfigureServer(container *dig.Container) {
	container.Provide(NewSynchronizer)
}
