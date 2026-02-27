package services

import (
	"fmt"
	"time"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/google/uuid"
)

const (
	WALLET_OWNERSHIP_CHALLENGE_EXPIRATION    = 5 * 60 // 5 minutes in seconds
	WALLET_OWNERSHIP_CHALLENGE_PATH_TEMPLATE = "wallet_ownership_challenge:%s"
)

type WalletOwnershipChallenge struct {
	cacheStorage interfaces.CacheStorageInterface
}

func NewWalletOwnershipChallenge(
	cacheStorage interfaces.CacheStorageInterface,
) *WalletOwnershipChallenge {
	return &WalletOwnershipChallenge{
		cacheStorage: cacheStorage,
	}
}

func (w *WalletOwnershipChallenge) IssueChallenge(userId string) (*uuid.UUID, error) {
	nonce := uuid.New()

	duration := time.Duration(WALLET_OWNERSHIP_CHALLENGE_EXPIRATION) * time.Second

	err := w.cacheStorage.Store(
		fmt.Sprintf(WALLET_OWNERSHIP_CHALLENGE_PATH_TEMPLATE, nonce.String()),
		userId,
		&duration,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to store challenge for nonce %s: %w", nonce.String(), err)
	}

	return &nonce, nil
}

func (w *WalletOwnershipChallenge) ConsumeChallenge(nonce *uuid.UUID) (*string, error) {
	if nonce == nil {
		return nil, fmt.Errorf("nonce cannot be nil")
	}

	userId, err := w.cacheStorage.ReadString(fmt.Sprintf(WALLET_OWNERSHIP_CHALLENGE_PATH_TEMPLATE, nonce.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to read challenge for nonce %s: %w", nonce.String(), err)
	}
	if userId == nil {
		return nil, fmt.Errorf("challenge for nonce %s not found or expired", nonce.String())
	}

	return userId, nil
}
