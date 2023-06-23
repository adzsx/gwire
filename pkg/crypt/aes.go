package crypt

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"
)

func EncryptAES(key []byte, plaintext string) string {
	// create cipher
	c, _ := aes.NewCipher(key)

	// allocate space for ciphered data
	out := make([]byte, len(plaintext))

	// encrypt
	c.Encrypt(out, []byte(plaintext))
	// return hex string
	return hex.EncodeToString(out)
}

func DecryptAES(key []byte, ct string) {
	ciphertext, _ := hex.DecodeString(ct)

	c, _ := aes.NewCipher(key)

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	s := string(pt[:])
	fmt.Println("DECRYPTED:", s)
}
