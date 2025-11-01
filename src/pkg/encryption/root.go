package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
)

type Encryptor struct {
	gcm cipher.AEAD
}

func NewEncryptor(key string) (*Encryptor, error) {
	keyHash := sha256.Sum256([]byte(key))

	block, err := aes.NewCipher(keyHash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	return &Encryptor{gcm: gcm}, nil
}

func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	nonce := make([]byte, e.gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := e.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (e *Encryptor) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	nonceSize := e.gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, cipherData := data[:nonceSize], data[nonceSize:]
	plaintext, err := e.gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

func (e *Encryptor) EncryptStringSlice(data []string) ([]string, error) {
	encrypted := make([]string, len(data))
	for i, item := range data {
		enc, err := e.Encrypt(item)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt item %d: %w", i, err)
		}
		encrypted[i] = enc
	}
	return encrypted, nil
}

func (e *Encryptor) DecryptStringSlice(data []string) ([]string, error) {
	decrypted := make([]string, len(data))
	for i, item := range data {
		dec, err := e.Decrypt(item)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt item %d: %w", i, err)
		}
		decrypted[i] = dec
	}
	return decrypted, nil
}
