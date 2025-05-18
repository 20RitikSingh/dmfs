package encryption

import (
	"encoding/hex"
	"io"
	"os"

	"github.com/zeebo/blake3"
)

type FileHasher interface {
	hashFile() (string, error)
}

type Blake3Hasher struct {
}

func (h *Blake3Hasher) hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := blake3.New()
	buf := make([]byte, 32*1024) // 32KB buffer

	for {
		n, err := file.Read(buf)
		if n > 0 {
			hasher.Write(buf[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
	}

	sum := hasher.Sum(nil)
	return hex.EncodeToString(sum), nil
}

func (h *Blake3Hasher) Verify(path string, expectedHash string) bool {
	actualHash, err := h.hashFile(path)
	if err != nil {
		return false
	}

	return actualHash == expectedHash
}

func NewBlake3Hasher() *Blake3Hasher {
	return &Blake3Hasher{}
}
