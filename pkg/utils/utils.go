package utils

import (
	"log"
	"os"
	"strings"
)

func FilterIp(ip string) string {
	var final string

	for _, element := range ip {
		if string(element) != ":" {
			final += string(element)
		} else {
			return final
		}
	}

	return final
}

func FilterPort(ip string) string {
	var final []string

	for index, element := range ip {
		if string(element) == ":" {
			final = append(final, ip[index+1:])
		}
	}

	return strings.Join(final, "")
}

func Err(err error, critical bool) {
	if err != nil {
		log.Println("Error: ", err)
		if critical {
			os.Exit(0)
		}
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

// Verbose print
func Print(v any, level int) {
	if time {
		log.SetFlags(log.Ltime)
	}

	if verbose >= level {
		log.Print(v)
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
