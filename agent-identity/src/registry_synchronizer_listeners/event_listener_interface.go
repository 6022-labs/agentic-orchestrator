package registry_synchronizer_listeners

type EventListenerInterface interface {
	SubscribeAll() error
	Listen() error
}
