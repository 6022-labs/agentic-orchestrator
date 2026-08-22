package registry_synchronizer_listeners

import (
	"context"

	"sync"

	"github.com/ethereum/go-ethereum"
)

type AbstractUpdatableSubscriptionListener struct {
	subscription     ethereum.Subscription
	oldSubscriptions []ethereum.Subscription
	lastBlockSeen    uint64
	needsRebuild     bool

	rebuildSubscription func(ctx context.Context) error

	mu    sync.Mutex
	errCh chan error
	wg    sync.WaitGroup
}

func NewAbstractUpdatableSubscriptionListener(
	rebuildSubscription func(ctx context.Context) error,
) *AbstractUpdatableSubscriptionListener {
	return &AbstractUpdatableSubscriptionListener{
		oldSubscriptions: []ethereum.Subscription{},
		errCh:            make(chan error, 8),

		rebuildSubscription: rebuildSubscription,
	}
}

func (l *AbstractUpdatableSubscriptionListener) Listen(ctx context.Context) error {
	var err error
	select {
	case err = <-l.errCh:
	case <-ctx.Done():
		// Graceful stop, but a subscription error already queued still wins.
		select {
		case err = <-l.errCh:
		default:
		}
	}

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

func (l *AbstractUpdatableSubscriptionListener) OnBlockProcessed(ctx context.Context, blockNumber uint64) {
	l.onBlockSeen(blockNumber)

	l.mu.Lock()
	if !l.needsRebuild {
		l.mu.Unlock()
		return
	}
	l.needsRebuild = false
	l.mu.Unlock()

	if err := l.rebuildSubscription(ctx); err != nil {
		select {
		case l.errCh <- err:
		default:
		}
	}
}

func (l *AbstractUpdatableSubscriptionListener) markNeedsRebuild() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.needsRebuild = true
}

func (l *AbstractUpdatableSubscriptionListener) onBlockSeen(blockNumber uint64) {
	l.mu.Lock()
	defer l.mu.Unlock()

	last := l.lastBlockSeen
	if blockNumber > last {
		l.lastBlockSeen = blockNumber
	}

	if len(l.oldSubscriptions) > 0 && blockNumber >= last {
		for _, sub := range l.oldSubscriptions {
			sub.Unsubscribe()
		}
		l.oldSubscriptions = nil
	}
}
