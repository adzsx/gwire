package test

import (
	"fmt"
	"log"

	"github.com/adzsx/g-wire/pkg/utils"
)

func Test() {
	var testSlice [][]string
	for i := 1; i <= 10; i++ {
		testSlice = append(testSlice, []string{fmt.Sprint(i), " element"})
	}

	log.Printf("TestSlice: %v\n Removed: %v\n", testSlice, utils.Remove(testSlice, 10000))

}
