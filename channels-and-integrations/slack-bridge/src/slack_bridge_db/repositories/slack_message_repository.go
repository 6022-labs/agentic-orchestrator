package repositories

import (
	"errors"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"gorm.io/gorm"
)

type SlackMessageRepository struct {
	AbstractRepository[models.SlackMessage]
}

func NewSlackMessageRepository(
	client *slack_bridge_db.DatabaseClient,
) *SlackMessageRepository {
	return &SlackMessageRepository{
		AbstractRepository: *NewAbstractRepository[models.SlackMessage](client),
	}
}

func (r *SlackMessageRepository) GetConversationByMessageId(id string) ([]*models.SlackMessage, error) {
	var msg models.SlackMessage
	err := r.client.GetDB().First(&msg, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*models.SlackMessage{}, nil
		}
		return nil, err
	}

	var conversation []*models.SlackMessage
	err = r.client.GetDB().
		Where("thread_id = ? AND channel_id = ?", msg.ThreadId, msg.ChannelId).
		Order("id ASC").
		Find(&conversation).Error
	if err != nil {
		return nil, err
	}

	return conversation, nil
}

func (r *SlackMessageRepository) GetByThreadIdAndChannelId(threadId, channelId string) ([]*models.SlackMessage, error) {
	var messages []*models.SlackMessage
	err := r.client.GetDB().
		Where("thread_id = ? AND channel_id = ?", threadId, channelId).
		Order("id ASC").
		Find(&messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *SlackMessageRepository) SaveMany(entities []*models.SlackMessage) ([]*models.SlackMessage, error) {
	if len(entities) == 0 {
		return nil, nil
	}

	err := r.client.GetDB().Create(&entities).Error
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *SlackMessageRepository) GetProcessableMessagesForAgentAndThread(agent *models.SlackAgent, threadId, channelId string) ([]*models.SlackMessage, error) {
	query := r.client.GetDB().Model(&models.SlackMessage{}).
		Where("thread_id = ? AND channel_id = ?", threadId, channelId)

	if agent.Role == enums.AgentRoleExpert || agent.Role == enums.AgentRoleClone {
		// Get messages that only mention the agent
		query = query.Where("content ILIKE ?", "%"+agent.EnsDomain+"%")
	}

	if agent.Role == enums.AgentRoleHuman {
		// Get messages that are only from this human agent (for storing data from them)
		query = query.
			Where("(author_collection_address = ? AND author_collection_token_id = ?) OR content ILIKE ?",
				agent.CollectionAddress,
				agent.CollectionTokenId,
				"%"+agent.EnsDomain+"%",
			)
	}

	var messages []*models.SlackMessage
	err := query.Order("id ASC").Find(&messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}
