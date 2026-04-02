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

package engine

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mosesgameli/ztvs-sdk-go/rpc"
	"github.com/mosesgameli/ztvs-sdk-go/sdk"
	"github.com/mosesgameli/ztvs/internal/config"
	"github.com/mosesgameli/ztvs/internal/pluginhost"
	"github.com/mosesgameli/ztvs/internal/policy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEngine_All(t *testing.T) {
	t.Run("Scan_FullSuccess_Interactive", func(t *testing.T) {
		cfg := &config.Config{}
		mockHost := new(MockHost)
		mockReporter := new(MockReporter)
		mockReg := new(MockRegistry)
		e := New(cfg, mockHost, mockReporter, mockReg)
		e.Interactive = true

		mockReg.On("CheckAndUpdateAll", mock.Anything, mockHost, mock.Anything).Return(nil)
		mockHost.On("Discover", mock.Anything).Return([]string{"/p1"}, nil)
		mockHost.On("GetPluginInfo", "/p1").Return(&pluginhost.PluginInfo{Enabled: true}, true)
		mockHost.On("GetManifest", "/p1").Return(&sdk.Manifest{Name: "p1", Capabilities: []string{"read"}}, true)
		mockHost.On("Handshake", mock.Anything, "/p1").Return(&rpc.HandshakeResponse{
			Name: "p1", APIVersion: 1, ChecksSupported: []string{"c1"},
		}, nil)
		mockHost.On("RunCheck", mock.Anything, "/p1", "c1").Return(&rpc.RunCheckResponse{
			Finding: &rpc.Finding{ID: "f1"},
		}, nil)
		mockReporter.On("AddFinding", "p1", mock.Anything).Return()
		mockReporter.On("Flush").Return(nil)

		err := e.Scan()
		assert.NoError(t, err)
	})

	t.Run("Scan_DiscoveryError", func(t *testing.T) {
		cfg := &config.Config{}
		mockHost := new(MockHost)
		mockReg := new(MockRegistry)
		e := New(cfg, mockHost, nil, mockReg)
		e.Interactive = true

		mockReg.On("CheckAndUpdateAll", mock.Anything, mockHost, mock.Anything).Return(nil)
		mockHost.On("Discover", mock.Anything).Return(nil, errors.New("remote error"))

		err := e.Scan()
		assert.Error(t, err)
	})

	t.Run("Scan_NoAllowedPlugins", func(t *testing.T) {
		cfg := &config.Config{}
		cfg.Policy.BlockedCapabilities = []string{"net"}
		mockHost := new(MockHost)
		mockReg := new(MockRegistry)
		mockReporter := new(MockReporter)
		e := New(cfg, mockHost, mockReporter, mockReg)
		e.Interactive = true
		e.stdin = strings.NewReader("n\n") // Deny

		mockReg.On("CheckAndUpdateAll", mock.Anything, mockHost, mock.Anything).Return(nil)
		mockHost.On("Discover", mock.Anything).Return([]string{"/p1"}, nil)
		mockHost.On("GetPluginInfo", "/p1").Return(&pluginhost.PluginInfo{Enabled: true}, true)
		mockHost.On("GetManifest", "/p1").Return(&sdk.Manifest{Name: "p1", Capabilities: []string{"net"}}, true)
		mockReporter.On("Flush").Return(nil)

		err := e.Scan()
		assert.NoError(t, err)
	})

	t.Run("preflightCapabilityCheck_GrantSuccess", func(t *testing.T) {
		tmp := t.TempDir()
		oldHome := os.Getenv("HOME")
		os.Setenv("HOME", tmp)
		defer os.Setenv("HOME", oldHome)

		cfg := config.DefaultConfig()
		cfg.Policy.BlockedCapabilities = []string{"net"}
		mockHost := new(MockHost)
		e := New(cfg, mockHost, nil, nil)
		e.Interactive = true
		e.stdin = strings.NewReader("y\n")

		mockHost.On("GetPluginInfo", "/p1").Return(&pluginhost.PluginInfo{Enabled: true}, true)
		mockHost.On("GetManifest", "/p1").Return(&sdk.Manifest{Name: "p1", Capabilities: []string{"net"}}, true)

		res := e.preflightCapabilityCheck([]string{"/p1"})
		assert.Len(t, res, 1)
		assert.Contains(t, cfg.Policy.AllowedCapabilities, "net")
	})

	t.Run("preflightCapabilityCheck_SaveError", func(t *testing.T) {
		tmpHome := t.TempDir()
		oldHome := os.Getenv("HOME")
		os.Setenv("HOME", tmpHome)
		defer os.Setenv("HOME", oldHome)

		// Create .ztvs as a FILE to force MkdirAll failure
		configDir := filepath.Join(tmpHome, ".ztvs")
		os.WriteFile(configDir, []byte("not-a-dir"), 0644)

		cfg := config.DefaultConfig()
		cfg.Policy.BlockedCapabilities = []string{"net"}
		mockHost := new(MockHost)
		e := New(cfg, mockHost, nil, nil)
		e.Interactive = true
		e.stdin = strings.NewReader("y\n")

		mockHost.On("GetPluginInfo", "/p1").Return(&pluginhost.PluginInfo{Enabled: true}, true)
		mockHost.On("GetManifest", "/p1").Return(&sdk.Manifest{Name: "p1", Capabilities: []string{"net"}}, true)

		res := e.preflightCapabilityCheck([]string{"/p1"})
		assert.Empty(t, res) // Should fail MkdirAll and skip
	})

	t.Run("scanPlugin_RaceAndVersions", func(t *testing.T) {
		cfg := &config.Config{}
		mockHost := new(MockHost)
		mockReporter := new(MockReporter)
		e := &Engine{host: mockHost, policy: policy.New(cfg), cfg: cfg, reporter: mockReporter}

		// 1. Race condition (IsAllowed fails after Discovery)
		mockHost.On("GetManifest", "/p1").Return(&sdk.Manifest{Name: "p1", Capabilities: []string{"net"}}, true)
		cfg.Policy.BlockedCapabilities = []string{"net"}
		e.policy.Reload(cfg)
		e.scanPlugin(context.Background(), "/p1")

		// 2. Handshake Error
		mockHost.On("GetManifest", "/p2").Return(&sdk.Manifest{Name: "p2"}, true)
		cfg.Policy.BlockedCapabilities = nil
		e.policy.Reload(cfg)
		mockHost.On("Handshake", mock.Anything, "/p2").Return(nil, errors.New("handshake fail"))
		e.scanPlugin(context.Background(), "/p2")

		// 3. API Version mismatch
		mockHost.On("GetManifest", "/p3").Return(&sdk.Manifest{Name: "p3"}, true)
		mockHost.On("Handshake", mock.Anything, "/p3").Return(&rpc.HandshakeResponse{APIVersion: 2}, nil)
		e.scanPlugin(context.Background(), "/p3")

		// 4. RunCheck error
		mockHost.On("GetManifest", "/p4").Return(&sdk.Manifest{Name: "p4"}, true)
		mockHost.On("Handshake", mock.Anything, "/p4").Return(&rpc.HandshakeResponse{APIVersion: 1, ChecksSupported: []string{"c1"}}, nil)
		mockHost.On("RunCheck", mock.Anything, "/p4", "c1").Return(nil, errors.New("check fail"))
		e.scanPlugin(context.Background(), "/p4")
	})

	t.Run("RunLoop_Coverage", func(t *testing.T) {
		cfg := &config.Config{}
		mockHost := new(MockHost)
		mockReporter := new(MockReporter)
		mockReg := new(MockRegistry)
		e := New(cfg, mockHost, mockReporter, mockReg)

		mockReg.On("CheckAndUpdateAll", mock.Anything, mockHost, mock.Anything).Return(nil)
		mockHost.On("Discover", mock.Anything).Return([]string{}, nil)
		mockReporter.On("Flush").Return(nil)

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(250 * time.Millisecond) // Wait for at least 2 ticks if interval is 100ms
			cancel()
		}()

		err := e.RunLoop(ctx, 100*time.Millisecond)
		assert.ErrorIs(t, err, context.Canceled)
	})
}
