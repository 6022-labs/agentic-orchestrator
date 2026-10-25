package interfaces

// SignedRequestMessageBuilderInterface renders the canonical text a caller signs
// with EIP-191 to authenticate one request.
type SignedRequestMessageBuilderInterface interface {
	Build(method string, path string, query string, timestamp int64, nonce string, body []byte) string
}
