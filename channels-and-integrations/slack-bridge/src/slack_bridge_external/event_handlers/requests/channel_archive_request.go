package requests

type ChannelArchiveRequest EventRequest[ChannelArchiveEvent]

type ChannelArchiveEvent struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
}
