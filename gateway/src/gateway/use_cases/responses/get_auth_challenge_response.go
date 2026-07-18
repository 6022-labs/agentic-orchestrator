package responses

import "time"

type GetAuthChallengeResponse struct {
	Message   string    `json:"message"`
	ExpiresAt time.Time `json:"expiresAt"`
}
