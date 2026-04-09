package use_cases_test

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
)

func generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest() *requests.AgentByCollectionTokenIdRequest {
	authenticatedWallet := common_tests.RandomEthAddress()
	collectionAddress := common_tests.RandomEthAddress()

	return &requests.AgentByCollectionTokenIdRequest{
		AgentCollectionTokenId: 1,
		AgentCollectionAddress: collectionAddress.Hex(),
		AuthenticatedWallet:    &authenticatedWallet,
	}
}
