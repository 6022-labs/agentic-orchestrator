package use_cases_test

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
)

func generateValidAgentByCollectionTokenIdRequest() *requests.AgentByCollectionTokenIdRequest {
	return &requests.AgentByCollectionTokenIdRequest{
		AgentCollectionTokenId: 1,
		AgentCollectionAddress: common_tests.RandomEthAddress().Hex(),
	}
}
