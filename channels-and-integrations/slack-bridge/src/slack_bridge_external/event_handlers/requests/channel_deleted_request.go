package requests

type ChannelDeletedRequest EventRequest[ChannelDeletedEvent]

type ChannelDeletedEvent struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
}
