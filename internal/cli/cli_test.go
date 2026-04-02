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
	"os"
	"path/filepath"
	"testing"

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
}

func TestRunAgent(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	_ = runPluginInit()

	// Use a cancelled context to ensure it returns immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := runAgent(ctx)
	assert.NoError(t, err)
}

func TestRunScan(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	_ = runPluginInit()

	ctx := context.Background()
	// This will try to run a real scan, but with no plugins it should be quick
	err := runScan(ctx)
	assert.NoError(t, err)
}
