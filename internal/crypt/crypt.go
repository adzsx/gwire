package crypt

import (
	"crypto/rand"
	"math/big"

	"github.com/adzsx/gwire/internal/utils"
)

func GenPasswd() (string, error) {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!$%&/()=?+*#-_.:,;"
	ret := make([]byte, 32)
	for i := 0; i < 32; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		utils.Err(err, true)
		ret[i] = chars[num.Int64()]
	}

	return string(ret), nil
}
