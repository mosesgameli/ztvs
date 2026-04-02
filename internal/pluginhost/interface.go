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

	"github.com/mosesgameli/ztvs-sdk-go/rpc"
	"github.com/mosesgameli/ztvs-sdk-go/sdk"
	"github.com/mosesgameli/ztvs/pkg/registry"
)

// PluginHost defines the interface for plugin discovery and execution.
// This supports dependency injection and comprehensive unit testing.
type PluginHost interface {
	Discover(ctx context.Context) ([]string, error)
	GetManifest(pluginPath string) (*sdk.Manifest, bool)
	Handshake(ctx context.Context, pluginPath string) (*rpc.HandshakeResponse, error)
	RunCheck(ctx context.Context, pluginPath string, checkID string) (*rpc.RunCheckResponse, error)
	RegisterRunner(r Runner)
	Lockfile() *registry.Lockfile
	GetPluginInfo(pluginPath string) (*PluginInfo, bool)
}

type Registry interface {
	FetchIndex(ctx context.Context) (*registry.Index, error)
	Search(ctx context.Context, query string) ([]registry.PluginMetadata, error)
	GetInfo(ctx context.Context, name string) (*registry.PluginMetadata, error)
	Install(ctx context.Context, pluginName string, host PluginHost) error
	CheckAndUpdateAll(ctx context.Context, host PluginHost, mode string) error
	PerformAtomicUpdate(ctx context.Context, name string, host PluginHost, meta *registry.PluginMetadata) error
}

// Ensure Host implements PluginHost
var _ PluginHost = (*Host)(nil)
