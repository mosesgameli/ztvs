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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifySignature_Trust(t *testing.T) {
	// 1. Setup - generate a temporary key for the test
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	require.NoError(t, err)

	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})

	// 2. Mock the RegistryRootKey for this test
	oldKey := RegistryRootKey
	RegistryRootKey = string(pubPEM)
	defer func() { RegistryRootKey = oldKey }()

	data := []byte("test data")
	h := sha256.Sum256(data)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, h[:])
	require.NoError(t, err)

	signature, err := asn1.Marshal(ecdsaSignature{r, s})
	require.NoError(t, err)

	// 3. Test valid signature
	err = VerifySignature(data, signature)
	assert.NoError(t, err)

	// 4. Test invalid signature
	err = VerifySignature([]byte("tampered data"), signature)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid signature")

	// 5. Test malformed signature
	err = VerifySignature(data, []byte("not a signature"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unmarshal signature")

	// 6. Test invalid public key
	RegistryRootKey = "not a key"
	err = VerifySignature(data, signature)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse root public key")
}
