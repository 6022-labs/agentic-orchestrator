package container_builders

import (
	"context"
	"log"
	"sync"

	vault "github.com/hashicorp/vault/api"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	vaultContainer     testcontainers.Container
	vaultToken         string
	vaultHost          string
	vaultPort          string
	vaultContainerMu   sync.Mutex
	vaultConnectionsMu sync.Mutex
	vaultConnections   []*vault.Client
)

func startVault() (string, string, string) {
	vaultContainerMu.Lock()
	defer vaultContainerMu.Unlock()

	if vaultContainer == nil {

		vaultToken = "wBGKZCJjfIGqIKB3HMmvcQVt"

		ctx := context.Background()
		req := testcontainers.ContainerRequest{
			Image:        "hashicorp/vault:latest",
			ExposedPorts: []string{"8200/tcp"},
			Env: map[string]string{
				"VAULT_DEV_ROOT_TOKEN_ID": vaultToken,
			},
			WaitingFor: wait.ForAll(
				wait.ForListeningPort("8200/tcp"),
				wait.ForLog("Vault server started!"),
			),
		}

		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			log.Fatalf("Failed to start MySQL container: %v", err)
		}

		vaultHost, _ = container.Host(ctx)
		port, _ := container.MappedPort(ctx, "8200")
		vaultPort = port.Port()

		// Mount the path secret with kvV1
		cmd := []string{
			"/bin/sh", "-c",
			"VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=" + vaultToken + " vault secrets disable secret && " +
				"VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=" + vaultToken + " vault secrets enable -version=1 -path=secret kv",
		}
		exec, _, err := container.Exec(ctx, cmd)
		if err != nil {
			log.Fatalf("Failed to mount secret path: %v", err)
		}
		if exec != 0 {
			log.Fatalf("Failed to mount secret path, exit code: %d", exec)
		}

		log.Println("Vault container started and ready")
	}

	return vaultToken, vaultHost, vaultPort
}

func stopVault() {
	vaultContainerMu.Lock()
	defer vaultContainerMu.Unlock()

	if vaultContainer != nil {
		vaultContainer.Terminate(context.Background())
		vaultContainer = nil

		log.Println("Vault container stopped")
	}
}

func GetVaultConnection() *vault.Client {
	vaultToken, vaultHost, vaultPort := startVault()

	client, err := vault.NewClient(&vault.Config{
		Address: "http://" + vaultHost + ":" + vaultPort,
	})
	if err != nil {
		log.Fatalf("Failed to create Vault client: %v", err)
	}

	client.SetToken(vaultToken)

	vaultConnectionsMu.Lock()
	vaultConnections = append(vaultConnections, client)
	vaultConnectionsMu.Unlock()

	return client
}

func CloseVaultConnection(client *vault.Client) {
	// Just remove from the list, the container will be destroyed after all tests
	for i, connection := range vaultConnections {
		if connection == client {
			vaultConnectionsMu.Lock()
			vaultConnections = append(vaultConnections[:i], vaultConnections[i+1:]...)
			vaultConnectionsMu.Unlock()
			break
		}
	}

	log.Printf("Vault client removed from tracking")

	// Stop the container if no more databases are in use
	if len(vaultConnections) == 0 {
		stopVault()
	}
}
