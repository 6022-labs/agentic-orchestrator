package responses

import "time"

type GetAuthChallengeResponse struct {
	Nonce     string    `json:"nonce"`
	ExpiresAt time.Time `json:"expiresAt"`
}
