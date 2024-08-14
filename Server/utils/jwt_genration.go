package jwtToken

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func newCipherBlock(key []byte) (cipher.Block, error) {
	return aes.NewCipher(key)
}

// Encrypt encrypts plaintext using AES with the given key
func Encrypt(key []byte, plaintext string) (string, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return "", err
	}

	// Convert plaintext to bytes and pad it to block size
	plainBytes := []byte(plaintext)
	padding := block.BlockSize() - len(plainBytes)%block.BlockSize()
	paddedPlainBytes := append(plainBytes, byte(padding))
	for i := 1; i < padding; i++ {
		paddedPlainBytes = append(paddedPlainBytes, byte(padding))
	}

	ciphertext := make([]byte, aes.BlockSize+len(paddedPlainBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], paddedPlainBytes)

	// Encode IV and encrypted data
	ivEncoded := base64.URLEncoding.EncodeToString(iv)
	ciphertextEncoded := base64.URLEncoding.EncodeToString(ciphertext[aes.BlockSize:])
	return ivEncoded + ":" + ciphertextEncoded, nil
}
