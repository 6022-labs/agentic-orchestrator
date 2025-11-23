package interfaces

type SlackThreadConcludeInterface interface {
	Conclude(threadId, channelId string) error
	ConcludeIfExists(threadId, channelId string) error
}
