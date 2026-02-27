package clients

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector_http_mcp/models"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const (
	listToolsMethod = "tools/list"
	callToolMethod  = "tools/call"
)

type McpClientInterface interface {
	ListTools(mcpServerUrl string, headers map[string]string) (*models.JsonRPCResponse[models.ExternalMcpToolsResponse], error)
	CallTool(mcpServerUrl string, toolName string, input map[string]any, headers map[string]string) (*models.JsonRPCResponse[json.RawMessage], error)
}

type McpClient struct {
	logger     *zap.Logger
	httpClient *http.Client
}

type newMcpClientParams struct {
	dig.In

	Logger     *zap.Logger
	HttpClient *http.Client `name:"mcp-http-client"`
}

func NewMcpClient(
	params newMcpClientParams,
) *McpClient {
	return &McpClient{
		logger:     params.Logger,
		httpClient: params.HttpClient,
	}
}

func (c *McpClient) ListTools(mcpServerUrl string, headers map[string]string) (*models.JsonRPCResponse[models.ExternalMcpToolsResponse], error) {
	requestBody, err := json.Marshal(models.JsonRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  listToolsMethod,
		Params:  map[string]any{},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal MCP list tools request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, mcpServerUrl, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP list tools request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	c.logger.Debug("Sending MCP list tools request", zap.String("url", mcpServerUrl))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute MCP list tools request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > http.StatusPartialContent {
		respBody, _ := io.ReadAll(resp.Body)
		c.logger.Warn("MCP list tools returned unexpected status code",
			zap.String("url", mcpServerUrl),
			zap.Int("status_code", resp.StatusCode),
			zap.ByteString("response_body", respBody),
		)

		return nil, fmt.Errorf("unexpected MCP list tools status code: %d", resp.StatusCode)
	}

	rpcResponse, err := decodeRPCResponse[models.ExternalMcpToolsResponse](c.logger, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode MCP list tools response: %w", err)
	}

	return rpcResponse, nil
}

func (c *McpClient) CallTool(mcpServerUrl string, toolName string, input map[string]any, headers map[string]string) (*models.JsonRPCResponse[json.RawMessage], error) {
	requestBody, err := json.Marshal(models.JsonRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  callToolMethod,
		Params: map[string]any{
			"name":      toolName,
			"arguments": input,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal MCP call tool request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, mcpServerUrl, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP call tool request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	c.logger.Debug("Sending MCP call tool request", zap.String("url", mcpServerUrl), zap.String("tool", toolName))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute MCP call tool request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > http.StatusPartialContent {
		respBody, _ := io.ReadAll(resp.Body)
		c.logger.Warn("MCP call tool returned unexpected status code",
			zap.String("url", mcpServerUrl),
			zap.Int("status_code", resp.StatusCode),
			zap.ByteString("response_body", respBody),
		)

		return nil, fmt.Errorf("unexpected MCP call tool status code: %d", resp.StatusCode)
	}

	rpcResponse, err := decodeRPCResponse[json.RawMessage](c.logger, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode MCP call tool response: %w", err)
	}

	return rpcResponse, nil
}

func decodeRPCResponse[T any](logger *zap.Logger, resp *http.Response) (*models.JsonRPCResponse[T], error) {
	contentType := resp.Header.Get("Content-Type")

	if strings.HasPrefix(strings.ToLower(contentType), "text/event-stream") {
		logger.Debug("Decoding MCP response as SSE stream")
		return decodeEventStream[T](logger, resp.Body)
	}

	return decodeApplicationJSON[T](resp.Body)
}

func decodeApplicationJSON[T any](body io.Reader) (*models.JsonRPCResponse[T], error) {
	var rpcResponse models.JsonRPCResponse[T]
	if err := json.NewDecoder(body).Decode(&rpcResponse); err != nil {
		return nil, err
	}
	return &rpcResponse, nil
}

const (
	sseDataPrefix = "data:"
)

func decodeEventStream[T any](logger *zap.Logger, body io.Reader) (*models.JsonRPCResponse[T], error) {
	reader := bufio.NewReader(body)

	var (
		payloadBuilder strings.Builder
		lastDecoded    *models.JsonRPCResponse[T]
	)

	flushPayload := func() {
		if payloadBuilder.Len() == 0 {
			return
		}

		var candidate models.JsonRPCResponse[T]
		if err := json.Unmarshal([]byte(payloadBuilder.String()), &candidate); err != nil {
			logger.Debug("SSE event did not decode to JsonRPCResponse", zap.Error(err))
			payloadBuilder.Reset()
			return
		}

		payloadBuilder.Reset()
		lastDecoded = &candidate
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("failed reading SSE stream: %w", err)
		}

		line = strings.TrimRight(line, "\r\n")

		if strings.TrimSpace(line) == "" {
			flushPayload()
		} else if strings.HasPrefix(line, sseDataPrefix) {
			data := strings.TrimSpace(strings.TrimPrefix(line, sseDataPrefix))
			if payloadBuilder.Len() > 0 {
				payloadBuilder.WriteByte('\n')
			}
			payloadBuilder.WriteString(data)
		}

		if errors.Is(err, io.EOF) {
			break
		}
	}

	flushPayload()

	if lastDecoded == nil {
		return nil, fmt.Errorf("no JSON-RPC payload found in SSE stream")
	}

	return lastDecoded, nil
}
