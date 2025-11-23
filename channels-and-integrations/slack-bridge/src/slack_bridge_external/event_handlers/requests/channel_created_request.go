package requests

type ChannelCreatedRequest EventRequest[ChannelCreatedEvent]

type ChannelCreatedEvent struct {
	EventTs string                     `json:"event_ts"`
	Type    string                     `json:"type"`
	Channel ChannelCreatedEventChannel `json:"channel"`
}

type ChannelCreatedEventChannel struct {
	Id                      string      `json:"id"`
	IsChannel               bool        `json:"is_channel"`
	IsGroup                 bool        `json:"is_group"`
	IsIm                    bool        `json:"is_im"`
	IsMpim                  bool        `json:"is_mpim"`
	IsExtShared             bool        `json:"is_ext_shared"`
	IsOrgShared             bool        `json:"is_org_shared"`
	IsShared                bool        `json:"is_shared"`
	IsPrivate               bool        `json:"is_private"`
	IsArchived              bool        `json:"is_archived"`
	IsFrozen                bool        `json:"is_frozen"`
	IsGeneral               bool        `json:"is_general"`
	IsPendingExtShared      bool        `json:"is_pending_ext_shared"`
	SharedTeamIds           []string    `json:"shared_team_ids"`
	PendingShared           []string    `json:"pending_shared"`
	PendingConnectedTeamIds []string    `json:"pending_connected_team_ids"`
	ParentConversation      interface{} `json:"parent_conversation"`
	PreviousNames           []string    `json:"previous_names"`
	Name                    string      `json:"name"`
	NameNormalized          string      `json:"name_normalized"`
	Purpose                 struct {
		Created int64  `json:"created"`
		Creator string `json:"creator"`
		LastSet int64  `json:"last_set"`
		Value   string `json:"value"`
	} `json:"purpose"`
	Topic struct {
		Created int64  `json:"created"`
		Creator string `json:"creator"`
		LastSet int64  `json:"last_set"`
		Value   string `json:"value"`
	} `json:"topic"`
	Unlinked int64 `json:"unlinked"`
	Updated  int64 `json:"updated"`
}
