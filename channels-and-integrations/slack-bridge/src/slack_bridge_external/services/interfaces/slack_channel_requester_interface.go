package interfaces

type ChannelInformation struct {
	Id          string
	Name        string
	Members     []string
	Description string
	IsMember    bool
}

type SlackChannelRequesterInterface interface {
	JoinChannel(channelId string) error
	GetChannelInformation(channelId string) (*ChannelInformation, error)
}
