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

package cli

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/mosesgameli/ztvs-sdk-go/rpc"
	"github.com/mosesgameli/ztvs-sdk-go/sdk"
	"github.com/mosesgameli/ztvs/internal/config"
	"github.com/mosesgameli/ztvs/internal/pluginhost"
	"github.com/mosesgameli/ztvs/pkg/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRegistry is a mock type for the pluginhost.Registry interface
type MockRegistry struct {
	mock.Mock
}

func (m *MockRegistry) FetchIndex(ctx context.Context) (*registry.Index, error) {
	args := m.Called(ctx)
	return args.Get(0).(*registry.Index), args.Error(1)
}

func (m *MockRegistry) Search(ctx context.Context, query string) ([]registry.PluginMetadata, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]registry.PluginMetadata), args.Error(1)
}

func (m *MockRegistry) GetInfo(ctx context.Context, name string) (*registry.PluginMetadata, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*registry.PluginMetadata), args.Error(1)
}

func (m *MockRegistry) Install(ctx context.Context, name string, host pluginhost.PluginHost) error {
	args := m.Called(ctx, name, host)
	return args.Error(0)
}

func (m *MockRegistry) CheckAndUpdateAll(ctx context.Context, host pluginhost.PluginHost, mode string) error {
	args := m.Called(ctx, host, mode)
	return args.Error(0)
}

func (m *MockRegistry) PerformAtomicUpdate(ctx context.Context, name string, host pluginhost.PluginHost, meta *registry.PluginMetadata) error {
	args := m.Called(ctx, name, host, meta)
	return args.Error(0)
}

// MockHost is a mock type for the pluginhost.PluginHost interface
type MockHost struct {
	mock.Mock
}

func (m *MockHost) Discover(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockHost) Handshake(ctx context.Context, path string) (*rpc.HandshakeResponse, error) {
	args := m.Called(ctx, path)
	return args.Get(0).(*rpc.HandshakeResponse), args.Error(1)
}

func (m *MockHost) RunCheck(ctx context.Context, path string, checkID string) (*rpc.RunCheckResponse, error) {
	args := m.Called(ctx, path, checkID)
	return args.Get(0).(*rpc.RunCheckResponse), args.Error(1)
}

func (m *MockHost) GetManifest(path string) (*sdk.Manifest, bool) {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(*sdk.Manifest), args.Bool(1)
}

func (m *MockHost) GetPluginInfo(path string) (*pluginhost.PluginInfo, bool) {
	args := m.Called(path)
	return args.Get(0).(*pluginhost.PluginInfo), args.Bool(1)
}

func (m *MockHost) Lockfile() *registry.Lockfile {
	args := m.Called()
	return args.Get(0).(*registry.Lockfile)
}

func (m *MockHost) RegisterRunner(r pluginhost.Runner) {
	m.Called(r)
}

func TestRunPluginInit(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	err := runPluginInit()
	assert.NoError(t, err)

	expectedPath := filepath.Join(tmpDir, ".ztvs", "config.yaml")
	assert.FileExists(t, expectedPath)
}

func TestRunPluginList(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	oldHost := host
	host = pluginhost.New()
	defer func() { host = oldHost }()

	err := runPluginList()
	assert.NoError(t, err)
}

func TestRunPluginSearch(t *testing.T) {
	mockReg := new(MockRegistry)
	ctx := mock.Anything
	mockReg.On("Search", ctx, "test").Return([]registry.PluginMetadata{
		{Name: "test-plugin", LatestVersion: "1.0.0", AuditStatus: "verified"},
	}, nil)

	err := runPluginSearch(mockReg, "test")
	assert.NoError(t, err)
	mockReg.AssertExpectations(t)
}

func TestRunPluginInfo(t *testing.T) {
	mockReg := new(MockRegistry)
	ctx := mock.Anything
	mockReg.On("GetInfo", ctx, "test-plugin").Return(&registry.PluginMetadata{
		Name:          "test-plugin",
		LatestVersion: "1.0.0",
		Repo:          "github.com/user/test",
		AuditStatus:   "verified",
		Checksum:      "sha256:123",
	}, nil)

	err := runPluginInfo(mockReg, "test-plugin")
	assert.NoError(t, err)
	mockReg.AssertExpectations(t)
}

func TestRunPluginToggle(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	oldHost := host
	host = pluginhost.New()
	defer func() { host = oldHost }()

	_ = runPluginInit()

	err := runPluginToggle("ghost", true)
	assert.Error(t, err)
}

func TestRunPluginInstall(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	oldHost := host
	host = pluginhost.New()
	defer func() { host = oldHost }()

	mockReg := new(MockRegistry)
	ctx := mock.Anything
	mockReg.On("Install", ctx, "test-plugin", host).Return(nil)

	err := runPluginInstall(mockReg, "test-plugin")
	assert.NoError(t, err)
	mockReg.AssertExpectations(t)
}

func TestRunPluginUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	_ = runPluginInit()

	oldHost := host
	host = pluginhost.New()
	defer func() { host = oldHost }()

	mockReg := new(MockRegistry)
	mockReg.On("CheckAndUpdateAll", mock.Anything, host, "safe").Return(nil)

	err := runPluginUpdate(mockReg)
	assert.NoError(t, err)
	mockReg.AssertExpectations(t)

	// Test update error
	mockReg2 := new(MockRegistry)
	mockReg2.On("CheckAndUpdateAll", mock.Anything, mock.Anything, "safe").Return(assert.AnError)
	err = runPluginUpdate(mockReg2)
	assert.Error(t, err)
}

func TestRunPluginInstall_Error(t *testing.T) {
	mockReg := new(MockRegistry)
	mockReg.On("Install", mock.Anything, "bad-plugin", mock.Anything).Return(assert.AnError)
	err := runPluginInstall(mockReg, "bad-plugin")
	assert.Error(t, err)
}

func TestRunPluginSearch_Error(t *testing.T) {
	mockReg := new(MockRegistry)
	mockReg.On("Search", mock.Anything, "error").Return([]registry.PluginMetadata{}, assert.AnError)
	err := runPluginSearch(mockReg, "error")
	assert.Error(t, err)
}

func TestRunPluginInfo_Error(t *testing.T) {
	mockReg := new(MockRegistry)
	mockReg.On("GetInfo", mock.Anything, "ghost").Return(&registry.PluginMetadata{}, assert.AnError)
	err := runPluginInfo(mockReg, "ghost")
	assert.Error(t, err)
}

func TestRunAgent(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	_ = runPluginInit()

	mockHost := new(MockHost)
	mockReg := new(MockRegistry)
	
	oldHost := host
	host = mockHost
	defer func() { host = oldHost }()
	
	oldReg := registryClient
	registryClient = mockReg
	defer func() { registryClient = oldReg }()

	mockHost.On("Discover", mock.Anything).Return([]string{}, nil)
	mockReg.On("CheckAndUpdateAll", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Use a cancelled context to ensure it returns immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := runAgent(ctx)
	assert.NoError(t, err)

	// Test invalid interval fallback
	cfg, _ := config.Load()
	cfg.Agent.Interval = "invalid"
	_ = cfg.Save()
	
	err = runAgent(ctx)
	assert.NoError(t, err)
}

func TestRunPluginToggle_Errors(t *testing.T) {
	originalHost := host
	mockHost := new(MockHost)
	host = mockHost
	defer func() { host = originalHost }()

	tmpDir := t.TempDir()
	lf := registry.NewLockfile(filepath.Join(tmpDir, "lock.yaml"))
	mockHost.On("Lockfile").Return(lf)

	// Test: Plugin not found in lockfile and not found in Discovery
	mockHost.On("Discover", mock.Anything).Return([]string{}, nil)
	err := runPluginToggle("missing", true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plugin missing not found")

	// Test: Plugin found in Discovery but lockfile save fails
	mockHost.ExpectedCalls = nil
	mockHost.On("Lockfile").Return(lf)
	mockHost.On("Discover", mock.Anything).Return([]string{"p1"}, nil)
	mockHost.On("GetManifest", "p1").Return(&sdk.Manifest{Name: "p1", Version: "1.0"}, true)
	
	// Force Save error by making the directory non-writable
	_ = os.Chmod(tmpDir, 0555) // Read and execute only
	defer func() { _ = os.Chmod(tmpDir, 0755) }()
	
	err = runPluginToggle("p1", true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error saving lockfile")
}

func TestRunPluginList_Errors(t *testing.T) {
	originalHost := host
	mockHost := new(MockHost)
	host = mockHost
	defer func() { host = originalHost }()

	// Test: Discover returns error
	mockHost.On("Discover", mock.Anything).Return([]string{}, errors.New("discovery failed"))
	err := runPluginList()
	assert.Error(t, err)

	// Test: GetManifest returns false
	mockHost.ExpectedCalls = nil
	mockHost.On("Discover", mock.Anything).Return([]string{"/bad-plugin"}, nil)
	mockHost.On("GetManifest", "/bad-plugin").Return(nil, false)
	err = runPluginList()
	assert.NoError(t, err) // It just prints "unknown"
}

func TestRunPluginInstall_ManifestHandling(t *testing.T) {
	originalRegistry := registryClient
	mockReg := new(MockRegistry)
	registryClient = mockReg
	defer func() { registryClient = originalRegistry }()

	tmpDir := t.TempDir()
	_ = os.Setenv("ZTVS_HOME", tmpDir)
	defer func() { _ = os.Unsetenv("ZTVS_HOME") }()

	// Test: Successful install but manifest file missing (no error, just skipped info)
	mockReg.On("Install", mock.Anything, "p1", mock.Anything).Return(nil)
	err := runPluginInstall(mockReg, "p1")
	assert.NoError(t, err)

	// Test: Successful install with valid manifest
	pluginDir := filepath.Join(tmpDir, "plugins", "p2")
	_ = os.MkdirAll(pluginDir, 0755)
	
	yamlData := "name: p2\nversion: 1.0\ncapabilities:\n  - net"
	_ = os.WriteFile(filepath.Join(pluginDir, "plugin.yaml"), []byte(yamlData), 0644)
	
	mockReg.On("Install", mock.Anything, "p2", mock.Anything).Return(nil)
	err = runPluginInstall(mockReg, "p2")
	assert.NoError(t, err)

	// Test: Install error
	mockReg.ExpectedCalls = nil
	mockReg.On("Install", mock.Anything, "p3", mock.Anything).Return(errors.New("install failed"))
	err = runPluginInstall(mockReg, "p3")
	assert.Error(t, err)

	mockReg.ExpectedCalls = nil
	mockReg.On("Install", mock.Anything, "p4", mock.Anything).Return(errors.New("install failed"))
	err = runPluginInstall(mockReg, "p4")
	assert.Error(t, err)

	// Test: Manifest read error (already covers successful Install but missing manifest)
	mockReg.ExpectedCalls = nil
	mockReg.On("Install", mock.Anything, "p5", mock.Anything).Return(nil)
	err = runPluginInstall(mockReg, "p5")
	assert.NoError(t, err) // Implementation returns nil even if manifest read fails
}

func TestRunPluginToggle_Unknown(t *testing.T) {
	mockHost := new(MockHost)
	oldHost := host
	host = mockHost
	defer func() { host = oldHost }()

	mockHost.On("Discover", mock.Anything).Return([]string{"p1"}, nil)
	mockHost.On("GetManifest", "p1").Return(&sdk.Manifest{Name: "p1"}, true)
	mockHost.On("Lockfile").Return(registry.NewLockfile(filepath.Join(t.TempDir(), "lock")))
	err := runPluginToggle("unknown-plugin", true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown-plugin not found")

	// Test: Discover fails
	mockHost.ExpectedCalls = nil
	mockHost.On("Lockfile").Return(registry.NewLockfile(filepath.Join(t.TempDir(), "lock")))
	mockHost.On("Discover", mock.Anything).Return([]string{}, errors.New("discover failed"))
	err = runPluginToggle("any", true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "discover failed")
}

func TestRunPluginList_Error(t *testing.T) {
	mockHost := new(MockHost)
	oldHost := host
	host = mockHost
	defer func() { host = oldHost }()

	mockHost.On("Discover", mock.Anything).Return([]string{}, errors.New("list fail"))
	err := runPluginList()
	assert.Error(t, err)
}

func TestRunPluginList_Success(t *testing.T) {
	mockHost := new(MockHost)
	oldHost := host
	host = mockHost
	defer func() { host = oldHost }()

	mockHost.On("Discover", mock.Anything).Return([]string{"p1"}, nil)
	mockHost.On("GetManifest", "p1").Return(&sdk.Manifest{Name: "p1", Version: "1.0.0"}, true)
	mockHost.On("Lockfile").Return(registry.NewLockfile(filepath.Join(t.TempDir(), "lock")))

	err := runPluginList()
	assert.NoError(t, err)
}

func TestRoot_Commands(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	mockHost := new(MockHost)
	mockReg := new(MockRegistry)
	
	oldHost := host
	host = mockHost
	defer func() { host = oldHost }()
	
	oldReg := registryClient
	registryClient = mockReg
	defer func() { registryClient = oldReg }()

	mockHost.On("Discover", mock.Anything).Return([]string{}, nil)
	mockHost.On("Lockfile").Return(registry.NewLockfile(filepath.Join(tmpDir, "lock")))
	mockReg.On("Search", mock.Anything, mock.Anything).Return([]registry.PluginMetadata{}, nil)
	mockReg.On("GetInfo", mock.Anything, mock.Anything).Return(&registry.PluginMetadata{}, nil)
	mockReg.On("Install", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockReg.On("CheckAndUpdateAll", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	tests := [][]string{
		{"plugin", "init"},
		{"plugin", "list"},
		{"plugin", "search", "test"},
		{"plugin", "info", "test"},
		{"plugin", "enable", "test"},
		{"plugin", "disable", "test"},
		{"plugin", "install", "test"},
		{"plugin", "update"},
		{"scan"},
		{"agent"},
	}

	for _, args := range tests {
		rootCmd.SetArgs(args)
		// We use a cancelled context for agent to avoid hangs
		ctx, cancel := context.WithCancel(context.Background())
		if args[0] == "agent" {
			cancel()
		}
		_ = rootCmd.ExecuteContext(ctx)
		cancel()
	}
}

func TestRunScan_Mocks(t *testing.T) {
	oldHost := host
	mockHost := new(MockHost)
	host = mockHost
	defer func() { host = oldHost }()

	oldReg := registryClient
	mockReg := new(MockRegistry)
	registryClient = mockReg
	defer func() { registryClient = oldReg }()

	mockHost.On("Discover", mock.Anything).Return([]string{}, nil)
	mockReg.On("CheckAndUpdateAll", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	
	err := runScan(context.Background())
	assert.NoError(t, err)

	// Test: engine.New error (via unreadable config file)
	tmpHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", originalHome)
	
	configDir := filepath.Join(tmpHome, ".ztvs")
	os.MkdirAll(configDir, 0755)
	// Create a directory where the file should be to make ReadFile fail
	os.MkdirAll(filepath.Join(configDir, "config.yaml"), 0755)
	
	err = runScan(context.Background())
	assert.Error(t, err)

	// Test: CheckAndUpdateAll failure (should be caught and logged as warning, not return error)
	os.Setenv("HOME", tmpHome) // Reset to avoid config error first
	os.RemoveAll(filepath.Join(tmpHome, ".ztvs", "config.yaml"))
	mockReg.ExpectedCalls = nil
	mockReg.On("CheckAndUpdateAll", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("update failed"))
	err = runScan(context.Background())
	assert.NoError(t, err)
}
