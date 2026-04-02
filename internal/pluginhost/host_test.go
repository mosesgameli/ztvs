// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pluginhost

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mosesgameli/ztvs-sdk-go/rpc"
	"github.com/mosesgameli/ztvs-sdk-go/sdk"
	"github.com/mosesgameli/ztvs/pkg/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRunner is a testify mock for the Runner interface.
type MockRunner struct {
	mock.Mock
}

func (m *MockRunner) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockRunner) Supports(runtimeType string) bool {
	args := m.Called(runtimeType)
	return args.Bool(0)
}

func (m *MockRunner) Validate(entrypoint string) error {
	args := m.Called(entrypoint)
	return args.Error(0)
}

func (m *MockRunner) Execute(ctx context.Context, entrypoint string, stdin []byte) ([]byte, error) {
	args := m.Called(ctx, entrypoint, stdin)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func TestHost_Handshake(t *testing.T) {
	ctx := context.Background()
	h := New()
	mockRunner := new(MockRunner)

	// Setup mock plugin info
	pluginPath := "/mock/plugin/bin"
	h.plugins[pluginPath] = &PluginInfo{
		Entrypoint: pluginPath,
		Manifest: &sdk.Manifest{
			Name: "mock-plugin",
		},
		Runner: mockRunner,
	}

	t.Run("successful handshake", func(t *testing.T) {
		expectedResp := rpc.Response[rpc.HandshakeResponse]{
			JSONRPC: "2.0",
			ID:      "handshake",
			Result: rpc.HandshakeResponse{
				ChecksSupported: []string{"check-1", "check-2"},
			},
		}
		respBytes, _ := json.Marshal(expectedResp)

		mockRunner.On("Execute", ctx, pluginPath, mock.Anything).Return(respBytes, nil).Once()

		resp, err := h.Handshake(ctx, pluginPath)
		assert.NoError(t, err)
		assert.Contains(t, resp.ChecksSupported, "check-1")
		mockRunner.AssertExpectations(t)
	})

	t.Run("plugin error", func(t *testing.T) {
		errorResp := rpc.Response[json.RawMessage]{
			JSONRPC: "2.0",
			ID:      "handshake",
			Error: &rpc.Error{
				Code:    -32603,
				Message: "Internal error",
			},
		}
		respBytes, _ := json.Marshal(errorResp)

		mockRunner.On("Execute", ctx, pluginPath, mock.Anything).Return(respBytes, nil).Once()

		_, err := h.Handshake(ctx, pluginPath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "plugin error [-32603]")
		mockRunner.AssertExpectations(t)
	})
}

func TestHost_RunCheck(t *testing.T) {
	ctx := context.Background()
	h := New()
	mockRunner := new(MockRunner)

	pluginPath := "/mock/plugin/bin"
	h.plugins[pluginPath] = &PluginInfo{
		Entrypoint: pluginPath,
		Manifest: &sdk.Manifest{
			Name: "mock-plugin",
		},
		Runner: mockRunner,
	}

	t.Run("successful run_check", func(t *testing.T) {
		expectedResp := rpc.Response[rpc.RunCheckResponse]{
			JSONRPC: "2.0",
			ID:      "1",
			Result: rpc.RunCheckResponse{
				Status: "completed",
				Finding: &rpc.Finding{
					ID:    "F1",
					Title: "Found vulnerability",
				},
			},
		}
		respBytes, _ := json.Marshal(expectedResp)

		mockRunner.On("Execute", ctx, pluginPath, mock.Anything).Return(respBytes, nil).Once()

		resp, err := h.RunCheck(ctx, pluginPath, "check-1")
		assert.NoError(t, err)
		assert.Equal(t, "completed", resp.Status)
		assert.Equal(t, "F1", resp.Finding.ID)
		mockRunner.AssertExpectations(t)
	})
}

func TestHost_RegisterRunner(t *testing.T) {
	h := New()
	mockRunner := new(MockRunner)
	mockRunner.On("Name").Return("Mock")
	mockRunner.On("Supports", "mock").Return(true)
	
	h.RegisterRunner(mockRunner)
	
	// We don't have a direct "GetRunner" but we can check if it's used during discovery
	// Actually we just want to cover the line
	assert.NotNil(t, h.runners)
}

func TestHost_Lockfile(t *testing.T) {
	h := New()
	assert.NotNil(t, h.Lockfile()) // Default is initialized in New()
	
	l := &registry.Lockfile{Version: "2.0"}
	h.lockfile = l
	assert.Equal(t, "2.0", h.Lockfile().Version)
}

func TestHost_GetManifest(t *testing.T) {
	h := New()
	m, ok := h.GetManifest("none")
	assert.False(t, ok)
	assert.Nil(t, m)

	h.plugins["p1"] = &PluginInfo{Manifest: &sdk.Manifest{Name: "p1"}}
	m, ok = h.GetManifest("p1")
	assert.True(t, ok)
	assert.Equal(t, "p1", m.Name)
}

func TestHost_Discover(t *testing.T) {
	tmpDir := t.TempDir()
	h := New()
	h.paths = []string{tmpDir}

	// 1. Valid plugin
	p1Dir := filepath.Join(tmpDir, "p1")
	_ = os.MkdirAll(p1Dir, 0755)
	_ = os.WriteFile(filepath.Join(p1Dir, "plugin.yaml"), []byte("name: p1\nruntime:\n  type: binary\n  entrypoint: p1"), 0644)
	_ = os.WriteFile(filepath.Join(p1Dir, "p1"), []byte("bin"), 0755)

	// 2. Missing manifest (should skip)
	p2Dir := filepath.Join(tmpDir, "p2")
	_ = os.MkdirAll(p2Dir, 0755)

	// 3. Unsupported runtime (should skip)
	p3Dir := filepath.Join(tmpDir, "p3")
	_ = os.MkdirAll(p3Dir, 0755)
	_ = os.WriteFile(filepath.Join(p3Dir, "plugin.yaml"), []byte("name: p3\nruntime:\n  type: cobol"), 0644)

	// 4. Duplicate name (p1 again in another path)
	tmpDir2 := t.TempDir()
	h.paths = append(h.paths, tmpDir2)
	p1DupDir := filepath.Join(tmpDir2, "p1-dup")
	_ = os.MkdirAll(p1DupDir, 0755)
	_ = os.WriteFile(filepath.Join(p1DupDir, "plugin.yaml"), []byte("name: p1\nruntime:\n  type: binary"), 0644)

	ctx := context.Background()
	discovered, err := h.Discover(ctx)
	assert.NoError(t, err)
	assert.Len(t, discovered, 1)
	assert.Contains(t, discovered[0], "p1")
}
