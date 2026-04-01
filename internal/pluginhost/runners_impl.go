package pluginhost

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// BinaryRunner executes native binaries (Go, Rust).
type BinaryRunner struct{}

func (r *BinaryRunner) Name() string { return "Binary" }
func (r *BinaryRunner) Supports(runtimeType string) bool {
	return runtimeType == "binary"
}

func (r *BinaryRunner) Validate(entrypoint string) error {
	info, err := os.Stat(entrypoint)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("entrypoint is a directory: %s", entrypoint)
	}
	// Check for executable bit (owner/group/others)
	if info.Mode()&0111 == 0 {
		return fmt.Errorf("entrypoint is not executable: %s", entrypoint)
	}
	return nil
}

func (r *BinaryRunner) Execute(ctx context.Context, entrypoint string, stdin []byte) ([]byte, error) {
	cmd := exec.CommandContext(ctx, entrypoint, "--rpc")
	cmd.Stdin = bytes.NewReader(stdin)

	// Apply process isolation
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	return cmd.CombinedOutput()
}

// PythonRunner is a placeholder for Python support.
type PythonRunner struct{}

func (r *PythonRunner) Name() string { return "Python" }
func (r *PythonRunner) Supports(runtimeType string) bool {
	return runtimeType == "python"
}

func (r *PythonRunner) Validate(entrypoint string) error {
	if _, err := os.Stat(entrypoint); err != nil {
		return fmt.Errorf("Python script not found: %s", entrypoint)
	}
	// Check for python3 in PATH
	if _, err := exec.LookPath("python3"); err != nil {
		return fmt.Errorf("Python 3 runtime not found in PATH")
	}
	return nil
}

func (r *PythonRunner) Execute(ctx context.Context, entrypoint string, stdin []byte) ([]byte, error) {
	return nil, fmt.Errorf("Python runtime is not yet supported in this version of ZTVS")
}

// NodeRunner is a placeholder for Node.js support.
type NodeRunner struct{}

func (r *NodeRunner) Name() string { return "Node.js" }
func (r *NodeRunner) Supports(runtimeType string) bool {
	return runtimeType == "node"
}

func (r *NodeRunner) Validate(entrypoint string) error {
	if _, err := os.Stat(entrypoint); err != nil {
		return fmt.Errorf("Node.js script not found: %s", entrypoint)
	}
	// Check for node in PATH
	if _, err := exec.LookPath("node"); err != nil {
		return fmt.Errorf("Node.js runtime not found in PATH")
	}
	return nil
}

func (r *NodeRunner) Execute(ctx context.Context, entrypoint string, stdin []byte) ([]byte, error) {
	return nil, fmt.Errorf("Node.js runtime is not yet supported in this version of ZTVS")
}

// JavaRunner is a placeholder for Java support.
type JavaRunner struct{}

func (r *JavaRunner) Name() string { return "Java" }
func (r *JavaRunner) Supports(runtimeType string) bool {
	return runtimeType == "java"
}

func (r *JavaRunner) Validate(entrypoint string) error {
	// For now, Java is purely placeholder
	return fmt.Errorf("Java runtime is not yet supported in this version of ZTVS")
}

func (r *JavaRunner) Execute(ctx context.Context, entrypoint string, stdin []byte) ([]byte, error) {
	return nil, fmt.Errorf("Java runtime is not yet supported in this version of ZTVS")
}
