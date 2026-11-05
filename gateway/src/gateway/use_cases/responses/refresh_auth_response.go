package responses

import "time"

type RefreshAuthResponse struct {
	Token            string    `json:"token"`
	ExpiresAt        time.Time `json:"expiresAt"`
	RefreshToken     string    `json:"refreshToken"`
	RefreshExpiresAt time.Time `json:"refreshExpiresAt"`
}
