package utils

import (
	"log"
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

func Err(err error) {
	if err != nil {
		log.Panic("Error: ", err)
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
func VPrint(v any) {
	if time {
		log.SetFlags(log.Ltime)
	}

	if verbose {
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
