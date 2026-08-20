package firmware

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Manifest struct {
	Name      string
	Version   string
	Digest    string
	Size      int64
	Signature string
}

func NewManifest(name, version string, data []byte) Manifest {
	sum := sha256.Sum256(data)
	return Manifest{Name: name, Version: version, Digest: hex.EncodeToString(sum[:]), Size: int64(len(data))}
}
func (m Manifest) Verify(data []byte) bool {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]) == m.Digest && int64(len(data)) == m.Size
}
func (m Manifest) Validate() error {
	if m.Name == "" || m.Version == "" || m.Digest == "" || m.Size <= 0 {
		return fmt.Errorf("invalid firmware manifest")
	}
	return nil
}
