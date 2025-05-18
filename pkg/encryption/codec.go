package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
)

// Codec defines an interface for encryption/decryption
type Codec interface {
	Encrypt([]byte)
	Decrypt([]byte)
}

// AESGCM implements the Codec interface using AES-GCM
type AESGCM struct {
	key []byte
}

// NewAESGCM returns a new AESGCM instance
func NewAESGCM(key []byte) (*AESGCM, error) {
	if len(key) != 32 { // AES-256 requires 32-byte key
		return nil, ErrInvalidKeyLength
	}
	return &AESGCM{key: key}, nil
}

// Encrypt encrypts the input data using AES-GCM in-place
func (a *AESGCM) Encrypt(data []byte) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		log.Printf("failed to create AES cipher: %v", err)
		return
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		log.Printf("failed to create GCM: %v", err)
		return
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Printf("failed to generate nonce: %v", err)
		return
	}

	ciphertext := aead.Seal(nonce, nonce, data, nil)

	// In-place overwrite
	copy(data, ciphertext)
	data = data[:len(ciphertext)]
}

// Decrypt decrypts the input data using AES-GCM in-place
func (a *AESGCM) Decrypt(data []byte) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		log.Printf("failed to create AES cipher: %v", err)
		return
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		log.Printf("failed to create GCM: %v", err)
		return
	}

	nonceSize := aead.NonceSize()
	if len(data) < nonceSize {
		log.Printf("ciphertext too short")
		return
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Printf("failed to decrypt: %v", err)
		return
	}

	// In-place overwrite
	copy(data, plaintext)
	data = data[:len(plaintext)]
}

// Optional error for invalid key length
var ErrInvalidKeyLength = &InvalidKeyError{}

type InvalidKeyError struct{}

func (e *InvalidKeyError) Error() string {
	return "invalid key length: expected 32 bytes for AES-256"
}
