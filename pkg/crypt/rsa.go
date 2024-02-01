package crypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"

	"github.com/adzsx/gwire/internal/utils"
)

func GenKeys() rsa.PrivateKey {
	//Generate RSA public and private keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	utils.Err(err, true)

	return *privateKey
}

func EncryptRSA(publicKey rsa.PublicKey, message []byte) []byte {

	// Use SHA and PublicKey to enrypt a []byte message
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &publicKey, message, nil)

	if err != nil {
		panic(err)
	}

	return encryptedBytes
}

func DecryptRSA(privateKey rsa.PrivateKey, encrypted []byte) string {

	// Use SHA and privatekey ro decrypt []byte message
	decryptedBytes, err := privateKey.Decrypt(nil, encrypted, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		panic(err)
	}

	return string(decryptedBytes)
}
