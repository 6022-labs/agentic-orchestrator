package requests

type ChannelLeftRequest EventRequest[ChannelLeftEvent]

type ChannelLeftEvent struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
}
