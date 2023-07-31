package crypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"

	"github.com/adzsx/gwire/pkg/utils"
)

// Generate RSA keys (Private and Public) publicKey = privateKey.PublicKey
func GenKeys() rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	utils.Err(err)

	return *privateKey
}

func EncryptRSA(publicKey rsa.PublicKey, message []byte) []byte {

	// Encrypt message with public Key and hash function
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &publicKey, message, nil)

	if err != nil {
		panic(err)
	}

	return encryptedBytes
}

func DecryptRSA(privateKey rsa.PrivateKey, encrypted []byte) string {

	// Decrypt RSA with private key encrypted bytes and hash function
	decryptedBytes, err := privateKey.Decrypt(nil, encrypted, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		panic(err)
	}

	// We get back the original information in the form of bytes, which we
	// the cast to a string and print
	return string(decryptedBytes)
}
