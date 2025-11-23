package models

type ExternalGetConversationMembersResponse struct {
	Ok               bool     `json:"ok"`
	Members          []string `json:"members"`
	ResponseMetadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
}
