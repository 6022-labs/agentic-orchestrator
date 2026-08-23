package container_builders

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	AnvilChainId     = "31337"
	AnvilChainIdUint = uint64(31337)
)

// AnvilAccountSettings is a throwaway chain endpoint pair plus a pre-funded key
// owned by a single test.
type AnvilAccountSettings struct {
	ChainId       string
	ChainIdUint   uint64
	HttpUrl       string
	WsUrl         string
	PrivateKeyHex string
	Address       common.Address
}

var (
	anvilContainer   testcontainers.Container
	anvilHttpUrl     string
	anvilWsUrl       string
	anvilRpcClient   *rpc.Client
	anvilContainerMu sync.Mutex
)

// Start the container once for all test packages
func startAnvil() *rpc.Client {
	anvilContainerMu.Lock()
	defer anvilContainerMu.Unlock()

	if anvilContainer == nil {
		ctx := context.Background()
		req := testcontainers.ContainerRequest{
			Image:        "ghcr.io/foundry-rs/foundry:v1.7.1",
			ExposedPorts: []string{"8545/tcp"},
			// The image entrypoint is a shell, so anvil has to be the explicit entrypoint.
			Entrypoint: []string{"anvil", "--host", "0.0.0.0"},
			WaitingFor: wait.ForListeningPort("8545/tcp"),
		}

		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			log.Fatalf("Failed to start anvil container: %v", err)
		}
		anvilContainer = container

		host, _ := container.Host(ctx)
		port, _ := container.MappedPort(ctx, "8545")
		anvilHttpUrl = fmt.Sprintf("http://%s:%s", host, port.Port())
		anvilWsUrl = fmt.Sprintf("ws://%s:%s", host, port.Port())

		client, err := rpc.DialContext(ctx, anvilHttpUrl)
		if err != nil {
			log.Fatalf("Failed to dial anvil at %s: %v", anvilHttpUrl, err)
		}
		var chainId string
		if err := client.CallContext(ctx, &chainId, "eth_chainId"); err != nil {
			log.Fatalf("Anvil did not answer eth_chainId: %v", err)
		}
		anvilRpcClient = client

		log.Println("Anvil container started and ready")
	}

	return anvilRpcClient
}

// Stop the container after all tests
func stopAnvil() {
	anvilContainerMu.Lock()
	defer anvilContainerMu.Unlock()

	if anvilContainer != nil {
		if anvilRpcClient != nil {
			anvilRpcClient.Close()
			anvilRpcClient = nil
		}
		anvilContainer.Terminate(context.Background())
		anvilContainer = nil

		log.Println("Anvil container stopped")
	}
}

// CreateAnvilAccount funds a freshly generated key so parallel tests never share
// a nonce sequence (anvil_setBalance is nonce-free, unlike funding from a faucet).
func CreateAnvilAccount() AnvilAccountSettings {
	client := startAnvil()

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("Failed to generate anvil test key: %v", err)
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	balance := new(big.Int).Exp(big.NewInt(10), big.NewInt(22), nil)
	if err := client.CallContext(
		context.Background(),
		nil,
		"anvil_setBalance",
		address,
		hexutil.EncodeBig(balance),
	); err != nil {
		log.Fatalf("Failed to fund %s on anvil: %v", address.Hex(), err)
	}

	return AnvilAccountSettings{
		ChainId:       AnvilChainId,
		ChainIdUint:   AnvilChainIdUint,
		HttpUrl:       anvilHttpUrl,
		WsUrl:         anvilWsUrl,
		PrivateKeyHex: common.Bytes2Hex(crypto.FromECDSA(privateKey)),
		Address:       address,
	}
}
