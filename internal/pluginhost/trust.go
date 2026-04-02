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
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
)

// RegistryRootKey is the pinned public key for the ZTVS Registry.
// In a production build, this would be a real ECDSA public key.
const RegistryRootKey = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEzt8mI+R0D/H7t1h6u1NfXpYv0h5q
kK7z9Xp3oY6v+9X8kLz7X9Xp3oY6v+9X8kLz7X9Xp3oY6v+9X8kLz7Q==
-----END PUBLIC KEY-----`

type ecdsaSignature struct {
	R, S *big.Int
}

// VerifySignature validates that the provided data was signed by the RegistryRootKey.
func VerifySignature(data []byte, signature []byte) error {
	block, _ := pem.Decode([]byte(RegistryRootKey))
	if block == nil {
		return errors.New("failed to parse root public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("parse public key: %v", err)
	}

	ecdsaPub, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("not an ECDSA public key")
	}

	var sig ecdsaSignature
	if _, err := asn1.Unmarshal(signature, &sig); err != nil {
		return fmt.Errorf("unmarshal signature: %v", err)
	}

	h := sha256.Sum256(data)
	if !ecdsa.Verify(ecdsaPub, h[:], sig.R, sig.S) {
		return errors.New("invalid signature")
	}

	return nil
}
