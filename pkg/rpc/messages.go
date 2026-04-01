package rpc

type Request struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type Response[T any] struct {
	JSONRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  T      `json:"result,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type HandshakeRequest struct {
	HostVersion string `json:"host_version"`
	APIVersion  int    `json:"api_version"`
}

type HandshakeResponse struct {
	Name            string   `json:"name"`
	Version         string   `json:"version"`
	APIVersion      int      `json:"api_version"`
	ChecksSupported []string `json:"checks_supported"`
}

type RunCheckRequest struct {
	CheckID string `json:"check_id"`
}

type RunCheckResponse struct {
	Status  string   `json:"status"`
	Finding *Finding `json:"finding"`
}

type Finding struct {
	ID          string                 `json:"id"`
	Severity    string                 `json:"severity"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Evidence    map[string]interface{} `json:"evidence"`
	Remediation string                 `json:"remediation"`
}
