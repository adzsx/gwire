package test

import (
	"fmt"

	"github.com/adzsx/g-wire/pkg/crypt"
)

func Test() {
	fmt.Println(crypt.EncryptAES([]byte("testtesttesdtest"), "This is my secre"))
}
