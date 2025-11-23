package models

type ExternalGetConversationInfoResponse struct {
	Ok      bool `json:"ok"`
	Channel struct {
		ID                      string        `json:"id"`
		Name                    string        `json:"name"`
		IsChannel               bool          `json:"is_channel"`
		IsGroup                 bool          `json:"is_group"`
		IsIm                    bool          `json:"is_im"`
		IsMpim                  bool          `json:"is_mpim"`
		IsPrivate               bool          `json:"is_private"`
		Created                 int           `json:"created"`
		IsArchived              bool          `json:"is_archived"`
		IsGeneral               bool          `json:"is_general"`
		Unlinked                int           `json:"unlinked"`
		NameNormalized          string        `json:"name_normalized"`
		IsShared                bool          `json:"is_shared"`
		IsFrozen                bool          `json:"is_frozen"`
		IsOrgShared             bool          `json:"is_org_shared"`
		IsPendingExtShared      bool          `json:"is_pending_ext_shared"`
		PendingShared           []interface{} `json:"pending_shared"`
		ContextTeamID           string        `json:"context_team_id"`
		Updated                 int64         `json:"updated"`
		ParentConversation      interface{}   `json:"parent_conversation"`
		Creator                 string        `json:"creator"`
		IsMember                *bool         `json:"is_member,omitempty"`
		IsExtShared             bool          `json:"is_ext_shared"`
		SharedTeamIds           []string      `json:"shared_team_ids"`
		PendingConnectedTeamIds []interface{} `json:"pending_connected_team_ids"`
		Topic                   struct {
			Value   string `json:"value"`
			Creator string `json:"creator"`
			LastSet int    `json:"last_set"`
		} `json:"topic"`
		Purpose struct {
			Value   string `json:"value"`
			Creator string `json:"creator"`
			LastSet int    `json:"last_set"`
		} `json:"purpose"`
		Properties struct {
			Tabs []struct {
				ID    string `json:"id"`
				Label string `json:"label"`
				Type  string `json:"type"`
			} `json:"tabs"`
		} `json:"properties"`
		PreviousNames []interface{} `json:"previous_names"`
	} `json:"channel"`
}
