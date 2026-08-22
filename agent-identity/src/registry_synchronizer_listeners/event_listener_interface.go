package registry_synchronizer_listeners

import "context"

type EventListenerInterface interface {
	SubscribeAll(ctx context.Context) error
	Listen(ctx context.Context) error
}
