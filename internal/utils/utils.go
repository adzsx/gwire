package utils

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
)

func Err(err error, critical bool) {
	Ansi("\033[31m")
	if err != nil {
		log.Println("Error:", err)
		if critical {
			os.Exit(0)
		}
	}
	Ansi("\033[0m")
}

func InSlice(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// Verbose print
func Print(v any, level int) {

	Ansi("\033[33m")

	log.SetFlags(log.Ltime)

	if verbose >= level {
		log.Print("System: ", v)
	}

	Ansi("\033[0m")
}

func getOS() string {
	return runtime.GOOS
}

func Ansi(inp string) {
	os := getOS()

	if os != "windows" {
		fmt.Print(inp)
	}
}

func Remove(list []string, str string) []string {
	var final []string
	for _, element := range list {
		if element != str {
			final = append(final, element)
		}
	}

	return final
}

func FilterChar(str string, char string, before bool) string {
	var final string

	for index, element := range str {
		if before {
			if string(element) != char {
				final += string(element)
			} else {
				return final
			}
		} else {

			if string(element) == char {
				final += str[index+1:]
			}
		}

	}

	return final
}

func Nothing() {

}

func aton(letter rune) int {
	uppercaseLetter := strings.ToUpper(string(letter))
	if len(uppercaseLetter) != 1 || uppercaseLetter < "A" || uppercaseLetter > "Z" {
		return 0 // Return 0 for non-letter characters or invalid input
	}

	num := int(uppercaseLetter[0] - 'A' + 1)
	return num
}

func GetRandomString(strings []string, username string) string {
	if len(strings) == 0 {
		return "" // Return an empty string if the input slice is empty
	}

	var seed int64

	for _, char := range username {
		seed += int64(aton(char))
	}

	rand.Seed(seed)

	randomIndex := rand.Intn(len(strings))
	return strings[randomIndex]
}
