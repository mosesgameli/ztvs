package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mosesgameli/ztvs/pkg/rpc"
)

func Run(meta Metadata, checks []Check) {
	var req rpc.Request
	err := json.NewDecoder(os.Stdin).Decode(&req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode request: %v\n", err)
		os.Exit(1)
	}

	switch req.Method {
	case "handshake":
		ids := make([]string, len(checks))
		for i, c := range checks {
			ids[i] = c.ID()
		}

		resp := rpc.Response[rpc.HandshakeResponse]{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: rpc.HandshakeResponse{
				Name:            meta.Name,
				Version:         meta.Version,
				APIVersion:      meta.APIVersion,
				ChecksSupported: ids,
			},
		}
		_ = json.NewEncoder(os.Stdout).Encode(resp)

	case "run_check":
		// Handle params unmarshaling
		var runReq rpc.RunCheckRequest
		b, _ := json.Marshal(req.Params)
		_ = json.Unmarshal(b, &runReq)

		for _, c := range checks {
			if c.ID() == runReq.CheckID {
				finding, err := c.Run(context.Background())
				if err != nil {
					sendError(req.ID, 5000, err.Error())
					return
				}

				resp := rpc.Response[rpc.RunCheckResponse]{
					JSONRPC: "2.0",
					ID:      req.ID,
					Result: rpc.RunCheckResponse{
						Status: "fail", // Defaulting to fail if a finding is returned
						Finding: &rpc.Finding{
							ID:          finding.ID,
							Severity:    finding.Severity,
							Title:       finding.Title,
							Description: finding.Description,
							Evidence:    finding.Evidence,
							Remediation: finding.Remediation,
						},
					},
				}

				_ = json.NewEncoder(os.Stdout).Encode(resp)
				return
			}
		}
		sendError(req.ID, 4002, "check not found")
	default:
		sendError(req.ID, -32601, "method not found")
	}
}

func sendError(id string, code int, msg string) {
	resp := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"error": map[string]interface{}{
			"code":    code,
			"message": msg,
		},
	}
	_ = json.NewEncoder(os.Stdout).Encode(resp)
}
