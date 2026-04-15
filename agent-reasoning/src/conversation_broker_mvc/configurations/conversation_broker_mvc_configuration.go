package configurations

import (
	"go.uber.org/dig"
)

// AddConversationBrokerMvcConfiguration registers HTTP controllers as
// dig.Group("controllers"). Empty in C1 — controllers land in C4.
func AddConversationBrokerMvcConfiguration(container *dig.Container) {
}
