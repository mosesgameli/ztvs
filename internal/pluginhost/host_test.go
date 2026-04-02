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
	"testing"

	"github.com/mosesgameli/ztvs-sdk-go/rpc"
	"github.com/mosesgameli/ztvs-sdk-go/sdk"
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
