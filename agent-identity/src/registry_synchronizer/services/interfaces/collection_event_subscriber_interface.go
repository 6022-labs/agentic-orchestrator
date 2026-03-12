package interfaces

import "github.com/ethereum/go-ethereum/common"

type CollectionEventSubscriberInterface interface {
	Subscribe(collectionAddress common.Address) error
}
