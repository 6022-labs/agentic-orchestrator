package models

type ExternalPostMessageResponse struct {
	Ok      bool   `json:"ok"`
	Channel string `json:"channel"`
	Ts      string `json:"ts"`
	Message struct {
		Text        string `json:"text"`
		Username    string `json:"username"`
		BotID       string `json:"bot_id"`
		Attachments []struct {
			Text     string `json:"text"`
			ID       int    `json:"id"`
			Fallback string `json:"fallback"`
		} `json:"attachments"`
		Type    string `json:"type"`
		Subtype string `json:"subtype"`
		Ts      string `json:"ts"`
	} `json:"message"`
}
