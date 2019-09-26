/*
 * cryptography functions
 */

package utility

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

// Hash builds the sha256 hash of a data and salt byte
func Hash(data, salt []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

// RandomIV returns <size> random bytes
func RandomIV(size int) ([]byte, error) {
	nonce := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return nonce, nil
}

// NewEncryptCFB aes encrypts bytes and returns the cipher followed by a new iv
func NewEncryptCFB(plaintext, passphrase []byte) ([]byte, []byte, error) {
	iv, err := RandomIV(16)
	if err != nil {
		return nil, nil, err
	}
	cipher, err := EncryptCFB(plaintext, iv, passphrase)
	if err != nil {
		return nil, nil, err
	}
	return cipher, iv, nil
}

// EncryptCFB aes encrypts bytes using a predefined iv
func EncryptCFB(plaintext, iv, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)
	return ciphertext, nil
}

// DecryptCFB aes decrypts a cipher when provided with key and iv
func DecryptCFB(ciphertext, iv, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}
