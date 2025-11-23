package requests

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type MessageCreatedRequest EventRequest[MessageCreatedEvent]

type MessageCreatedEvent struct {
	Ts       string  `json:"ts"`
	EventTs  string  `json:"event_ts"`
	ThreadTs *string `json:"thread_ts,omitempty"` // Filled if the message is part of a thread

	Text string `json:"text"`

	Type        string `json:"type"`
	Subtype     string `json:"subtype,omitempty"` // Will be empty for normal messages, "bot_message" for bot messages, "message_changed" for edited messages, "message_deleted" for deleted messages
	Channel     string `json:"channel"`
	ChannelType string `json:"channel_type"`

	User *string `json:"user,omitempty"` // Filled if the message is from a user

	Metadata *MessageCreatedMetadata `json:"metadata,omitempty"`

	BotId    *string `json:"bot_id,omitempty"` // Filled if the message is from a bot
	AppId    *string `json:"app_id,omitempty"` // Filled if the message is from an bot
	Username *string `json:"username"`         // Filled if the message is from a bot
}

type MessageCreatedMetadata struct {
	EventType    string         `json:"event_type"`
	EventPayload map[string]any `json:"event_payload"`
}

func (m *MessageCreatedMetadata) ToAgentIdentifiers() (*string, *uint64, error) {
	collectionTokenAddressRaw, ok := m.EventPayload["collection_token_address"]
	if !ok {
		return nil, nil, fmt.Errorf("metadata missing collection_token_address")
	}
	collectionTokenAddress, ok := collectionTokenAddressRaw.(string)
	if !ok {
		return nil, nil, fmt.Errorf("invalid collection_token_address in metadata, got type %T with value %v", collectionTokenAddressRaw, collectionTokenAddressRaw)
	}
	if !common.IsHexAddress(collectionTokenAddress) {
		return nil, nil, fmt.Errorf("invalid collection_token_address in metadata, got non-hex value %v", collectionTokenAddress)
	}

	collectionTokenIdRaw, ok := m.EventPayload["collection_token_id"]
	if !ok {
		return nil, nil, fmt.Errorf("metadata missing collection_token_id")
	}
	collectionTokenIdFloat, ok := collectionTokenIdRaw.(float64)
	if !ok {
		return nil, nil, fmt.Errorf("invalid collection_token_id in metadata, got type %T with value %v", collectionTokenIdRaw, collectionTokenIdRaw)
	}
	if collectionTokenIdFloat < 0 {
		return nil, nil, fmt.Errorf("invalid collection_token_id in metadata, got negative value %v", collectionTokenIdFloat)
	}
	if collectionTokenIdFloat != float64(uint64(collectionTokenIdFloat)) {
		return nil, nil, fmt.Errorf("invalid collection_token_id in metadata, got non-integer value %v", collectionTokenIdFloat)
	}

	collectionTokenIdUint64 := uint64(collectionTokenIdFloat)

	return &collectionTokenAddress, &collectionTokenIdUint64, nil
}
