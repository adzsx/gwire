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

func Remove(slice [][]string, index int) [][]string {
	var final [][]string

	for rn, element := range slice {
		if rn != index {
			final = append(final, element)
		}
	}

	return final
}

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

// Verbose print
func VPrint(v any, level int) {
	if time {
		log.SetFlags(log.Ltime)
	}

	if verbose >= level {
		log.Print(v)
	}
}
