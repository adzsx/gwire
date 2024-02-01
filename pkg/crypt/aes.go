package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func EncryptAES(text string, key []byte) string {
	plaintext := []byte(text)

	// Create new cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Generate IV
	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		panic(err)
	}

	// Encrypt
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	//IDK what this does. I copied it off stack overflow or something
	copy(ciphertext[:aes.BlockSize], iv)

	//Encode cipher to base64
	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)
	return string(encodedCiphertext)
}

func DecryptAES(ciphertext string, key []byte) string {
	// Decode base64
	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		panic(err)
	}

	// Create new cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Extract the IV from the decoded ciphertext
	iv := decodedCiphertext[:aes.BlockSize]
	ciphertextData := decodedCiphertext[aes.BlockSize:]

	// Apply AES decryption to the ciphertext
	plaintext := make([]byte, len(ciphertextData))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext, ciphertextData)

	return string(plaintext)
}
