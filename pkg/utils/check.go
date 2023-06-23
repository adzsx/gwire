package utils

import (
	"log"
)

func Err(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func InSlice(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
