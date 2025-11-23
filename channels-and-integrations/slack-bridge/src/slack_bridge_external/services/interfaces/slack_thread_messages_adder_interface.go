package interfaces

type SlackThreadMessagesAdderInterface interface {
	Add(
		channelId, parentMessageTs, messageId, text string,
		authorUserId, authorCollectionAddress *string, authorCollectionTokenId *uint64,
	) error
}
