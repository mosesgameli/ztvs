package pluginhost

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// VerifyIntegrity checks if the plug binary's SHA-256 hash matches the provided checksum.
// checksum should be in format "f2416982..." or "sha256:f2416982..."
func VerifyIntegrity(pluginPath, expectedChecksum string) error {
	f, err := os.Open(pluginPath)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	actual := hex.EncodeToString(h.Sum(nil))

	// Handle optional "sha256:" prefix
	if len(expectedChecksum) > 7 && expectedChecksum[:7] == "sha256:" {
		expectedChecksum = expectedChecksum[7:]
	}

	if actual != expectedChecksum {
		return fmt.Errorf("integrity violation: expected %s, got %s", expectedChecksum, actual)
	}

	return nil
}
