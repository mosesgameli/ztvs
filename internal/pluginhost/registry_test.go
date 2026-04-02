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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/mosesgameli/ztvs/pkg/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGit struct {
	mock.Mock
}

func (m *MockGit) Clone(url, dest string) error {
	args := m.Called(url, dest)
	return args.Error(0)
}

func (m *MockGit) Pull(dest string) error {
	args := m.Called(dest)
	return args.Error(0)
}

func (m *MockGit) Remove(dest string) error {
	args := m.Called(dest)
	return args.Error(0)
}

func TestRegistry_Search(t *testing.T) {
	idx := registry.Index{
		Plugins: []registry.PluginMetadata{
			{Name: "vulnerability-scanner"},
			{Name: "compliance-checker"},
		},
	}
	idxData, _ := json.Marshal(idx)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(idxData)
	}))
	defer server.Close()

	reg := &hostRegistry{BaseURL: server.URL}
	
	t.Run("found", func(t *testing.T) {
		res, err := reg.Search(context.Background(), "scanner")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "vulnerability-scanner", res[0].Name)
	})

	t.Run("not found", func(t *testing.T) {
		res, err := reg.Search(context.Background(), "unknown")
		assert.NoError(t, err)
		assert.Empty(t, res)
	})
}

func TestRegistry_Install(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	// 1. Setup mock index
	pluginName := "test-plugin"
	idx := registry.Index{
		Plugins: []registry.PluginMetadata{
			{
				Name:          pluginName,
				LatestVersion: "1.0.0",
				Repo:          "https://github.com/m/p",
			},
		},
	}
	idxData, _ := json.Marshal(idx)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(idxData)
	}))
	defer server.Close()

	// 2. Setup Host with Lockfile
	lockPath := filepath.Join(tmpDir, ".ztvs", "plugins.lock")
	host := New() // Includes a default registry
	host.lockfile = registry.NewLockfile(lockPath)

	// 3. Mock Git
	mockGit := new(MockGit)
	reg := &hostRegistry{
		BaseURL: server.URL,
		git:     mockGit,
	}

	// 4. Setup mock repo directory structure to simulate successful download and build-not-required
	mockGit.On("Clone", "https://github.com/m/p", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dest := args.String(1)
		pluginDir := filepath.Join(dest, pluginName)
		os.MkdirAll(pluginDir, 0755)
		os.WriteFile(filepath.Join(pluginDir, "plugin.yaml"), []byte("runtime:\n  type: python\n  entrypoint: main.py"), 0644)
		os.WriteFile(filepath.Join(pluginDir, "main.py"), []byte("print('ok')"), 0644)
	})

	err := reg.Install(ctx, pluginName, host)
	assert.NoError(t, err)

	// Verify installation
	installDir := filepath.Join(tmpDir, ".ztvs", "plugins", pluginName)
	_, err = os.Stat(installDir)
	assert.NoError(t, err)

	// Verify lockfile
	l, ok := host.lockfile.Get(pluginName)
	assert.True(t, ok)
	assert.Equal(t, "1.0.0", l.Version)
}

func TestVerifySignature_Logic(t *testing.T) {
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	data := []byte("hello world")
	h := sha256.Sum256(data)
	r, s, _ := ecdsa.Sign(rand.Reader, priv, h[:])
	sig, _ := asn1.Marshal(ecdsaSignature{R: r, S: s})

	pubBytes, _ := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})

	block, _ := pem.Decode(pubPEM)
	pub, _ := x509.ParsePKIXPublicKey(block.Bytes)
	ecdsaPub := pub.(*ecdsa.PublicKey)

	var sigStruct ecdsaSignature
	asn1.Unmarshal(sig, &sigStruct)

	if !ecdsa.Verify(ecdsaPub, h[:], sigStruct.R, sigStruct.S) {
		t.Error("signature verification logic failed")
	}

	// Also call the actual VerifySignature function (it will fail because key is pinned)
	err := VerifySignature(data, sig)
	if err == nil {
		t.Error("expected signature violation for pinned key, got success")
	}
}
