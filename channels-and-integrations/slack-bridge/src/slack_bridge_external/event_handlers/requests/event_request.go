package requests

type EventRequest[T any] struct {
	ApiAppId       string `json:"api_app_id"`
	Authorizations []struct {
		EnterpriseId        interface{} `json:"enterprise_id"`
		IsBot               bool        `json:"is_bot"`
		IsEnterpriseInstall bool        `json:"is_enterprise_install"`
		TeamId              string      `json:"team_id"`
		UserId              string      `json:"user_id"`
	} `json:"authorizations"`
	Event              T      `json:"event"`
	EventId            string `json:"event_id"`
	EventTime          int64  `json:"event_time"`
	IsExtSharedChannel bool   `json:"is_ext_shared_channel"`
	TeamId             string `json:"team_id"`
	Token              string `json:"token"`
	Type               string `json:"type"`
}
