package registry_synchronizer_listeners

import (
	"sync"

	"github.com/ethereum/go-ethereum"
)

type AbstractUpdatableSubscriptionListener struct {
	subscription     ethereum.Subscription
	oldSubscriptions []ethereum.Subscription
	lastBlockSeen    uint64

	mu    sync.Mutex
	errCh chan error
	wg    sync.WaitGroup
}

func NewAbstractUpdatableSubscriptionListener() *AbstractUpdatableSubscriptionListener {
	return &AbstractUpdatableSubscriptionListener{
		oldSubscriptions: []ethereum.Subscription{},
		errCh:            make(chan error, 8),
	}
}

func (l *AbstractUpdatableSubscriptionListener) Listen() error {
	err := <-l.errCh

	l.mu.Lock()
	if l.subscription != nil {
		l.subscription.Unsubscribe()
	}
	for _, sub := range l.oldSubscriptions {
		sub.Unsubscribe()
	}
	l.mu.Unlock()

	l.wg.Wait()

	return err
}

func (l *AbstractUpdatableSubscriptionListener) HandleBlockProcessed(blockNumber uint64) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	prevLast := l.lastBlockSeen
	if blockNumber > prevLast {
		l.lastBlockSeen = blockNumber
	}

	if len(l.oldSubscriptions) > 0 && blockNumber >= prevLast {
		for _, sub := range l.oldSubscriptions {
			sub.Unsubscribe()
		}
		l.oldSubscriptions = nil
	}

	return nil
}
