package requests

type MessageChangedRequest EventRequest[MessageChangedEvent]

type MessageChangedEvent struct {
	Channel         string                             `json:"channel"`
	ChannelType     string                             `json:"channel_type"`
	EventTS         string                             `json:"event_ts"`
	Hidden          bool                               `json:"hidden"`
	Message         MessageChangedEventMessage         `json:"message"`
	PreviousMessage MessageChangedEventPreviousMessage `json:"previous_message"`
	Subtype         string                             `json:"subtype"`
	TS              string                             `json:"ts"`
	Type            string                             `json:"type"`
}

type MessageChangedEventMessage struct {
	Hidden   bool   `json:"hidden"`
	IsLocked bool   `json:"is_locked"`
	Team     string `json:"team"`
	Text     string `json:"text"`
	ThreadTS string `json:"thread_ts"`
	TS       string `json:"ts"`
	Type     string `json:"type"`
}

type MessageChangedEventPreviousMessage struct {
	Hidden   bool   `json:"hidden"`
	Team     string `json:"team"`
	Text     string `json:"text"`
	ThreadTS string `json:"thread_ts"`
	TS       string `json:"ts"`
	Type     string `json:"type"`
}
