package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

/*
 * cryptography functions
 */
func hash(data, salt []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}
func randomIV(size int) ([]byte, error) {
	nonce := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return nonce, nil
}
func newEncryptCFB(plaintext []byte, passphrase string) ([]byte, []byte, error) {
	iv, err := randomIV(16)
	if err != nil {
		return nil, nil, err
	}
	cipher, err := encryptCFB(plaintext, iv, hash([]byte(passphrase), iv))
	if err != nil {
		return nil, nil, err
	}
	return cipher, iv, nil
}
func encryptCFB(plaintext, iv, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)
	return ciphertext, nil
}
func decryptCFB(ciphertext, iv, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}
