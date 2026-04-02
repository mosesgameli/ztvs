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

package policy

import (
	"testing"

	"github.com/mosesgameli/ztvs/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestPolicy_IsAllowed(t *testing.T) {
	t.Run("allow all when no policy", func(t *testing.T) {
		p := &Policy{}
		cap, err := p.IsAllowed("test-plugin", []string{"read_files", "network"})
		assert.NoError(t, err)
		assert.Empty(t, cap)
	})

	t.Run("block explicit capability", func(t *testing.T) {
		p := &Policy{
			BlockedCapabilities: []string{"network"},
		}
		cap, err := p.IsAllowed("test-plugin", []string{"read_files", "network"})
		assert.Error(t, err)
		assert.Equal(t, "network", cap)
		assert.Contains(t, err.Error(), "requested blocked capability")
	})

	t.Run("allow only authorized list", func(t *testing.T) {
		p := &Policy{
			AllowedCapabilities: []string{"read_files"},
		}
		// Valid
		cap, err := p.IsAllowed("test-plugin", []string{"read_files"})
		assert.NoError(t, err)
		assert.Empty(t, cap)

		// Unauthorized
		cap, err = p.IsAllowed("test-plugin", []string{"network"})
		assert.Error(t, err)
		assert.Equal(t, "network", cap)
		assert.Contains(t, err.Error(), "requested unauthorized capability")
	})

	t.Run("mixed policy priority", func(t *testing.T) {
		// Blocked takes precedence over Allowed
		p := &Policy{
			AllowedCapabilities: []string{"read_files", "network"},
			BlockedCapabilities: []string{"network"},
		}
		cap, err := p.IsAllowed("test-plugin", []string{"network"})
		assert.Error(t, err)
		assert.Equal(t, "network", cap)
		assert.Contains(t, err.Error(), "requested blocked capability")
	})
}

func TestPolicy_Reload(t *testing.T) {
	p := &Policy{
		AllowedCapabilities: []string{"none"},
	}

	cfg := &config.Config{}
	cfg.Policy.AllowedCapabilities = []string{"read_files"}
	cfg.Policy.BlockedCapabilities = []string{"network"}

	p.Reload(cfg)
	assert.Equal(t, []string{"read_files"}, p.AllowedCapabilities)
	assert.Equal(t, []string{"network"}, p.BlockedCapabilities)
}

func TestNew(t *testing.T) {
	cfg := &config.Config{}
	cfg.Policy.AllowedCapabilities = []string{"read_files"}
	p := New(cfg)
	assert.NotNil(t, p)
	assert.Equal(t, []string{"read_files"}, p.AllowedCapabilities)
}
