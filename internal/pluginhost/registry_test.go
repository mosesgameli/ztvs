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
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mosesgameli/ztvs/pkg/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	assert.NotNil(t, r)
	assert.Contains(t, r.(*hostRegistry).BaseURL, "githubusercontent.com")
}

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
		_, _ = w.Write(idxData)
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
	originalHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", tmpDir)
	defer func() { _ = os.Setenv("HOME", originalHome) }()

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
		_, _ = w.Write(idxData)
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
		_ = os.WriteFile(filepath.Join(pluginDir, "plugin.yaml"), []byte("runtime:\n  type: python\n  entrypoint: main.py"), 0644)
		_ = os.WriteFile(filepath.Join(pluginDir, "main.py"), []byte("print('ok')"), 0644)
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

func TestRegistry_CheckAndUpdateAll(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", tmpDir)
	defer func() { _ = os.Setenv("HOME", originalHome) }()

	pluginName := "test-plugin"
	idx := registry.Index{
		Plugins: []registry.PluginMetadata{
			{
				Name:          pluginName,
				LatestVersion: "1.1.0",
				Repo:          "https://github.com/m/p",
			},
		},
	}
	idxData, _ := json.Marshal(idx)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(idxData)
	}))
	defer server.Close()

	host := New()
	lockPath := filepath.Join(tmpDir, ".ztvs", "plugins.lock")
	host.lockfile = registry.NewLockfile(lockPath)
	host.lockfile.Set(pluginName, registry.PluginLock{
		Version: "1.0.0",
		Enabled: true,
	})

	// Ensure the plugins directory exists for atomic rename
	os.MkdirAll(filepath.Join(tmpDir, ".ztvs", "plugins"), 0755)

	mockGit := new(MockGit)
	reg := &hostRegistry{
		BaseURL: server.URL,
		git:     mockGit,
	}

	mockGit.On("Clone", "https://github.com/m/p", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dest := args.String(1)
		pluginDir := filepath.Join(dest, pluginName)
		os.MkdirAll(pluginDir, 0755)
		_ = os.WriteFile(filepath.Join(pluginDir, "plugin.yaml"), []byte("runtime:\n  type: python"), 0644)
	})

	err := reg.CheckAndUpdateAll(ctx, host, "safe")
	assert.NoError(t, err)

	l, _ := host.lockfile.Get(pluginName)
	assert.Equal(t, "1.1.0", l.Version)
}

func TestRegistry_PerformAtomicUpdate(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", tmpDir)
	defer func() { _ = os.Setenv("HOME", originalHome) }()

	pluginName := "test-plugin"
	meta := &registry.PluginMetadata{
		Name:          pluginName,
		LatestVersion: "1.1.0",
		Repo:          "https://github.com/m/p",
	}

	host := New()
	lockPath := filepath.Join(tmpDir, ".ztvs", "plugins.lock")
	host.lockfile = registry.NewLockfile(lockPath)

	// Create existing plugin dir
	pluginDir := filepath.Join(tmpDir, ".ztvs", "plugins", pluginName)
	os.MkdirAll(pluginDir, 0755)
	_ = os.WriteFile(filepath.Join(pluginDir, "old.txt"), []byte("old"), 0644)

	mockGit := new(MockGit)
	reg := &hostRegistry{git: mockGit}

	mockGit.On("Clone", "https://github.com/m/p", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dest := args.String(1)
		pluginSrcDir := filepath.Join(dest, pluginName)
		os.MkdirAll(pluginSrcDir, 0755)
		_ = os.WriteFile(filepath.Join(pluginSrcDir, "new.txt"), []byte("new"), 0644)
		_ = os.WriteFile(filepath.Join(pluginSrcDir, "plugin.yaml"), []byte("runtime:\n  type: python"), 0644)
	})

	err := reg.PerformAtomicUpdate(ctx, pluginName, host, meta)
	assert.NoError(t, err)

	// Verify new file exists, old file is gone
	assert.FileExists(t, filepath.Join(pluginDir, "new.txt"))
	assert.NoFileExists(t, filepath.Join(pluginDir, "old.txt"))

	l, _ := host.lockfile.Get(pluginName)
	assert.Equal(t, "1.1.0", l.Version)
}

func TestRegistry_BuildPlugin(t *testing.T) {
	oldExec := execCommand
	defer func() { execCommand = oldExec }()

	execCommand = func(command string, args ...string) *exec.Cmd {
		return exec.Command(os.Args[0], append([]string{"-test.run=TestHelperProcess", "--", command}, args...)...)
	}

	tmpDir := t.TempDir()
	err := os.MkdirAll(filepath.Join(tmpDir, "cmd"), 0755)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(tmpDir, "cmd", "main.go"), []byte("package main\nfunc main(){}"), 0644)
	assert.NoError(t, err)

	reg := &hostRegistry{}
	err = reg.buildPlugin(tmpDir, "test-plugin")
	assert.NoError(t, err)

	// Fallback case
	_ = os.RemoveAll(filepath.Join(tmpDir, "cmd"))
	err = reg.buildPlugin(tmpDir, "test-plugin")
	assert.NoError(t, err)
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

func TestRegistry_Errors(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", tmpDir)
	defer func() { _ = os.Setenv("HOME", originalHome) }()
	
	t.Run("FetchIndex_HTTPError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()
		reg := &hostRegistry{BaseURL: server.URL}
		_, err := reg.FetchIndex(ctx)
		assert.Error(t, err)
	})

	t.Run("FetchIndex_InvalidJSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("not-json"))
		}))
		defer server.Close()
		reg := &hostRegistry{BaseURL: server.URL}
		_, err := reg.FetchIndex(ctx)
		assert.Error(t, err)
	})

	t.Run("Install_GitError", func(t *testing.T) {
		mockGit := new(MockGit)
		reg := &hostRegistry{git: mockGit}
		mockGit.On("Clone", mock.Anything, mock.Anything).Return(errors.New("git clone failed"))
		
		host := New()
		err := reg.Install(ctx, "test", host)
		assert.Error(t, err)
	})

	t.Run("Install_MissingFromIndex", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"plugins":[]}`))
		}))
		defer server.Close()
		reg := &hostRegistry{BaseURL: server.URL}
		host := New()
		err := reg.Install(ctx, "missing", host)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "plugin missing not found in registry")
	})

	t.Run("Install_SubdirNotFound", func(t *testing.T) {
		idx := registry.Index{
			Plugins: []registry.PluginMetadata{
				{Name: "test", LatestVersion: "1.0.0", Repo: "url"},
			},
		}
		idxData, _ := json.Marshal(idx)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(idxData)
		}))
		defer server.Close()

		mockGit := new(MockGit)
		reg := &hostRegistry{BaseURL: server.URL, git: mockGit}
		mockGit.On("Clone", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			// Don't create the plugin subdir
			dest := args.String(1)
			os.MkdirAll(dest, 0755)
		})
		
		host := New()
		err := reg.Install(ctx, "test", host)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "plugin subdirectory \"test\" not found")
	})

	t.Run("Install_GitCloneError", func(t *testing.T) {
		mockGit := new(MockGit)
		reg := &hostRegistry{BaseURL: "url", git: mockGit}
		mockGit.On("Clone", mock.Anything, mock.Anything).Return(fmt.Errorf("clone failed"))
		host := New()
		err := reg.Install(ctx, "test", host)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "clone failed")
	})

	t.Run("Install_AlreadyInstalled", func(t *testing.T) {
		tmpHome := t.TempDir()
		oldHome := os.Getenv("HOME")
		os.Setenv("HOME", tmpHome)
		defer os.Setenv("HOME", oldHome)
		
		idx := registry.Index{Plugins: []registry.PluginMetadata{{Name: "test-plugin", Repo: "url"}}}
		idxData, _ := json.Marshal(idx)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write(idxData) }))
		defer server.Close()

		pluginDir := filepath.Join(tmpHome, ".ztvs", "plugins", "test-plugin")
		os.MkdirAll(pluginDir, 0755)
		
		mockGit := new(MockGit)
		reg := &hostRegistry{BaseURL: server.URL, git: mockGit}
		host := New()
		err := reg.Install(ctx, "test-plugin", host)
		assert.NoError(t, err)
	})

	t.Run("Install_NonGoRuntime", func(t *testing.T) {
		idx := registry.Index{Plugins: []registry.PluginMetadata{{Name: "py", Repo: "url", LatestVersion: "1"}}}
		idxData, _ := json.Marshal(idx)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write(idxData) }))
		defer server.Close()

		mockGit := new(MockGit)
		reg := &hostRegistry{BaseURL: server.URL, git: mockGit}
		mockGit.On("Clone", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			dest := args.String(1)
			pDir := filepath.Join(dest, "py")
			os.MkdirAll(pDir, 0755)
			_ = os.WriteFile(filepath.Join(pDir, "plugin.yaml"), []byte("runtime:\n  type: python"), 0644)
		})
		
		host := New()
		err := reg.Install(ctx, "py", host)
		assert.NoError(t, err)
	})

	t.Run("Install_BuildError", func(t *testing.T) {
		oldExec := execCommand
		defer func() { execCommand = oldExec }()
		execCommand = func(command string, args ...string) *exec.Cmd {
			return exec.Command("false")
		}

		idx := registry.Index{Plugins: []registry.PluginMetadata{{Name: "test", Repo: "url", LatestVersion: "1"}}}
		idxData, _ := json.Marshal(idx)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write(idxData) }))
		defer server.Close()

		mockGit := new(MockGit)
		reg := &hostRegistry{BaseURL: server.URL, git: mockGit}
		mockGit.On("Clone", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			dest := args.String(1)
			pDir := filepath.Join(dest, "test")
			os.MkdirAll(pDir, 0755)
			_ = os.WriteFile(filepath.Join(pDir, "Makefile"), []byte("all:\n\tfalse"), 0644)
		})
		
		host := New()
		err := reg.Install(ctx, "test", host)
		assert.Error(t, err)
	})

	t.Run("PerformAtomicUpdate_BuildError", func(t *testing.T) {
		oldExec := execCommand
		defer func() { execCommand = oldExec }()
		execCommand = func(command string, args ...string) *exec.Cmd {
			cmd := exec.Command(os.Args[0], append([]string{"-test.run=TestHelperProcessFailure", "--", command}, args...)...)
			cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
			return cmd
		}

		mockGit := new(MockGit)
		reg := &hostRegistry{git: mockGit}
		mockGit.On("Clone", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			dest := args.String(1)
			os.MkdirAll(filepath.Join(dest, "test"), 0755)
		})

		host := New()
		err := reg.PerformAtomicUpdate(ctx, "test", host, &registry.PluginMetadata{Name: "test", Repo: "url"})
		assert.Error(t, err)
	})

	t.Run("DownloadFile_Errors", func(t *testing.T) {
		err := DownloadFile("/non-existent/path", "http://invalid-url")
		assert.Error(t, err)
		
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("data"))
		}))
		defer server.Close()
		err = DownloadFile("/non-existent/path/file", server.URL)
		assert.Error(t, err)
	})

	t.Run("ExtractRuntimeType_EdgeCases", func(t *testing.T) {
		assert.Equal(t, "", extractRuntimeType([]byte("not-yaml")))
		assert.Equal(t, "python", extractRuntimeType([]byte("runtime:\n  type: python\nother: key")))
		assert.Equal(t, "", extractRuntimeType([]byte("runtime:\nother: key")))
		assert.Equal(t, "python", extractRuntimeType([]byte("runtime:\n  type: python\nstop: true\n  type: nested")))
		// Test stop condition
		assert.Equal(t, "", extractRuntimeType([]byte("runtime:\nstop: true\n  type: nested")))
	})

	t.Run("FetchIndex_MkdirError", func(t *testing.T) {
		tmpDir := t.TempDir()
		// Make cache dir a file to force MkdirAll to fail
		_ = os.WriteFile(filepath.Join(tmpDir, "cache"), []byte("file"), 0644)
		reg := &hostRegistry{BaseURL: "http://any"}
		
		originalHome := os.Getenv("HOME")
		_ = os.Setenv("HOME", tmpDir)
		defer func() { _ = os.Setenv("HOME", originalHome) }()

		_, err := reg.FetchIndex(ctx)
		assert.Error(t, err)
	})

	t.Run("Install_VerifyIntegrityError", func(t *testing.T) {
		idx := registry.Index{Plugins: []registry.PluginMetadata{{Name: "bad", Repo: "url", Checksum: "sha256:wrong"}}}
		idxData, _ := json.Marshal(idx)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write(idxData) }))
		defer server.Close()

		mockGit := new(MockGit)
		reg := &hostRegistry{BaseURL: server.URL, git: mockGit}
		mockGit.On("Clone", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			dest := args.String(1)
			pDir := filepath.Join(dest, "bad")
			os.MkdirAll(pDir, 0755)
			_ = os.WriteFile(filepath.Join(pDir, "bad"), []byte("data"), 0644)
			_ = os.WriteFile(filepath.Join(pDir, "plugin.yaml"), []byte("runtime:\n  type: python"), 0644)
		})
		
		host := New()
		err := reg.Install(ctx, "bad", host)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "integrity check failed")
	})

	t.Run("CheckAndUpdateAll_Locked", func(t *testing.T) {
		reg := &hostRegistry{}
		err := reg.CheckAndUpdateAll(ctx, nil, "locked")
		assert.NoError(t, err)
	})

	t.Run("FetchIndex_SignatureWarning", func(t *testing.T) {
		tmpHome := t.TempDir()
		oldHome := os.Getenv("HOME")
		os.Setenv("HOME", tmpHome)
		defer os.Setenv("HOME", oldHome)

		idx := registry.Index{Plugins: []registry.PluginMetadata{}}
		idxData, _ := json.Marshal(idx)
		
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "index.json") {
				_, _ = w.Write(idxData)
			} else if strings.HasSuffix(r.URL.Path, ".sig") {
				_, _ = w.Write([]byte("invalid-sig"))
			}
		}))
		defer server.Close()

		reg := &hostRegistry{BaseURL: server.URL}
		_, err := reg.FetchIndex(ctx)
		assert.NoError(t, err) // Warning should be printed to stderr, not return error
	})
}
