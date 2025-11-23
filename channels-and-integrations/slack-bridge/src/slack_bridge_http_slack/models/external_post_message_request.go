package models

type ExternalPostMessageRequest struct {
	Channel  string  `json:"channel"`
	Text     string  `json:"text"`
	ThreadTs *string `json:"thread_ts,omitempty"`
	Username *string `json:"username"`
	IconUrl  *string `json:"icon_url"`
	Metadata *string `json:"metadata"`
}
