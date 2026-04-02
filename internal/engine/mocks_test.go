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
	"io"

	"github.com/mosesgameli/ztvs-sdk-go/rpc"
	"github.com/mosesgameli/ztvs-sdk-go/sdk"
	"github.com/mosesgameli/ztvs/internal/pluginhost"
	"github.com/mosesgameli/ztvs/pkg/registry"
	"github.com/stretchr/testify/mock"
)

// MockRunner
type MockRunner struct {
	mock.Mock
}

func (m *MockRunner) Name() string { return "mock" }
func (m *MockRunner) Supports(runtime string) bool { return runtime == "mock" }
func (m *MockRunner) Validate(path string) error {
	args := m.Called(path)
	return args.Error(0)
}
func (m *MockRunner) Execute(ctx context.Context, path string, payload []byte) ([]byte, error) {
	args := m.Called(ctx, path, payload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

// MockReporter
type MockReporter struct {
	mock.Mock
}

func (m *MockReporter) AddFinding(pluginName string, finding *rpc.Finding) {
	m.Called(pluginName, finding)
}
func (m *MockReporter) Flush() error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockReporter) SetOutput(w io.Writer) {
	m.Called(w)
}

// MockHost
type MockHost struct {
	mock.Mock
}

func (m *MockHost) Discover(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}
func (m *MockHost) GetManifest(pluginPath string) (*sdk.Manifest, bool) {
	args := m.Called(pluginPath)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(*sdk.Manifest), args.Bool(1)
}
func (m *MockHost) Handshake(ctx context.Context, pluginPath string) (*rpc.HandshakeResponse, error) {
	args := m.Called(ctx, pluginPath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rpc.HandshakeResponse), args.Error(1)
}
func (m *MockHost) RunCheck(ctx context.Context, pluginPath string, checkID string) (*rpc.RunCheckResponse, error) {
	args := m.Called(ctx, pluginPath, checkID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rpc.RunCheckResponse), args.Error(1)
}
func (m *MockHost) RegisterRunner(r pluginhost.Runner) {
	m.Called(r)
}
func (m *MockHost) Lockfile() *registry.Lockfile {
	args := m.Called()
	return args.Get(0).(*registry.Lockfile)
}
func (m *MockHost) GetPluginInfo(pluginPath string) (*pluginhost.PluginInfo, bool) {
	args := m.Called(pluginPath)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(*pluginhost.PluginInfo), args.Bool(1)
}

// MockRegistry
type MockRegistry struct {
	mock.Mock
}

func (m *MockRegistry) FetchIndex(ctx context.Context) (*registry.Index, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*registry.Index), args.Error(1)
}

func (m *MockRegistry) Search(ctx context.Context, query string) ([]registry.PluginMetadata, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]registry.PluginMetadata), args.Error(1)
}

func (m *MockRegistry) GetInfo(ctx context.Context, name string) (*registry.PluginMetadata, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*registry.PluginMetadata), args.Error(1)
}

func (m *MockRegistry) Install(ctx context.Context, pluginName string, host pluginhost.PluginHost) error {
	args := m.Called(ctx, pluginName, host)
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
