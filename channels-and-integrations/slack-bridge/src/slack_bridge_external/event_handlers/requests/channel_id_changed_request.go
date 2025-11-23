package requests

type ChannelIdChangedRequest EventRequest[ChannelIdChangedEvent]

type ChannelIdChangedEvent struct {
	Type         string `json:"type"`
	OldChannelId string `json:"old_channel_id"`
	NewChannelId string `json:"new_channel_id"`
}
