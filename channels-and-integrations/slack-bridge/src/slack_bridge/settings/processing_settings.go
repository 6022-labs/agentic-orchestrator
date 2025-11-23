package settings

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type ProcessingSettings struct {
	AgentProcessingInterval                     uint64
	FacilitatorTickInterval                     uint64
	ConversationForceConcludeByFacilitatorAfter uint64
	ConversationForceConcludeAfter              uint64
}

func NewProcessingSettings(
	logger *zap.Logger,
) *ProcessingSettings {
	facilitatorTickIntervalStr := os.Getenv("FACILITATOR_TICK_INTERVAL")
	if len(strings.TrimSpace(facilitatorTickIntervalStr)) == 0 {
		facilitatorTickIntervalStr = "30"

		logger.Info("FACILITATOR_TICK_INTERVAL is not set, using default value",
			zap.String("defaultFacilitatorTickInterval", facilitatorTickIntervalStr),
		)
	}

	facilitatorTickInterval, err := strconv.ParseUint(facilitatorTickIntervalStr, 10, 32)
	if err != nil {
		logger.Fatal("Failed to parse FACILITATOR_TICK_INTERVAL", zap.Error(err))
	}

	conversationForceConcludeByFacilitatorAfterStr := os.Getenv("CONVERSATION_FORCE_CONCLUDE_BY_FACILITATOR_AFTER")
	if len(strings.TrimSpace(conversationForceConcludeByFacilitatorAfterStr)) == 0 {
		conversationForceConcludeByFacilitatorAfterStr = "300"

		logger.Info("CONVERSATION_FORCE_CONCLUDE_BY_FACILITATOR_AFTER is not set, using default value",
			zap.String("defaultConversationForceConcludeByFacilitatorAfter", conversationForceConcludeByFacilitatorAfterStr),
		)
	}

	conversationForceConcludeByFacilitatorAfter, err := strconv.ParseUint(conversationForceConcludeByFacilitatorAfterStr, 10, 32)
	if err != nil {
		logger.Fatal("Failed to parse CONVERSATION_FORCE_CONCLUDE_BY_FACILITATOR_AFTER", zap.Error(err))
	}

	conversationForceConcludeAfterStr := os.Getenv("CONVERSATION_FORCE_CONCLUDE_AFTER")
	if len(strings.TrimSpace(conversationForceConcludeAfterStr)) == 0 {
		conversationForceConcludeAfterStr = fmt.Sprintf("%d", conversationForceConcludeByFacilitatorAfter*2)

		logger.Info("CONVERSATION_FORCE_CONCLUDE_AFTER is not set, using default value (x2 of CONVERSATION_FORCE_CONCLUDE_BY_FACILITATOR_AFTER)",
			zap.String("defaultConversationForceConcludeAfter", conversationForceConcludeAfterStr),
		)
	}

	conversationForceConcludeAfter, err := strconv.ParseUint(conversationForceConcludeAfterStr, 10, 32)
	if err != nil {
		logger.Fatal("Failed to parse CONVERSATION_FORCE_CONCLUDE_AFTER", zap.Error(err))
	}

	agentProcessingIntervalStr := os.Getenv("AGENT_PROCESSING_INTERVAL")
	if len(strings.TrimSpace(agentProcessingIntervalStr)) == 0 {
		agentProcessingIntervalStr = "20"

		logger.Info("AGENT_PROCESSING_INTERVAL is not set, using default value",
			zap.String("defaultAgentProcessingInterval", agentProcessingIntervalStr),
		)
	}

	agentProcessingInterval, err := strconv.ParseUint(agentProcessingIntervalStr, 10, 32)
	if err != nil {
		logger.Fatal("Failed to parse AGENT_PROCESSING_INTERVAL", zap.Error(err))
	}

	return &ProcessingSettings{
		AgentProcessingInterval:                     agentProcessingInterval,
		FacilitatorTickInterval:                     facilitatorTickInterval,
		ConversationForceConcludeAfter:              conversationForceConcludeAfter,
		ConversationForceConcludeByFacilitatorAfter: conversationForceConcludeByFacilitatorAfter,
	}
}
