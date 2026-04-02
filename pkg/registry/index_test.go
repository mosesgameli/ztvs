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

package registry

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndex_JSON(t *testing.T) {
	jsonData := `
{
  "version": "1.0",
  "plugins": [
    {
      "name": "test-plugin",
      "latest_version": "1.2.3",
      "repo": "https://github.com/example/test",
      "checksum": "abc12345",
      "signature": "sig-789",
      "audit_status": "verified"
    }
  ]
}
`
	var idx Index
	err := json.Unmarshal([]byte(jsonData), &idx)
	require.NoError(t, err)

	assert.Equal(t, "1.0", idx.Version)
	require.Len(t, idx.Plugins, 1)
	assert.Equal(t, "test-plugin", idx.Plugins[0].Name)
	assert.Equal(t, "1.2.3", idx.Plugins[0].LatestVersion)
	assert.Equal(t, "https://github.com/example/test", idx.Plugins[0].Repo)
	assert.Equal(t, "abc12345", idx.Plugins[0].Checksum)
	assert.Equal(t, "sig-789", idx.Plugins[0].Signature)
	assert.Equal(t, "verified", idx.Plugins[0].AuditStatus)

	// Test empty index
	var emptyIdx Index
	err = json.Unmarshal([]byte(`{"version":"2.0","plugins":[]}`), &emptyIdx)
	assert.NoError(t, err)
	assert.Equal(t, "2.0", emptyIdx.Version)
	assert.Empty(t, emptyIdx.Plugins)
}
