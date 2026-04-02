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

func TestCallRPC_Errors(t *testing.T) {
	ctx := context.Background()
	h := New()
	mockRunner := new(MockRunner)
	
	pluginPath := "/mock/plugin"
	h.plugins[pluginPath] = &PluginInfo{
		Entrypoint: pluginPath,
		Manifest:   &sdk.Manifest{Name: "test"},
		Runner:     mockRunner,
	}

	req := rpc.Request{Method: "testMethod"}
	var result struct{ Status string }

	t.Run("InfoNotFound", func(t *testing.T) {
		err := h.callRPC(ctx, "missing-plugin", req, &result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "plugin info not found")
	})

	t.Run("Execute_Error", func(t *testing.T) {
		mockRunner.On("Execute", ctx, pluginPath, mock.Anything).Return(nil, assert.AnError).Once()
		err := h.callRPC(ctx, pluginPath, req, &result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "execute plugin /mock/plugin")
	})

	t.Run("Invalid_JSON_Response", func(t *testing.T) {
		mockRunner.On("Execute", ctx, pluginPath, mock.Anything).Return([]byte("not { json"), nil).Once()
		err := h.callRPC(ctx, pluginPath, req, &result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unmarshal envelope from /mock/plugin")
	})

	t.Run("Plugin_Error_Response", func(t *testing.T) {
		resp := `{"jsonrpc":"2.0","id":"1","error":{"code":-1,"message":"bad"}}`
		mockRunner.On("Execute", ctx, pluginPath, mock.Anything).Return([]byte(resp), nil).Once()
		err := h.callRPC(ctx, pluginPath, req, &result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "plugin error [-1]")
	})

	t.Run("Unmarshal_Result_Error", func(t *testing.T) {
		resp := rpc.Response[json.RawMessage]{
			JSONRPC: "2.0",
			ID:      "1",
			Result:  json.RawMessage(`{"status": 123}`), // status should be string
		}
		respBytes, _ := json.Marshal(resp)
		mockRunner.On("Execute", ctx, pluginPath, mock.Anything).Return(respBytes, nil).Once()
		err := h.callRPC(ctx, pluginPath, req, &result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unmarshal result from /mock/plugin")
	})
}
