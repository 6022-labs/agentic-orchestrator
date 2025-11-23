package requests

type MessageDeletedRequest EventRequest[MessageDeletedEvent]

type MessageDeletedEvent struct {
	Channel         string                             `json:"channel"`
	ChannelType     string                             `json:"channel_type"`
	DeletedTS       string                             `json:"deleted_ts"`
	EventTS         string                             `json:"event_ts"`
	Hidden          bool                               `json:"hidden"`
	PreviousMessage MessageDeletedEventPreviousMessage `json:"previous_message"`
	Subtype         string                             `json:"subtype"`
	TS              string                             `json:"ts"`
	Type            string                             `json:"type"`
}

type MessageDeletedEventPreviousMessage struct {
	Team     string `json:"team"`
	Text     string `json:"text"`
	ThreadTS string `json:"thread_ts"`
	TS       string `json:"ts"`
	Type     string `json:"type"`
}
