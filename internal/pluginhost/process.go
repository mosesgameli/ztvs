package pluginhost

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"syscall"

	"github.com/mosesgameli/ztvs/pkg/rpc"
)

func (h *Host) Handshake(
	ctx context.Context,
	pluginPath string,
) (*rpc.HandshakeResponse, error) {

	req := rpc.Request{
		JSONRPC: "2.0",
		ID:      "handshake",
		Method:  "handshake",
		Params: rpc.HandshakeRequest{
			HostVersion: "1.0.0",
			APIVersion:  1,
		},
	}

	var resp rpc.HandshakeResponse
	err := h.callRPC(ctx, pluginPath, req, &resp)
	return &resp, err
}

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

	var resp rpc.RunCheckResponse
	err := h.callRPC(ctx, pluginPath, req, &resp)
	return &resp, err
}

func (h *Host) callRPC(
	ctx context.Context,
	pluginPath string,
	req rpc.Request,
	result interface{},
) error {
	// 1. Verify integrity before every call (Phase 3)
	if manifest, ok := h.GetManifest(pluginPath); ok && manifest.Checksum != "" {
		if err := VerifyIntegrity(pluginPath, manifest.Checksum); err != nil {
			return fmt.Errorf("security violation: %v", err)
		}
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %v", err)
	}

	cmd := exec.CommandContext(ctx, pluginPath, "--rpc")
	cmd.Stdin = bytes.NewReader(payload)

	// 2. Apply basic sandboxing (Process Isolation)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true, // Separate process group
	}

	out, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("plugin %s timed out", pluginPath)
		}
		return fmt.Errorf("execute plugin %s: %v", pluginPath, err)
	}

	var r rpc.Response[json.RawMessage]
	if err := json.Unmarshal(out, &r); err != nil {
		return fmt.Errorf("unmarshal envelope from %s: %v", pluginPath, err)
	}

	if r.Error != nil {
		return fmt.Errorf("plugin error [%d]: %s", r.Error.Code, r.Error.Message)
	}

	if err := json.Unmarshal(r.Result, result); err != nil {
		return fmt.Errorf("unmarshal result from %s: %v", pluginPath, err)
	}

	return nil
}
