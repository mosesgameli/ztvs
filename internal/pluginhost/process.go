package pluginhost

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/mosesgameli/ztvs/pkg/rpc"
)

func (h *Host) runCheckProcess(
	ctx context.Context,
	pluginPath string,
	checkID string,
) (*rpc.RunCheckResponse, error) {

	req := rpc.Request{
		JSONRPC: "2.0",
		ID:      "1",
		Method:  "run_check",
		Params: rpc.RunCheckRequest{
			CheckID: checkID,
		},
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %v", err)
	}

	cmd := exec.CommandContext(
		ctx,
		pluginPath,
		"--rpc",
	)

	cmd.Stdin = bytes.NewReader(payload)

	// In a real implementation, we'd handle stderr logs separately
	out, err := cmd.Output()
	if err != nil {
		// Try to capture more context from stderr if needed
		return nil, fmt.Errorf("execute plugin %s: %v", pluginPath, err)
	}

	var resp rpc.Response[rpc.RunCheckResponse]
	err = json.Unmarshal(out, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response from %s: %v", pluginPath, err)
	}

	return &resp.Result, nil
}
