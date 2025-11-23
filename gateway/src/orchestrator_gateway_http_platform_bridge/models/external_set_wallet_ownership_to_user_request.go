package models

type ExternalSetWalletOwnershipToUserRequest struct {
	Challenge     string `json:"challenge"`
	WalletAddress string `json:"walletAddress"`
}
