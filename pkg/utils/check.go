package utils

import (
	"log"
)

func Err(err error) {
	if err != nil {
		log.Panic(err)
	}
}
