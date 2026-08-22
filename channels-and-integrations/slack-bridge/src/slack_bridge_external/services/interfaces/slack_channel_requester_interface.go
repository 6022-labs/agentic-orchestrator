package interfaces

import "context"

type ChannelInformation struct {
	Id          string
	Name        string
	UserIds     []string
	Description string
	IsMember    bool
}

type SlackChannelRequesterInterface interface {
	JoinChannel(ctx context.Context, channelId string) error
	GetChannelInformation(ctx context.Context, channelId string) (*ChannelInformation, error)
}
