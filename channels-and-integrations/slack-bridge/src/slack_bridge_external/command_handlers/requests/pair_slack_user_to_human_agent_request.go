package requests

type PairSlackUserToHumanAgentRequest struct {
	SlackUserId                 string `json:"slack_user_id"`
	RequesterSlackUserId        string `json:"requester_slack_user_id"`
	HumanAgentCollectionAddress string `form:"human_agent_collection_address"`
	HumanAgentCollectionTokenId uint64 `form:"human_agent_collection_token_id"`
}
