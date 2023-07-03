package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/adzsx/gwire/pkg/utils"
)

func EncryptAES(text string, key []byte) string {
	plaintext := []byte(text)

	// Create a new AES cipher block based on the key
	block, err := aes.NewCipher(key)
	utils.Err(err)

	// Generate a new random IV (initialization vector)
	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	utils.Err(err)

	// Apply AES encryption to the plaintext
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// Concatenate the IV and ciphertext
	copy(ciphertext[:aes.BlockSize], iv)

	// Encode the ciphertext in base64 for easier storage and transmission
	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)
	return string(encodedCiphertext)
}

func DecryptAES(ciphertext string, key []byte) string {
	// Decode the base64-encoded ciphertext
	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	utils.Err(err)

	// Create a new AES cipher block based on the key
	block, err := aes.NewCipher(key)
	utils.Err(err)
	// Extract the IV from the decoded ciphertext
	iv := decodedCiphertext[:aes.BlockSize]
	ciphertextData := decodedCiphertext[aes.BlockSize:]

	// Apply AES decryption to the ciphertext
	plaintext := make([]byte, len(ciphertextData))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext, ciphertextData)

	return string(plaintext)
}
