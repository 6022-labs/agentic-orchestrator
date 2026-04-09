package models

type JsonRPCResponse[T any] struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      interface{}   `json:"id"`
	Result  T             `json:"result"`
	Error   *JsonRPCError `json:"error,omitempty"`
}
